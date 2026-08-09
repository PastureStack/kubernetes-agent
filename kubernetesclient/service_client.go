package kubernetesclient

import (
	"context"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ServiceOperations interface {
	ByName(namespace string, name string) (*v1.Service, error)
	CreateService(namespace string, resource *v1.Service) (*v1.Service, error)
	ReplaceService(namespace string, resource *v1.Service) (*v1.Service, error)
	DeleteService(namespace string, name string) error
}

func newServiceClient(client *Client) *ServiceClient {
	return &ServiceClient{
		client: client,
	}
}

type ServiceClient struct {
	client *Client
}

func (c *ServiceClient) ByName(namespace string, name string) (*v1.Service, error) {
	return c.client.K8sClient.CoreV1().Services(namespace).Get(context.Background(), name, metav1.GetOptions{})
}

func (c *ServiceClient) CreateService(namespace string, resource *v1.Service) (*v1.Service, error) {
	return c.client.K8sClient.CoreV1().Services(namespace).Create(context.Background(), resource, metav1.CreateOptions{})
}

func (c *ServiceClient) ReplaceService(namespace string, resource *v1.Service) (*v1.Service, error) {
	return c.client.K8sClient.CoreV1().Services(namespace).Update(context.Background(), resource, metav1.UpdateOptions{})
}

func (c *ServiceClient) DeleteService(namespace string, name string) error {
	return c.client.K8sClient.CoreV1().Services(namespace).Delete(context.Background(), name, metav1.DeleteOptions{})
}
