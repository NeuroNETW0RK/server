package bykubernetes

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/metrics/pkg/client/clientset/versioned"
	"neuronet/pkg/log"
)

type Options struct {
	KubeConfigPath string
}

type IClientSet interface {
	GetKubeConfigPath() string
	GetClient() *kubernetes.Clientset
	GetMetricsClient() *versioned.Clientset
}

var _ IClientSet = (*ClientSet)(nil)

func NewClientSet(opt Options) *ClientSet {
	return &ClientSet{kubeConfigPath: opt.KubeConfigPath}
}

type ClientSet struct {
	kubeConfigPath string
}

func (c *ClientSet) GetKubeConfigPath() string {
	return c.kubeConfigPath
}

func (c *ClientSet) GetClient() *kubernetes.Clientset {
	config, err := clientcmd.BuildConfigFromFlags("", c.kubeConfigPath)
	if err != nil {
		log.Panicf("Can't create clientSet %v", err)
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Panicf("Can't create clientSet %v", err)
	}
	return clientSet
}

func (c *ClientSet) GetMetricsClient() *versioned.Clientset {
	config, err := clientcmd.BuildConfigFromFlags("", c.kubeConfigPath)
	if err != nil {
		log.Panicf("Can't create clientSet %v", err)
	}

	clientSet, err := versioned.NewForConfig(config)
	if err != nil {
		log.Panicf("Can't create clientSet %v", err)
	}
	return clientSet
}
