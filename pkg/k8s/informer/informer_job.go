package informer

import (
	"context"
	batchv1 "k8s.io/api/batch/v1"
	k8serror "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	batchinformers "k8s.io/client-go/informers/batch/v1"
	"neuronet/internal/pkg/code"
	"neuronet/pkg/errors"
	"neuronet/pkg/k8s/meta"
)

var _ JobAction = (*jobs)(nil)

type Job interface {
	InformerJobs() JobAction
}

type JobAction interface {
	List(ctx context.Context, options meta.ListOptions) ([]*batchv1.Job, error)
	Get(ctx context.Context, options meta.GetOptions) (*batchv1.Job, error)
}

type jobs struct {
	client batchinformers.JobInformer
}

func newJobs(client batchinformers.JobInformer) *jobs {
	return &jobs{
		client: client,
	}
}

func (d *jobs) List(ctx context.Context, options meta.ListOptions) ([]*batchv1.Job, error) {
	var (
		list     = make([]*batchv1.Job, 0)
		selector labels.Selector
		err      error
	)
	if len(options.IndexMap) != 0 {
		// 只根据第一个索引查询
		for key, value := range options.IndexMap {
			jobs, err := d.client.Informer().GetIndexer().ByIndex(key, value)
			if err != nil {
				return nil, err
			}
			if options.Namespace != "" {
				for _, tmpDeployment := range jobs {
					job := tmpDeployment.(*batchv1.Job)
					if options.Namespace == job.Namespace {
						list = append(list, job)
					}
				}
			} else {
				for _, tmpDeployment := range jobs {
					job := tmpDeployment.(*batchv1.Job)
					list = append(list, job)
				}
			}
			if len(list) == 0 {
				return nil, errors.WithCode(code.ErrDataNotFound, "jobList not found")
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
	list, err = d.client.Lister().Jobs(options.Namespace).List(selector)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, errors.WithCode(code.ErrDataNotFound, "jobList with %s label in %s namespace not found", options.Label, options.Namespace)
	}
	return list, nil
}

func (d *jobs) Get(ctx context.Context, options meta.GetOptions) (*batchv1.Job, error) {
	job, err := d.client.Lister().Jobs(options.Namespace).Get(options.ObjectName)
	if err != nil {
		if k8serror.IsNotFound(err) {
			return nil, errors.WithCode(code.ErrDataNotFound, "job %s in %s namespace not found", options.ObjectName, options.Namespace)
		}
		return nil, err
	}
	return job, nil
}
