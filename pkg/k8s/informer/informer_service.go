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

var _ ServiceAction = (*services)(nil)

type Service interface {
	InformerServices() ServiceAction
}

type ServiceAction interface {
	ListAll(ctx context.Context) ([]*v1.Service, error)
	ListByLabel(ctx context.Context, namespace string, label string) ([]*v1.Service, error)
	Get(ctx context.Context, options meta.GetOptions) (*v1.Service, error)
}

type services struct {
	client coreinformers.ServiceInformer
}

func newServices(client coreinformers.ServiceInformer) *services {
	return &services{
		client: client,
	}
}

func (d *services) ListAll(ctx context.Context) ([]*v1.Service, error) {
	list, err := d.client.Lister().List(labels.Everything())
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, errors.WithCode(code.ErrDataNotFound, "nodeList not found")
	}
	return list, nil
}

func (d *services) ListByLabel(ctx context.Context, namespace string, label string) ([]*v1.Service, error) {
	selector, err := labels.Parse(label)
	if err != nil {
		return nil, err
	}
	list, err := d.client.Lister().Services(namespace).List(selector)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, errors.WithCode(code.ErrDataNotFound, "serviceList with %s label in %s namespace not found", label, namespace)
	}
	return list, nil
}

func (d *services) Get(ctx context.Context, options meta.GetOptions) (*v1.Service, error) {
	service, err := d.client.Lister().Services(options.Namespace).Get(options.ObjectName)
	if err != nil {
		if k8serror.IsNotFound(err) {
			return nil, errors.WithCode(code.ErrDataNotFound, "service %s in %s namespace not found", options.ObjectName, options.Namespace)
		}
		return nil, err
	}
	return service, nil
}
