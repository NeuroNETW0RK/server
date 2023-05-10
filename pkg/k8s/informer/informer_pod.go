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

var _ PodAction = (*pods)(nil)

type Pod interface {
	InformerPods() PodAction
}

type PodAction interface {
	ListAll(ctx context.Context) ([]*v1.Pod, error)
	ListByLabel(ctx context.Context, namespace string, label string) ([]*v1.Pod, error)
	Get(ctx context.Context, options meta.GetOptions) (*v1.Pod, error)
}

type pods struct {
	client coreinformers.PodInformer
}

func newPods(client coreinformers.PodInformer) *pods {
	return &pods{
		client: client,
	}
}

func (d *pods) ListAll(ctx context.Context) ([]*v1.Pod, error) {
	list, err := d.client.Lister().List(labels.Everything())
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, errors.WithCode(code.ErrDataNotFound, "nodeList not found")
	}
	return list, nil
}

func (d *pods) ListByLabel(ctx context.Context, namespace string, label string) ([]*v1.Pod, error) {
	selector, err := labels.Parse(label)
	if err != nil {
		return nil, err
	}
	list, err := d.client.Lister().Pods(namespace).List(selector)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, errors.WithCode(code.ErrDataNotFound, "podList with %s label in %s namespace not found", label, namespace)
	}
	return list, nil
}

func (d *pods) Get(ctx context.Context, options meta.GetOptions) (*v1.Pod, error) {
	pod, err := d.client.Lister().Pods(options.Namespace).Get(options.ObjectName)
	if err != nil {
		if k8serror.IsNotFound(err) {
			return nil, errors.WithCode(code.ErrDataNotFound, "pod %s in %s namespace not found", options.ObjectName, options.Namespace)
		}
		return nil, err
	}
	return pod, nil
}
