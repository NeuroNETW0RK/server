package informer

import (
	"context"
	v1 "k8s.io/api/core/v1"
	k8serror "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	coreinformers "k8s.io/client-go/informers/core/v1"
	"neuronet/internal/pkg/code"
	"neuronet/pkg/errors"
	"neuronet/pkg/k8s/meta"
)

var _ NodeAction = (*nodes)(nil)

type Node interface {
	InformerNodes() NodeAction
}

type NodeAction interface {
	ListAll(ctx context.Context) ([]*v1.Node, error)
	ListByLabel(ctx context.Context, label string) ([]*v1.Node, error)
	Get(ctx context.Context, options meta.GetOptions) (*v1.Node, error)
}

type nodes struct {
	client coreinformers.NodeInformer
}

func newNodes(client coreinformers.NodeInformer) *nodes {
	return &nodes{
		client: client,
	}
}

func (d *nodes) ListAll(ctx context.Context) ([]*v1.Node, error) {
	nodeList, err := d.client.Lister().List(labels.Everything())
	if err != nil {
		return nil, err
	}
	if len(nodeList) == 0 {
		return nil, errors.WithCode(code.ErrDataNotFound, "nodeList not found")
	}
	return nodeList, nil
}

func (d *nodes) ListByLabel(ctx context.Context, label string) ([]*v1.Node, error) {
	selector, err := labels.Parse(label)
	if err != nil {
		return nil, err
	}
	nodeList, err := d.client.Lister().List(selector)
	if err != nil {
		return nil, err
	}
	if len(nodeList) == 0 {
		return nil, errors.WithCode(code.ErrDataNotFound, "nodeList with %s label not found", label)
	}
	return nodeList, nil
}

func (d *nodes) Get(ctx context.Context, options meta.GetOptions) (*v1.Node, error) {
	node, err := d.client.Lister().Get(options.ObjectName)
	if err != nil {
		if k8serror.IsNotFound(err) {
			return nil, errors.WithCode(code.ErrDataNotFound, "node %s not found", options.ObjectName)
		}
		return nil, err
	}
	return node, nil
}
