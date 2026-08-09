package kubernetesclient

import (
	"context"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PodOperations interface {
	ByName(namespace string, name string) (*v1.Pod, error)
	CreatePod(namespace string, resource *v1.Pod) (*v1.Pod, error)
	DeletePod(namespace string, name string) error
}

func newPodClient(client *Client) *PodClient {
	return &PodClient{
		client: client,
	}
}

type PodClient struct {
	client *Client
}

func (c *PodClient) ByName(namespace string, name string) (*v1.Pod, error) {
	return c.client.K8sClient.CoreV1().Pods(namespace).Get(context.Background(), name, metav1.GetOptions{})
}

func (c *PodClient) CreatePod(namespace string, resource *v1.Pod) (*v1.Pod, error) {
	return c.client.K8sClient.CoreV1().Pods(namespace).Create(context.Background(), resource, metav1.CreateOptions{})
}

func (c *PodClient) DeletePod(namespace string, name string) error {
	return c.client.K8sClient.CoreV1().Pods(namespace).Delete(context.Background(), name, metav1.DeleteOptions{})
}
