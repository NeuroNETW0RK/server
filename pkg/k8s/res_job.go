package bykubernetes

import (
	"context"
	k8serror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"neuronet/internal/pkg/code"
	"neuronet/pkg/errors"
	"neuronet/pkg/k8s/informer"
	"neuronet/pkg/k8s/meta"

	batchv1 "k8s.io/api/batch/v1"
)

var _ IJobAction = (*jobs)(nil)

type IJob interface {
	Jobs() IJobAction
}

type IJobAction interface {
	Create(ctx context.Context, job *batchv1.Job, createOptions meta.CreateOptions) error
	Delete(ctx context.Context, deleteOptions meta.DeleteOptions) error
	Get(ctx context.Context, getOptions meta.GetOptions) (*batchv1.Job, error)
	List(ctx context.Context, listOptions meta.ListOptions) ([]*batchv1.Job, error)
}

type jobs struct {
	client   kubernetes.Interface
	informer informer.Storer
}

func newJobs(c kubernetes.Interface, informerStore informer.Storer) *jobs {
	return &jobs{
		client:   c,
		informer: informerStore,
	}
}

func (c *jobs) Create(ctx context.Context, job *batchv1.Job, options meta.CreateOptions) error {
	_, err := c.client.BatchV1().
		Jobs(options.Namespace).
		Create(ctx, job, metav1.CreateOptions{})
	if err != nil {
		if k8serror.IsAlreadyExists(err) {
			return errors.WithCode(code.ErrDataExisted, err.Error())
		}
		return err
	}

	return nil
}

func (c *jobs) Delete(ctx context.Context, deleteOptions meta.DeleteOptions) error {
	propagationPolicy := metav1.DeletePropagationBackground
	if err := c.client.BatchV1().
		Jobs(deleteOptions.Namespace).
		Delete(ctx, deleteOptions.ObjectName, metav1.DeleteOptions{PropagationPolicy: &propagationPolicy}); err != nil {
		return err
	}

	return nil
}

func (c *jobs) Get(ctx context.Context, getOptions meta.GetOptions) (*batchv1.Job, error) {

	return c.informer.InformerJobs().Get(ctx, getOptions)
}

func (c *jobs) List(ctx context.Context, options meta.ListOptions) ([]*batchv1.Job, error) {
	var (
		list []*batchv1.Job
		err  error
	)

	if options.Label != "" {
		list, err = c.informer.InformerJobs().ListByLabel(ctx, options.Namespace, options.Label)
		if err != nil {
			return nil, err
		}
		return list, nil

	} else if options.Annotations != "" {
		list, err = c.informer.InformerJobs().ListByAnnotations(ctx, options.Namespace, options.Annotations)
		if err != nil {
			return nil, err
		}
		return list, nil

	}

	list, err = c.informer.InformerJobs().ListAll(ctx)
	if err != nil {
		return nil, err
	}
	return list, nil
}
