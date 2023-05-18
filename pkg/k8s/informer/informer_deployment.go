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

var _ DeploymentAction = (*deployments)(nil)

type Deployment interface {
	InformerDeployments() DeploymentAction
}

type DeploymentAction interface {
	List(ctx context.Context, options meta.ListOptions) ([]*v1.Deployment, error)
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

func (d *deployments) List(ctx context.Context, options meta.ListOptions) ([]*v1.Deployment, error) {
	var (
		list     = make([]*v1.Deployment, 0)
		selector labels.Selector
		err      error
	)
	if len(options.IndexMap) != 0 {
		// 只根据第一个索引查询
		for key, value := range options.IndexMap {
			deployments, err := d.client.Informer().GetIndexer().ByIndex(key, value)
			if err != nil {
				return nil, err
			}
			if options.Namespace != "" {
				for _, tmpDeployment := range deployments {
					deployment := tmpDeployment.(*v1.Deployment)
					if options.Namespace == deployment.Namespace {
						list = append(list, deployment)
					}
				}
			} else {
				for _, tmpDeployment := range deployments {
					deployment := tmpDeployment.(*v1.Deployment)
					list = append(list, deployment)
				}
			}
			if len(list) == 0 {
				return nil, errors.WithCode(code.ErrDataNotFound, "deploymentList not found")
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
	list, err = d.client.Lister().Deployments(options.Namespace).List(selector)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, errors.WithCode(code.ErrDataNotFound, "deploymentList with %s label in %s namespace not found", options.Label, options.Namespace)
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
