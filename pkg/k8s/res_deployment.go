package k8s

import (
	"context"
	k8serror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"neuronet/internal/pkg/code"
	"neuronet/pkg/errors"
	"neuronet/pkg/k8s/informer"
	"neuronet/pkg/k8s/meta"
	"time"

	v1 "k8s.io/api/apps/v1"
)

var _ IDeploymentAction = (*deployments)(nil)

type IDeployment interface {
	Deployments(clusterName string) IDeploymentAction
}

type IDeploymentAction interface {
	Create(ctx context.Context, deployment *v1.Deployment, options meta.CreateOptions) error
	Update(ctx context.Context, deployment *v1.Deployment, options meta.UpdateOptions) error
	Delete(ctx context.Context, options meta.DeleteOptions) error
	List(ctx context.Context, options meta.ListOptions) ([]*v1.Deployment, error)
	Get(ctx context.Context, options meta.GetOptions) (*v1.Deployment, error)
	Restart(ctx context.Context, options meta.RestartOptions) error
}

type deployments struct {
	client   kubernetes.Interface
	informer informer.Storer
}

func newDeployments(c kubernetes.Interface, informerStore informer.Storer) *deployments {
	return &deployments{
		client:   c,
		informer: informerStore,
	}
}

func (d *deployments) Restart(ctx context.Context, options meta.RestartOptions) error {
	deployment, err := d.Get(ctx, meta.GetOptions{
		Namespace:  options.Namespace,
		ObjectName: options.ObjectName,
	})
	if err != nil {
		return err
	}
	annotations := deployment.Spec.Template.Annotations
	if annotations == nil {
		annotations = make(map[string]string)
	}
	annotations["kubectl.kubernetes.io/restartedAt"] = metav1.Now().Format(time.RFC3339)
	deployment.Spec.Template.Annotations = annotations
	err = d.Update(ctx, deployment, meta.UpdateOptions{
		Namespace:  options.Namespace,
		ObjectName: options.ObjectName,
	})
	if err != nil {
		return err
	}
	return nil
}

func (d *deployments) Create(ctx context.Context, deployment *v1.Deployment, options meta.CreateOptions) (err error) {
	if d.client == nil {
		return errors.WithCode(code.ErrInternalServer, "client is nil")
	}
	_, err = d.client.AppsV1().Deployments(options.Namespace).Create(ctx, deployment, metav1.CreateOptions{})
	if err != nil {
		if k8serror.IsAlreadyExists(err) {
			return errors.WithCode(code.ErrDataExisted, err.Error())
		}
		return err
	}
	return
}

func (d *deployments) Update(ctx context.Context, deployment *v1.Deployment, options meta.UpdateOptions) (err error) {
	if d.client == nil {
		return errors.WithCode(code.ErrInternalServer, "client is nil")
	}
	_, err = d.client.AppsV1().
		Deployments(options.Namespace).
		Update(ctx, deployment, metav1.UpdateOptions{})

	return
}

func (d *deployments) Delete(ctx context.Context, options meta.DeleteOptions) (err error) {
	if d.client == nil {
		return errors.WithCode(code.ErrInternalServer, "client is nil")
	}

	err = d.client.AppsV1().Deployments(options.Namespace).Delete(ctx, options.ObjectName, metav1.DeleteOptions{})
	if err != nil {
		if k8serror.IsNotFound(err) {
			return errors.WithCode(code.ErrDataExisted, err.Error())
		}
		return err
	}
	return
}

func (d *deployments) List(ctx context.Context, options meta.ListOptions) ([]*v1.Deployment, error) {
	if d.informer == nil {
		return nil, errors.WithCode(code.ErrInternalServer, "informer is nil")
	}

	var (
		list []*v1.Deployment
		err  error
	)

	if options.Label != "" {
		list, err = d.informer.InformerDeployments().ListByLabel(ctx, options.Namespace, options.Label)
		if err != nil {
			return nil, err
		}
		return list, nil

	} else if options.Annotations != "" {
		list, err = d.informer.InformerDeployments().ListByAnnotations(ctx, options.Namespace, options.Annotations)
		if err != nil {
			return nil, err
		}
		return list, nil

	}

	list, err = d.informer.InformerDeployments().ListAll(ctx)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (d *deployments) Get(ctx context.Context, options meta.GetOptions) (*v1.Deployment, error) {
	if d.informer == nil {
		return nil, errors.WithCode(code.ErrInternalServer, "informer is nil")
	}
	return d.informer.InformerDeployments().Get(ctx, options)
}
