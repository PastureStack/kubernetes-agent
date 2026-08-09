package watchevents

import (
	"time"

	"github.com/PastureStack/kubernetes-agent/kubernetesclient"
	"github.com/rancher/go-rancher/v3"
	"github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/tools/cache"
)

type namespaceHandler struct {
	rClient     *client.RancherClient
	kClient     *kubernetesclient.Client
	nsWatchChan chan struct{}
}

const (
	AddedEventType    = "ADDED"
	DeletedEventType  = "DELETED"
	ModifiedEventType = "MODIFIED"
)

func NewNamespaceHandler(rClient *client.RancherClient, kClient *kubernetesclient.Client) *namespaceHandler {
	nsHandler := &namespaceHandler{
		rClient: rClient,
		kClient: kClient,
	}
	return nsHandler
}

func (n *namespaceHandler) startNamespaceWatch() (chan struct{}, error) {
	watchlist := cache.NewListWatchFromClient(n.kClient.K8sClient.CoreV1().RESTClient(), "namespaces", "", fields.Everything())
	_, controller := cache.NewInformer(
		compatibleListWatch(watchlist),
		&v1.Namespace{},
		time.Second*0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: namespaceAddDelete(func(namespace v1.Namespace) {
				logrus.Infof("Skipping event: [%s] for namespace: %s", AddedEventType, namespace.Name)
			}),
			DeleteFunc: namespaceAddDelete(func(namespace v1.Namespace) {
				logrus.Infof("Received event: [%s] for Namespace: %s, Handling event.", DeletedEventType, namespace.Name)
				if err := n.delete(namespace); err != nil {
					logrus.Errorf("Error Handling event: [%s] for namespace: %v", DeletedEventType, err)
				}
			}),
			UpdateFunc: namespaceModify(func(namespace v1.Namespace) {
				logrus.Infof("Skipping event: [%s] for namespace: %s", ModifiedEventType, namespace.Name)
			}),
		},
	)

	return runAndSyncController(controller, "namespace")
}

func (n *namespaceHandler) delete(realNS v1.Namespace) error {
	var metadata metav1.ObjectMeta
	var kind string
	var prefix string
	var serviceEvent = &client.ExternalServiceEvent{}
	prefix = "kubernetes://"
	metadata = realNS.ObjectMeta
	kind = "kubernetesService"
	serviceEvent.Environment = &client.Stack{
		Kind: "environment",
	}
	serviceEvent.ExternalId = prefix + string(metadata.UID)
	serviceEvent.EventType = "stack.remove"

	service := client.Service{
		Kind: kind,
	}
	serviceEvent.Service = service

	_, err := n.rClient.ExternalServiceEvent.Create(serviceEvent)
	return err
}

func (n *namespaceHandler) Start() error {
	logrus.Infof("Starting namespace watch")
	stop, err := n.startNamespaceWatch()
	if err != nil {
		return err
	}
	n.nsWatchChan = stop
	return nil
}

func (n *namespaceHandler) Stop() {
	closeWatch(&n.nsWatchChan)
}

func namespaceAddDelete(f func(v1.Namespace)) func(interface{}) {
	return func(obj interface{}) {
		ns := obj.(*v1.Namespace)
		f(*ns)
	}
}

func namespaceModify(f func(v1.Namespace)) func(interface{}, interface{}) {
	return func(oldObj, newObj interface{}) {
		namespaceAddDelete(f)(newObj)
	}
}
