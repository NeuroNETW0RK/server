package bykubernetes

import (
	"context"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"neuronet/pkg/k8s/informer"
	"neuronet/pkg/k8s/meta"
)

type INode interface {
	Nodes() INodeAction
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
	return c.informer.InformerNodes().Get(ctx, getOptions)
}

func (c *nodes) List(ctx context.Context, options meta.ListOptions) ([]*v1.Node, error) {
	var (
		list []*v1.Node
		err  error
	)

	if options.Label != "" {
		list, err = c.informer.InformerNodes().ListByLabel(ctx, options.Label)
		if err != nil {
			return nil, err
		}
		return list, nil

	}

	list, err = c.informer.InformerNodes().ListAll(ctx)
	if err != nil {
		return nil, err
	}
	return list, nil
}
