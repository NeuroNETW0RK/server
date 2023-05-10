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

var _ DeploymentAction = (*deployments)(nil)

type Deployment interface {
	InformerDeployments() DeploymentAction
}

type DeploymentAction interface {
	ListAll(ctx context.Context) ([]*v1.Deployment, error)
	ListByLabel(ctx context.Context, namespace string, label string) ([]*v1.Deployment, error)
	ListByAnnotations(ctx context.Context, namespace string, annotations string) ([]*v1.Deployment, error)
	Get(ctx context.Context, options meta.GetOptions) (*v1.Deployment, error)
}

type deployments struct {
	client appsinformers.DeploymentInformer
}

func newDeployments(client appsinformers.DeploymentInformer) *deployments {
	return &deployments{
		client: client,
	}
}

func (d *deployments) ListAll(ctx context.Context) ([]*v1.Deployment, error) {
	list, err := d.client.Lister().List(labels.Everything())
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, errors.WithCode(code.ErrDataNotFound, "nodeList not found")
	}
	return list, nil
}

func (d *deployments) ListByLabel(ctx context.Context, namespace string, label string) ([]*v1.Deployment, error) {
	selector, err := labels.Parse(label)
	if err != nil {
		return nil, err
	}
	list, err := d.client.Lister().Deployments(namespace).List(selector)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, errors.WithCode(code.ErrDataNotFound, "deploymentList with %s label in %s namespace not found", label, namespace)
	}
	return list, nil
}

func (d *deployments) ListByAnnotations(ctx context.Context, namespace string, annotations string) ([]*v1.Deployment, error) {
	var list = make([]*v1.Deployment, 0)
	deployments, err := d.client.Informer().GetIndexer().ByIndex(index.DeploymentAnnotations, annotations)
	if err != nil {
		return nil, err
	}
	for _, tmpDeployment := range deployments {
		deployment := tmpDeployment.(*v1.Deployment)
		if namespace == deployment.Namespace {
			list = append(list, deployment)
		}
	}
	if len(list) == 0 {
		return nil, errors.WithCode(code.ErrDataNotFound, "deploymentList not found")
	}
	return list, nil
}

func (d *deployments) Get(ctx context.Context, options meta.GetOptions) (*v1.Deployment, error) {
	deployment, err := d.client.Lister().Deployments(options.Namespace).Get(options.ObjectName)
	if err != nil {
		if k8serror.IsNotFound(err) {
			return nil, errors.WithCode(code.ErrDataNotFound, "deployment %s in %s namespace not found", options.ObjectName, options.Namespace)
		}
		return nil, err
	}
	return deployment, nil
}
