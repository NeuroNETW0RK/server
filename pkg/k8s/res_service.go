package k8s

import (
	"context"
	v1 "k8s.io/api/core/v1"
	k8serror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"neuronet/internal/pkg/code"
	"neuronet/pkg/errors"
	"neuronet/pkg/k8s/informer"
	"neuronet/pkg/k8s/meta"
)

var _ IServiceAction = (*services)(nil)

type IService interface {
	Services(clusterName string) IServiceAction
}

type IServiceAction interface {
	Create(ctx context.Context, service *v1.Service, options meta.CreateOptions) error
	Update(ctx context.Context, service *v1.Service, options meta.UpdateOptions) error
	Delete(ctx context.Context, options meta.DeleteOptions) error
	Get(ctx context.Context, options meta.GetOptions) (*v1.Service, error)
	List(ctx context.Context, options meta.ListOptions) ([]*v1.Service, error)
}

type services struct {
	client   kubernetes.Interface
	informer informer.Storer
}

func newServices(c kubernetes.Interface, informerStore informer.Storer) *services {
	return &services{
		client:   c,
		informer: informerStore,
	}
}

func (c *services) List(ctx context.Context, options meta.ListOptions) ([]*v1.Service, error) {
	if c.informer == nil {
		return nil, errors.WithCode(code.ErrClusterNotFound, "informer is nil")
	}
	list, err := c.informer.InformerServices().List(ctx, options)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (c *services) Create(ctx context.Context, service *v1.Service, options meta.CreateOptions) error {
	if c.client == nil {
		return errors.WithCode(code.ErrClusterNotFound, "client is nil")
	}
	if _, err := c.client.CoreV1().
		Services(options.Namespace).
		Create(ctx, service, metav1.CreateOptions{}); err != nil {
		if k8serror.IsAlreadyExists(err) {
			return errors.WithCode(code.ErrDataExisted, err.Error())
		}
		return err
	}

	return nil
}

func (c *services) Update(ctx context.Context, service *v1.Service, options meta.UpdateOptions) error {
	if c.client == nil {
		return errors.WithCode(code.ErrClusterNotFound, "client is nil")
	}
	oldService, err := c.Get(ctx, meta.GetOptions{
		Namespace:  options.Namespace,
		ObjectName: options.ObjectName,
	})
	if err != nil {
		return err
	}
	service.SetResourceVersion(oldService.GetResourceVersion())
	service.Spec.ClusterIP = oldService.Spec.ClusterIP
	_, err = c.client.CoreV1().Services(options.Namespace).Update(ctx, service, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (c *services) Delete(ctx context.Context, deleteOptions meta.DeleteOptions) error {
	if c.client == nil {
		return errors.WithCode(code.ErrClusterNotFound, "client is nil")
	}
	err := c.client.CoreV1().
		Services(deleteOptions.Namespace).
		Delete(ctx, deleteOptions.ObjectName, metav1.DeleteOptions{})
	if err != nil {
		if k8serror.IsNotFound(err) {
			return errors.WithCode(code.ErrDataNotFound, err.Error())
		}
		return err
	}

	return err
}

func (c *services) Get(ctx context.Context, getOptions meta.GetOptions) (*v1.Service, error) {
	if c.informer == nil {
		return nil, errors.WithCode(code.ErrClusterNotFound, "informer is nil")
	}
	return c.informer.InformerServices().Get(ctx, getOptions)
}
