package k8s

import (
	"context"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"neuronet/internal/pkg/code"
	"neuronet/pkg/errors"
	"neuronet/pkg/k8s/informer"
	"neuronet/pkg/k8s/meta"
)

type INode interface {
	Nodes(clusterName string) INodeAction
}

type INodeAction interface {
	Get(ctx context.Context, getOptions meta.GetOptions) (*v1.Node, error)
	List(ctx context.Context, listOptions meta.ListOptions) ([]*v1.Node, error)
}

type nodes struct {
	client   kubernetes.Interface
	informer informer.Storer
}

func newNodes(c kubernetes.Interface, informerStore informer.Storer) *nodes {
	return &nodes{
		client:   c,
		informer: informerStore,
	}
}

func (c *nodes) Get(ctx context.Context, getOptions meta.GetOptions) (*v1.Node, error) {
	if c.informer == nil {
		return nil, errors.WithCode(code.ErrClusterNotFound, "informer is nil")
	}
	return c.informer.InformerNodes().Get(ctx, getOptions)
}

func (c *nodes) List(ctx context.Context, options meta.ListOptions) ([]*v1.Node, error) {
	if c.informer == nil {
		return nil, errors.WithCode(code.ErrClusterNotFound, "informer is nil")
	}
	list, err := c.informer.InformerNodes().List(ctx, options)
	if err != nil {
		return nil, err
	}
	return list, nil
}
