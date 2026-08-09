package kubernetesclient

import (
	"fmt"
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	kubeconfigLocation = "/etc/kubernetes/ssl/kubeconfig"
	inClusterConfig    = "INCLUSTER_CONFIG"
)

func GetK8sClientSet(apiURL string) (*kubernetes.Clientset, error) {
	var config *rest.Config
	var err error

	// used for backward compatibility for unit tests
	if inClusterEnv := os.Getenv(inClusterConfig); inClusterEnv == "false" {
		location := os.Getenv("KUBECONFIG")
		if location == "" {
			location = kubeconfigLocation
		}
		config, err = clientcmd.BuildConfigFromFlags("", location)
	} else {
		config, err = rest.InClusterConfig()
	}
	if err != nil {
		return nil, fmt.Errorf("build Kubernetes client configuration: %w", err)
	}
	if apiURL != "" {
		config.Host = apiURL
	}
	// JSON is supported by both Kubernetes 1.12 and current client-go releases.
	// Explicit negotiation keeps proxy and test-server behaviour deterministic.
	config.ContentType = "application/json"
	config.AcceptContentTypes = "application/json"
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("create Kubernetes client: %w", err)
	}
	return clientSet, nil
}
