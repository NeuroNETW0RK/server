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
	List(ctx context.Context, options meta.ListOptions) ([]*v1.Pod, error)
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

func (d *pods) List(ctx context.Context, options meta.ListOptions) ([]*v1.Pod, error) {
	var (
		list     = make([]*v1.Pod, 0)
		selector labels.Selector
		err      error
	)
	if len(options.IndexMap) != 0 {
		// 只根据第一个索引查询
		for key, value := range options.IndexMap {
			pods, err := d.client.Informer().GetIndexer().ByIndex(key, value)
			if err != nil {
				return nil, err
			}
			if options.Namespace != "" {
				for _, tmpPod := range pods {
					pod := tmpPod.(*v1.Pod)
					if options.Namespace == pod.Namespace {
						list = append(list, pod)
					}
				}
			} else {
				for _, tmpPod := range pods {
					pod := tmpPod.(*v1.Pod)
					list = append(list, pod)
				}
			}
			if len(list) == 0 {
				return nil, errors.WithCode(code.ErrDataNotFound, "podList not found")
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
	list, err = d.client.Lister().Pods(options.Namespace).List(selector)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, errors.WithCode(code.ErrDataNotFound, "podList with %s label in %s namespace not found", options.Label, options.Namespace)
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
