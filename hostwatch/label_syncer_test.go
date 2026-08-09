package hostwatch

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/PastureStack/kubernetes-agent/kubernetesclient"
	cache "github.com/patrickmn/go-cache"
	"github.com/rancher/go-rancher-metadata/metadata"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	metadataHandler *fakeMetadataHandler
	kubeHandler     *fakeKubeNodeHandler
	fakeMetadataURL string
	kubeURL         string
)

type fakeMetadataHandler struct {
	hosts []metadata.Host
}

func (f *fakeMetadataHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	hb, _ := json.Marshal(f.hosts)
	w.Write(hb)
}

type fakeKubeNodeHandler struct {
	nodes map[string]*v1.Node
}

func (f *fakeKubeNodeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// GET Node
	if r.Method == http.MethodGet {
		pathArray := strings.Split(r.URL.Path, "/")
		name := pathArray[len(pathArray)-1]
		w.Header().Set("Content-Type", "application/json")
		hb, _ := json.Marshal(f.nodes[name])
		w.Write(hb)
	}
	// Replace Node
	if r.Method == http.MethodPut {
		node := &v1.Node{}
		data, _ := io.ReadAll(r.Body)
		json.Unmarshal(data, node)
		f.nodes[node.Name] = node
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}
}

func TestMain(m *testing.M) {
	metadataHandler = &fakeMetadataHandler{
		hosts: []metadata.Host{},
	}

	kubeHandler = &fakeKubeNodeHandler{
		nodes: map[string]*v1.Node{},
	}

	metadataMux := http.NewServeMux()
	metadataMux.Handle("/2015-12-19/hosts/", metadataHandler)
	metadataServer := httptest.NewServer(metadataMux)
	fakeMetadataURL = metadataServer.URL + "/2015-12-19"

	kubeMux := http.NewServeMux()
	kubeMux.Handle("/api/v1/nodes/", kubeHandler)
	kubeServer := httptest.NewServer(kubeMux)
	kubeURL = kubeServer.URL

	returnVal := m.Run()
	kubeServer.Close()
	metadataServer.Close()
	os.Exit(returnVal)
}

func newTestKubernetesClient(t *testing.T) *kubernetesclient.Client {
	t.Helper()
	client, err := kubernetesclient.NewClient(kubeURL)
	if err != nil {
		t.Fatalf("create test Kubernetes client: %v", err)
	}
	return client
}

func TestDetectsRemoval(t *testing.T) {
	metadataClient := metadata.NewClient(fakeMetadataURL)
	kubeClient := newTestKubernetesClient(t)
	c := cache.New(1*time.Minute, 1*time.Minute)

	metadataHandler.hosts = []metadata.Host{
		{
			Name:     "test1",
			Hostname: "test1",
			Labels:   map[string]string{},
		},
	}

	kubeHandler.nodes["test1"] = &v1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"test1": "val1",
			},
			Annotations: map[string]string{
				"io.rancher.labels.test1": "",
			},
			Name: "test1",
		},
	}

	labelSync(kubeClient, metadataClient, c)

	if _, ok := kubeHandler.nodes["test1"].Labels["test1"]; ok {
		t.Error("Label test1 was not detected as removed")
	}
}

func TestDetectsAddition(t *testing.T) {
	metadataClient := metadata.NewClient(fakeMetadataURL)
	kubeClient := newTestKubernetesClient(t)
	c := cache.New(1*time.Minute, 1*time.Minute)

	metadataHandler.hosts = []metadata.Host{
		{
			Name:     "test2",
			Hostname: "test2",
			Labels: map[string]string{
				"test2": "val2",
			},
		},
	}

	kubeHandler.nodes["test2"] = &v1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"io.kubernetes.meta": "kube.val",
			},
			Annotations: map[string]string{
				"io.kube.test": "val",
			},
			Name: "test2",
		},
	}

	labelSync(kubeClient, metadataClient, c)

	if _, ok := kubeHandler.nodes["test2"].Labels["test2"]; !ok {
		t.Error("Label test2 was not detected as added")
	}

	if _, ok := kubeHandler.nodes["test2"].Annotations["io.rancher.labels.test2"]; !ok {
		t.Error("Annotation was not set on addition of new label")
	}
}

func TestDetectsChange(t *testing.T) {
	metadataClient := metadata.NewClient(fakeMetadataURL)
	kubeClient := newTestKubernetesClient(t)
	c := cache.New(1*time.Minute, 1*time.Minute)

	metadataHandler.hosts = []metadata.Host{
		{
			Name:     "test3",
			Hostname: "test3",
			Labels: map[string]string{
				"test3": "val3",
			},
		},
	}

	kubeHandler.nodes["test3"] = &v1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"test3": "valx",
			},
			Annotations: map[string]string{
				"io.rancher.labels.test3": "",
			},
			Name: "test3",
		},
	}

	labelSync(kubeClient, metadataClient, c)

	if val := kubeHandler.nodes["test3"].Labels["test3"]; val != "val3" {
		t.Error("Label test3 was not detected as changed")
	}

	if _, ok := kubeHandler.nodes["test3"].Annotations["io.rancher.labels.test3"]; !ok {
		t.Error("Annotation was not set on addition of new label")
	}
}
