package kubernetesclient

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
)

func NewClient(apiURL string) (*Client, error) {
	clientSet, err := GetK8sClientSet(apiURL)
	if err != nil {
		return nil, fmt.Errorf("initialize Kubernetes API client: %w", err)
	}

	client := &Client{
		K8sClient: clientSet,
	}

	client.Pod = newPodClient(client)
	client.Namespace = newNamespaceClient(client)
	client.Service = newServiceClient(client)
	client.Node = newNodeClient(client)

	return client, nil
}

type Client struct {
	K8sClient *kubernetes.Clientset
	Pod       PodOperations
	Namespace NamespaceOperations
	Service   ServiceOperations
	Node      NodeOperations
}
