package informer

import (
	"context"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	coreinformers "k8s.io/client-go/informers/core/v1"
	"neuronet/internal/pkg/code"
	"neuronet/pkg/errors"
)

var _ NamespaceAction = (*namespaces)(nil)

type Namespace interface {
	InformerNamespaces() NamespaceAction
}

type NamespaceAction interface {
	List(ctx context.Context) ([]*v1.Namespace, error)
}

type namespaces struct {
	client coreinformers.NamespaceInformer
}

func newNamespaces(client coreinformers.NamespaceInformer) *namespaces {
	return &namespaces{
		client: client,
	}
}

func (d *namespaces) List(ctx context.Context) ([]*v1.Namespace, error) {

	namespaceList, err := d.client.Lister().List(labels.Everything())
	if err != nil {
		return nil, err
	}
	if len(namespaceList) == 0 {
		return nil, errors.WithCode(code.ErrDataNotFound, "namespaceList not found")
	}
	return namespaceList, nil
}
