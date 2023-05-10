package informer

import (
	"context"
	batchv1 "k8s.io/api/batch/v1"
	k8serror "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	batchinformers "k8s.io/client-go/informers/batch/v1"
	"neuronet/internal/pkg/code"
	"neuronet/pkg/errors"
	"neuronet/pkg/k8s/informer/index"
	"neuronet/pkg/k8s/meta"
)

var _ JobAction = (*jobs)(nil)

type Job interface {
	InformerJobs() JobAction
}

type JobAction interface {
	ListAll(ctx context.Context) ([]*batchv1.Job, error)
	ListByLabel(ctx context.Context, namespace string, label string) ([]*batchv1.Job, error)
	ListByAnnotations(ctx context.Context, namespace string, annotations string) ([]*batchv1.Job, error)
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

func (d *jobs) ListAll(ctx context.Context) ([]*batchv1.Job, error) {
	list, err := d.client.Lister().List(labels.Everything())
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, errors.WithCode(code.ErrDataNotFound, "nodeList not found")
	}
	return list, nil
}

func (d *jobs) ListByLabel(ctx context.Context, namespace string, label string) ([]*batchv1.Job, error) {
	selector, err := labels.Parse(label)
	if err != nil {
		return nil, err
	}
	list, err := d.client.Lister().Jobs(namespace).List(selector)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, errors.WithCode(code.ErrDataNotFound, "jobList with %s label in %s namespace not found", label, namespace)
	}
	return list, nil
}

func (d *jobs) ListByAnnotations(ctx context.Context, namespace string, annotations string) ([]*batchv1.Job, error) {
	var list = make([]*batchv1.Job, 0)
	jobs, err := d.client.Informer().GetIndexer().ByIndex(index.JobAnnotations, annotations)
	if err != nil {
		return nil, err
	}
	for _, tmpJob := range jobs {
		job := tmpJob.(*batchv1.Job)
		if namespace == job.Namespace {
			list = append(list, job)
		}
	}
	if len(list) == 0 {
		return nil, errors.WithCode(code.ErrDataNotFound, "jobList not found")
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
