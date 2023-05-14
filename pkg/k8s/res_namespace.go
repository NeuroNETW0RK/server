package k8s

import (
	"context"
	"k8s.io/client-go/kubernetes"
	"neuronet/internal/pkg/code"
	"neuronet/pkg/errors"
	"neuronet/pkg/k8s/informer"

	"k8s.io/api/core/v1"
)

var _ INamespaceAction = (*namespaces)(nil)

type INamespace interface {
	Namespaces(clusterName string) INamespaceAction
}

type INamespaceAction interface {
	List(ctx context.Context) ([]*v1.Namespace, error)
}

type namespaces struct {
	client   kubernetes.Interface
	informer informer.Storer
}

func newNamespaces(c kubernetes.Interface, informerStore informer.Storer) *namespaces {
	return &namespaces{
		client:   c,
		informer: informerStore,
	}
}

func (c *namespaces) List(ctx context.Context) ([]*v1.Namespace, error) {
	if c.informer == nil {
		return nil, errors.WithCode(code.ErrInternalServer, "informer is nil")
	}
	return c.informer.InformerNamespaces().List(ctx)
}
