package bykubernetes

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/metrics/pkg/client/clientset/versioned"
	"neuronet/pkg/k8s/informer"
)

var (
	_      ICoreV1Store = (*coreV1)(nil)
	client ICoreV1Store
)

func GetClient() ICoreV1Store {
	return client
}

type ICoreV1Store interface {
	IPod
	IDeployment
	IService
	IStatefulSet
	IJob
	IEvent
	INode
	INamespace
}

type coreV1 struct {
	kubeConfigPath string
	informer       informer.Storer
	client         kubernetes.Interface
	metricsClient  versioned.Interface
}

func (c *coreV1) Pods() IPodAction {
	return newPods(c.client, c.metricsClient, c.informer)
}

func (c *coreV1) Deployments() IDeploymentAction {
	return newDeployments(c.client, c.informer)
}

func (c *coreV1) Services() IServiceAction {
	return newServices(c.client, c.informer)
}

func (c *coreV1) StatefulSets() IStatefulSetAction {
	return newStatefulSets(c.client, c.informer)
}

func (c *coreV1) Jobs() IJobAction {
	return newJobs(c.client, c.informer)
}

func (c *coreV1) Events() IEventAction {
	return newEvents(c.client)
}

func (c *coreV1) Nodes() INodeAction {
	return newNodes(c.client, c.informer)
}

func (c *coreV1) Namespaces() INamespaceAction {
	return newNamespaces(c.client, c.informer)
}

func NewCoreV1Store(c kubernetes.Interface, metricsClientSet versioned.Interface, informer informer.Storer, kubeConfigPath string) {
	client = &coreV1{kubeConfigPath, informer, c, metricsClientSet}
}
