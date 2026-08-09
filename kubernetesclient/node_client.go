package kubernetesclient

import (
	"context"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NodeOperations interface {
	ByName(name string) (*v1.Node, error)
	ReplaceNode(resource *v1.Node) (*v1.Node, error)
}

func newNodeClient(client *Client) *NodeClient {
	return &NodeClient{
		client: client,
	}
}

type NodeClient struct {
	client *Client
}

func (c *NodeClient) ByName(name string) (*v1.Node, error) {
	return c.client.K8sClient.CoreV1().Nodes().Get(context.Background(), name, metav1.GetOptions{})
}

func (c *NodeClient) ReplaceNode(resource *v1.Node) (*v1.Node, error) {
	return c.client.K8sClient.CoreV1().Nodes().Update(context.Background(), resource, metav1.UpdateOptions{})
}
