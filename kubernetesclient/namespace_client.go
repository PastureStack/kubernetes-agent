package kubernetesclient

import (
	"context"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NamespaceOperations interface {
	ByName(name string) (*v1.Namespace, error)
	CreateNamespace(resource *v1.Namespace) (*v1.Namespace, error)
	DeleteNamespace(namespace string) error
}

func newNamespaceClient(client *Client) *NamespaceClient {
	return &NamespaceClient{
		client: client,
	}
}

type NamespaceClient struct {
	client *Client
}

func (c *NamespaceClient) ByName(name string) (*v1.Namespace, error) {
	return c.client.K8sClient.CoreV1().Namespaces().Get(context.Background(), name, metav1.GetOptions{})
}

func (c *NamespaceClient) CreateNamespace(resource *v1.Namespace) (*v1.Namespace, error) {
	return c.client.K8sClient.CoreV1().Namespaces().Create(context.Background(), resource, metav1.CreateOptions{})

}

func (c *NamespaceClient) DeleteNamespace(name string) error {
	return c.client.K8sClient.CoreV1().Namespaces().Delete(context.Background(), name, metav1.DeleteOptions{})
}
