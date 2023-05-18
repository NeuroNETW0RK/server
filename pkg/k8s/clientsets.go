package k8s

import (
	"context"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/metrics/pkg/client/clientset/versioned"
	"neuronet/pkg/k8s/informer"
	"neuronet/pkg/log"
)

type ClientSets struct {
	K8sClient      kubernetes.Interface
	MetricsClient  versioned.Interface
	InformerClient informer.Storer
}

func NewClientSets(ctx context.Context, kubeConfigPath string, stop chan struct{}) (*ClientSets, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		log.C(ctx).Warnf("Can't create config %v", err)
		return nil, err
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.C(ctx).Warnf("Can't create clientSet %v", err)
		return nil, err
	}

	metricSet, err := versioned.NewForConfig(config)
	if err != nil {
		log.C(ctx).Warnf("Can't create metricSet %v", err)
		return nil, err
	}

	informerStore, err := informer.NewInformerStore(stop, clientSet)
	if err != nil {
		log.C(ctx).Warnf("create informer store error: %v", err)
		return nil, err
	}

	return &ClientSets{
		K8sClient:      clientSet,
		MetricsClient:  metricSet,
		InformerClient: informerStore,
	}, nil
}
