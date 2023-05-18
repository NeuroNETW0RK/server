package informer

import (
	"context"
	v1 "k8s.io/api/apps/v1"
	k8serror "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	appsinformers "k8s.io/client-go/informers/apps/v1"
	"neuronet/internal/pkg/code"
	"neuronet/pkg/errors"
	"neuronet/pkg/k8s/meta"
)

var _ StatefulSetAction = (*statefulSets)(nil)

type StatefulSet interface {
	InformerStatefulSets() StatefulSetAction
}

type StatefulSetAction interface {
	List(ctx context.Context, options meta.ListOptions) ([]*v1.StatefulSet, error)
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

func (d *statefulSets) List(ctx context.Context, options meta.ListOptions) ([]*v1.StatefulSet, error) {
	var (
		list     = make([]*v1.StatefulSet, 0)
		selector labels.Selector
		err      error
	)
	if len(options.IndexMap) != 0 {
		// 只根据第一个索引查询
		for key, value := range options.IndexMap {
			statefulSets, err := d.client.Informer().GetIndexer().ByIndex(key, value)
			if err != nil {
				return nil, err
			}
			if options.Namespace != "" {
				for _, tmpStatefulSet := range statefulSets {
					statefulSet := tmpStatefulSet.(*v1.StatefulSet)
					if options.Namespace == statefulSet.Namespace {
						list = append(list, statefulSet)
					}
				}
			} else {
				for _, tmpStatefulSet := range statefulSets {
					statefulSet := tmpStatefulSet.(*v1.StatefulSet)
					list = append(list, statefulSet)
				}
			}
			if len(list) == 0 {
				return nil, errors.WithCode(code.ErrDataNotFound, "statefulSetList not found")
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
	list, err = d.client.Lister().StatefulSets(options.Namespace).List(selector)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, errors.WithCode(code.ErrDataNotFound, "statefulSetList with %s label in %s namespace not found", options.Label, options.Namespace)
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
