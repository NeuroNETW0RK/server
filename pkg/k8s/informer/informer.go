package informer

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/informers"
	appsinformers "k8s.io/client-go/informers/apps/v1"
	batchinformers "k8s.io/client-go/informers/batch/v1"
	coreinformers "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	"neuronet/pkg/k8s/informer/index"
	"time"
)

var (
	_     Storer = (*Store)(nil)
	store *Store
)

func Get() *Store {
	return store
}

func Register(stopCh <-chan struct{}, client kubernetes.Interface) error {
	err := NewInformerStore(stopCh, client)
	if err != nil {
		return err
	}
	return nil
}

type Storer interface {
	Pod
	Deployment
	Service
	StatefulSet
	Job
	Node
	Namespace
}

type Store struct {
	namespaceInformer   coreinformers.NamespaceInformer
	eventInformer       coreinformers.EventInformer
	serviceInformer     coreinformers.ServiceInformer
	nodeInformer        coreinformers.NodeInformer
	podInformer         coreinformers.PodInformer
	jobInformer         batchinformers.JobInformer
	statefulSetInformer appsinformers.StatefulSetInformer
	deploymentInformer  appsinformers.DeploymentInformer
	informerFactory     informers.SharedInformerFactory
}

func NewInformerStore(stopCh <-chan struct{}, clients kubernetes.Interface) error {
	factory := informers.NewSharedInformerFactory(clients, time.Second*30)
	gvrs := []schema.GroupVersionResource{
		{Group: "", Version: "v1", Resource: "pods"},
		{Group: "", Version: "v1", Resource: "services"},
		{Group: "", Version: "v1", Resource: "namespaces"},
		{Group: "", Version: "v1", Resource: "nodes"},

		{Group: "apps", Version: "v1", Resource: "deployments"},
		{Group: "apps", Version: "v1", Resource: "statefulsets"},
		{Group: "apps", Version: "v1", Resource: "daemonsets"},

		{Group: "batch", Version: "v1", Resource: "jobs"},
	}

	for _, v := range gvrs {
		_, err := factory.ForResource(v)
		if err != nil {
			return err
		}
	}
	byIndex := index.New(factory)
	byIndex.Register()

	factory.Start(stopCh)
	factory.WaitForCacheSync(stopCh)

	store = &Store{
		namespaceInformer:   factory.Core().V1().Namespaces(),
		eventInformer:       factory.Core().V1().Events(),
		serviceInformer:     factory.Core().V1().Services(),
		nodeInformer:        factory.Core().V1().Nodes(),
		podInformer:         factory.Core().V1().Pods(),
		jobInformer:         factory.Batch().V1().Jobs(),
		statefulSetInformer: factory.Apps().V1().StatefulSets(),
		deploymentInformer:  factory.Apps().V1().Deployments(),
		informerFactory:     factory,
	}

	return nil
}

func (i *Store) InformerDeployments() DeploymentAction {
	return newDeployments(i.deploymentInformer)
}

func (i *Store) InformerPods() PodAction {
	return newPods(i.podInformer)
}

func (i *Store) InformerServices() ServiceAction {
	return newServices(i.serviceInformer)
}

func (i *Store) InformerStatefulSets() StatefulSetAction {
	return newStatefulSets(i.statefulSetInformer)
}

func (i *Store) InformerJobs() JobAction {
	return newJobs(i.jobInformer)
}

func (i *Store) InformerNodes() NodeAction {
	return newNodes(i.nodeInformer)
}

func (i *Store) InformerNamespaces() NamespaceAction {
	return newNamespaces(i.namespaceInformer)
}
