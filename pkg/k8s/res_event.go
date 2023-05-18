package k8s

import (
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"neuronet/internal/pkg/code"
	"neuronet/pkg/errors"
	apiv1 "neuronet/pkg/k8s/api/v1"
	"neuronet/pkg/k8s/meta"
)

var _ IEventAction = (*events)(nil)

type IEvent interface {
	Events(clusterName string) IEventAction
}

type IEventAction interface {
	List(ctx context.Context, options meta.ListOptions) ([]v1.Event, error)
	Get(ctx context.Context, args apiv1.Event, options meta.GetOptions) ([]v1.Event, error)
}

type events struct {
	client kubernetes.Interface
}

func newEvents(c kubernetes.Interface) *events {
	return &events{
		client: c,
	}
}

func (c *events) List(ctx context.Context, options meta.ListOptions) ([]v1.Event, error) {
	if c.client == nil {
		return nil, errors.WithCode(code.ErrClusterNotFound, "client is nil")
	}
	event, err := c.client.CoreV1().
		Events(options.Namespace).
		List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return event.Items, nil
}

func (c *events) Get(ctx context.Context, args apiv1.Event, options meta.GetOptions) ([]v1.Event, error) {
	if c.client == nil {
		return nil, errors.WithCode(code.ErrClusterNotFound, "client is nil")
	}
	event, err := c.client.CoreV1().
		Events(options.Namespace).
		List(ctx, metav1.ListOptions{
			FieldSelector: fmt.Sprintf("involvedObject.name=%v", options.ObjectName),
			TypeMeta:      metav1.TypeMeta{Kind: args.ResourceType},
		})
	if err != nil {
		return nil, err
	}

	return event.Items, nil
}
