package watchevents

import (
	"context"
	"encoding/json"
	"time"

	"github.com/PastureStack/kubernetes-agent/kubernetesclient"
	"github.com/rancher/go-rancher/v3"
	"gopkg.in/check.v1"
	"k8s.io/api/core/v1"
	k8sErr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NamespacehandlerTestSuite struct {
	kClient *kubernetesclient.Client
	events  chan client.ExternalServiceEvent
	handler *namespaceHandler
}

var _ = check.Suite(&NamespacehandlerTestSuite{})

func (s *NamespacehandlerTestSuite) SetUpSuite(c *check.C) {
	s.events = make(chan client.ExternalServiceEvent, 10)
	s.kClient = integrationKubernetesClient(c)
	mock := &MockServiceEventOperations{
		events: s.events,
	}
	mockRancherClient := &client.RancherClient{
		ExternalServiceEvent: mock,
	}

	s.handler = NewNamespaceHandler(mockRancherClient, s.kClient)
	c.Assert(s.handler.Start(), check.IsNil)
}

func (s *NamespacehandlerTestSuite) TearDownSuite(c *check.C) {
	if s.handler != nil {
		s.handler.Stop()
	}
}

func (s *NamespacehandlerTestSuite) TestHandler(c *check.C) {
	nsname := "test-ns-1"
	cleanup_ns(s.kClient, "namespace", "test-ns-1", c)

	meta := metav1.ObjectMeta{Name: nsname}
	ns := &v1.Namespace{
		ObjectMeta: meta,
	}

	respNs, err := s.kClient.Namespace.CreateNamespace(ns)
	if err != nil {
		c.Fatal(err)
	}

	err = s.kClient.Namespace.DeleteNamespace(nsname)
	if err != nil {
		c.Fatal(err)
	}

	// The isolated API gate intentionally does not run a namespace controller.
	// Follow the controller's finalize-then-delete sequence explicitly so the
	// API emits the same final delete event that a full cluster produces.
	terminating, err := s.kClient.K8sClient.CoreV1().Namespaces().Get(context.Background(), nsname, metav1.GetOptions{})
	if err != nil {
		c.Fatal(err)
	}
	payload, err := json.Marshal(map[string]interface{}{
		"apiVersion": "v1",
		"kind":       "Namespace",
		"metadata": map[string]string{
			"name":            nsname,
			"resourceVersion": terminating.ResourceVersion,
		},
		"spec": map[string]interface{}{
			"finalizers": []string{},
		},
	})
	if err != nil {
		c.Fatal(err)
	}
	err = s.kClient.K8sClient.CoreV1().RESTClient().Put().
		Resource("namespaces").
		Name(nsname).
		SubResource("finalize").
		Body(payload).
		Do(context.Background()).
		Error()
	if err != nil {
		c.Fatal(err)
	}
	err = s.kClient.Namespace.DeleteNamespace(nsname)
	if err != nil && !k8sErr.IsNotFound(err) {
		c.Fatal(err)
	}

	var gotDelete bool
	for !gotDelete {
		select {
		case event := <-s.events:
			c.Logf("%#v %+v", event, event)
			svc := event.Service
			service := svc.(client.Service)
			c.Logf("EXPECTED %s; EVENT %s", string(respNs.UID), event)
			if event.EventType == "stack.remove" {
				c.Assert(service.Kind, check.Equals, "kubernetesService")
				c.Assert(event.ExternalId, check.Equals, "kubernetes://"+string(respNs.UID))
				gotDelete = true
			}
		case <-time.After(time.Second * 5):
			c.Fatalf("Timed out waiting for event.")

		}
	}
}

func cleanup_ns(client *kubernetesclient.Client, resourceType string, namespace string, c *check.C) error {
	var err error
	switch resourceType {
	case "namespace":
		err = client.Namespace.DeleteNamespace(namespace)
	default:
		c.Fatalf("Unknown type for cleanup: %s", resourceType)
	}
	if err != nil {
		if k8sErr.IsNotFound(err) {
			return nil
		} else {
			return err
		}
	}
	return nil
}
