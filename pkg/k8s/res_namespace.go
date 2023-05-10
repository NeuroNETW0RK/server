package bykubernetes

import (
	"context"
	"k8s.io/client-go/kubernetes"
	"neuronet/pkg/k8s/informer"

	"k8s.io/api/core/v1"
)

var _ INamespaceAction = (*namespaces)(nil)

type INamespace interface {
	Namespaces() INamespaceAction
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
	return c.informer.InformerNamespaces().List(ctx)
}
