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
	List(ctx context.Context, options meta.ListOptions) ([]*v1.Service, error)
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

func (d *services) List(ctx context.Context, options meta.ListOptions) ([]*v1.Service, error) {
	var (
		list     = make([]*v1.Service, 0)
		selector labels.Selector
		err      error
	)
	if len(options.IndexMap) != 0 {
		// 只根据第一个索引查询
		for key, value := range options.IndexMap {
			services, err := d.client.Informer().GetIndexer().ByIndex(key, value)
			if err != nil {
				return nil, err
			}
			if options.Namespace != "" {
				for _, tmpService := range services {
					service := tmpService.(*v1.Service)
					if options.Namespace == service.Namespace {
						list = append(list, service)
					}
				}
			} else {
				for _, tmpService := range services {
					service := tmpService.(*v1.Service)
					list = append(list, service)
				}
			}
			if len(list) == 0 {
				return nil, errors.WithCode(code.ErrDataNotFound, "serviceList not found")
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
	list, err = d.client.Lister().Services(options.Namespace).List(selector)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, errors.WithCode(code.ErrDataNotFound, "serviceList with %s label in %s namespace not found", options.Label, options.Namespace)
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
