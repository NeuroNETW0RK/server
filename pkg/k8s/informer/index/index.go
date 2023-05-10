package index

import (
	"k8s.io/client-go/informers"
)

type Indexer interface {
	Register()
}

type Index struct {
	sharedInformerFactory informers.SharedInformerFactory
}

func New(sharedInformerFactory informers.SharedInformerFactory) *Index {
	return &Index{sharedInformerFactory}
}

func (b *Index) Register() {
	DeploymentRegister(b.sharedInformerFactory.Apps().V1().Deployments().Informer())
	JobRegister(b.sharedInformerFactory.Batch().V1().Jobs().Informer())
	NamespaceRegister(b.sharedInformerFactory.Core().V1().Namespaces().Informer())
	NodeRegister(b.sharedInformerFactory.Core().V1().Nodes().Informer())
	PodRegister(b.sharedInformerFactory.Core().V1().Pods().Informer())
	ServiceRegister(b.sharedInformerFactory.Core().V1().Services().Informer())
	StatefulsetRegister(b.sharedInformerFactory.Apps().V1().StatefulSets().Informer())
}
