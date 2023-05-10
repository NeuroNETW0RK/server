package index

import (
	"k8s.io/client-go/tools/cache"
)

func NamespaceRegister(informer cache.SharedIndexInformer) {
	informer.AddIndexers(cache.Indexers{})
}
