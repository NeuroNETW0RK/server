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
	List(ctx context.Context, options meta.ListOptions) ([]*v1.Node, error)
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

func (d *nodes) List(ctx context.Context, options meta.ListOptions) ([]*v1.Node, error) {
	var (
		list     = make([]*v1.Node, 0)
		selector labels.Selector
		err      error
	)
	if len(options.IndexMap) != 0 {
		// 只根据第一个索引查询
		for key, value := range options.IndexMap {
			nodes, err := d.client.Informer().GetIndexer().ByIndex(key, value)
			if err != nil {
				return nil, err
			}
			if options.Namespace != "" {
				for _, tmpNode := range nodes {
					node := tmpNode.(*v1.Node)
					if options.Namespace == node.Namespace {
						list = append(list, node)
					}
				}
			} else {
				for _, tmpNode := range nodes {
					node := tmpNode.(*v1.Node)
					list = append(list, node)
				}
			}
			if len(list) == 0 {
				return nil, errors.WithCode(code.ErrDataNotFound, "nodeList not found")
			}
			return list, nil
		}
	}

	if options.Label != "" {
		selector, err = labels.Parse(options.Label)
		if err != nil {
			return nil, errors.WithCode(code.ErrInternalServer, "label parse error")
		}
	} else {
		selector = labels.Everything()
	}
	list, err = d.client.Lister().List(selector)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, errors.WithCode(code.ErrDataNotFound, "nodeList with %s label not found", options.Label)
	}
	return list, nil
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
