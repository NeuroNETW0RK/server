package k8s

import (
	"context"
	k8serror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"neuronet/internal/pkg/code"
	"neuronet/pkg/errors"
	"neuronet/pkg/k8s/informer"
	"neuronet/pkg/k8s/meta"

	v1 "k8s.io/api/apps/v1"
	"k8s.io/client-go/kubernetes"
)

var _ IStatefulSetAction = (*statefulSets)(nil)

type IStatefulSet interface {
	StatefulSets(clusterName string) IStatefulSetAction
}

type IStatefulSetAction interface {
	Create(ctx context.Context, statefulSet *v1.StatefulSet, options meta.CreateOptions) error
	Update(ctx context.Context, statefulSet *v1.StatefulSet, options meta.UpdateOptions) error
	Delete(ctx context.Context, options meta.DeleteOptions) error
	Get(ctx context.Context, options meta.GetOptions) (*v1.StatefulSet, error)
	List(ctx context.Context, options meta.ListOptions) ([]*v1.StatefulSet, error)
}

type statefulSets struct {
	client   kubernetes.Interface
	informer informer.Storer
}

func newStatefulSets(c kubernetes.Interface, informerStore informer.Storer) *statefulSets {
	return &statefulSets{
		client:   c,
		informer: informerStore,
	}
}

func (c *statefulSets) Create(ctx context.Context, statefulSet *v1.StatefulSet, options meta.CreateOptions) error {
	if c.client == nil {
		return errors.WithCode(code.ErrClusterNotFound, "client is nil")
	}
	_, err := c.client.AppsV1().
		StatefulSets(options.Namespace).
		Create(ctx, statefulSet, metav1.CreateOptions{})
	if err != nil {
		if k8serror.IsAlreadyExists(err) {
			return errors.WithCode(code.ErrDataExisted, err.Error())
		}
		return err
	}

	return nil
}

func (c *statefulSets) Update(ctx context.Context, statefulSet *v1.StatefulSet, options meta.UpdateOptions) error {
	if c.client == nil {
		return errors.WithCode(code.ErrClusterNotFound, "client is nil")
	}
	if _, err := c.client.AppsV1().StatefulSets(options.Namespace).Update(ctx, statefulSet, metav1.UpdateOptions{}); err != nil {
		return err
	}

	return nil
}

func (c *statefulSets) Delete(ctx context.Context, options meta.DeleteOptions) error {
	if c.client == nil {
		return errors.WithCode(code.ErrClusterNotFound, "client is nil")
	}
	if err := c.client.AppsV1().StatefulSets(options.Namespace).Delete(ctx, options.ObjectName, metav1.DeleteOptions{}); err != nil {
		return err
	}

	return nil
}

func (c *statefulSets) Get(ctx context.Context, options meta.GetOptions) (*v1.StatefulSet, error) {
	if c.informer == nil {
		return nil, errors.WithCode(code.ErrClusterNotFound, "informer is nil")
	}
	return c.informer.InformerStatefulSets().Get(ctx, options)
}

func (c *statefulSets) List(ctx context.Context, options meta.ListOptions) ([]*v1.StatefulSet, error) {
	if c.informer == nil {
		return nil, errors.WithCode(code.ErrClusterNotFound, "informer is nil")
	}
	list, err := c.informer.InformerStatefulSets().List(ctx, options)
	if err != nil {
		return nil, err
	}
	return list, nil
}
