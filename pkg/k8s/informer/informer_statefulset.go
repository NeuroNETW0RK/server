package informer

import (
	"context"
	v1 "k8s.io/api/apps/v1"
	k8serror "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	appsinformers "k8s.io/client-go/informers/apps/v1"
	"neuronet/internal/pkg/code"
	"neuronet/pkg/errors"
	"neuronet/pkg/k8s/informer/index"
	"neuronet/pkg/k8s/meta"
)

var _ StatefulSetAction = (*statefulSets)(nil)

type StatefulSet interface {
	InformerStatefulSets() StatefulSetAction
}

type StatefulSetAction interface {
	ListAll(ctx context.Context) ([]*v1.StatefulSet, error)
	ListByLabel(ctx context.Context, namespace string, label string) ([]*v1.StatefulSet, error)
	ListByAnnotations(ctx context.Context, namespace string, annotations string) ([]*v1.StatefulSet, error)
	Get(ctx context.Context, options meta.GetOptions) (*v1.StatefulSet, error)
}

type statefulSets struct {
	client appsinformers.StatefulSetInformer
}

func newStatefulSets(client appsinformers.StatefulSetInformer) *statefulSets {
	return &statefulSets{
		client: client,
	}
}

func (d *statefulSets) ListAll(ctx context.Context) ([]*v1.StatefulSet, error) {
	list, err := d.client.Lister().List(labels.Everything())
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, errors.WithCode(code.ErrDataNotFound, "nodeList not found")
	}
	return list, nil
}

func (d *statefulSets) ListByLabel(ctx context.Context, namespace string, label string) ([]*v1.StatefulSet, error) {
	selector, err := labels.Parse(label)
	if err != nil {
		return nil, err
	}
	list, err := d.client.Lister().StatefulSets(namespace).List(selector)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, errors.WithCode(code.ErrDataNotFound, "statefulSetList with %s label in %s namespace not found", label, namespace)
	}
	return list, nil
}

func (d *statefulSets) ListByAnnotations(ctx context.Context, namespace string, annotations string) ([]*v1.StatefulSet, error) {
	var list = make([]*v1.StatefulSet, 0)
	statefulSets, err := d.client.Informer().GetIndexer().ByIndex(index.StatefulsetAnnotations, annotations)
	if err != nil {
		return nil, err
	}
	for _, tmpStatefulSet := range statefulSets {
		statefulSet := tmpStatefulSet.(*v1.StatefulSet)
		if namespace == statefulSet.Namespace {
			list = append(list, statefulSet)
		}
	}
	if len(list) == 0 {
		return nil, errors.WithCode(code.ErrDataNotFound, "statefulSetList not found")
	}
	return list, nil
}

func (d *statefulSets) Get(ctx context.Context, options meta.GetOptions) (*v1.StatefulSet, error) {
	statefulSet, err := d.client.Lister().StatefulSets(options.Namespace).Get(options.ObjectName)
	if err != nil {
		if k8serror.IsNotFound(err) {
			return nil, errors.WithCode(code.ErrDataNotFound, "statefulSet %s in %s namespace not found", options.ObjectName, options.Namespace)
		}
		return nil, err
	}
	return statefulSet, nil
}
