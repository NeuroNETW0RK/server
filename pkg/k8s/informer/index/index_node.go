package index

import (
	"k8s.io/client-go/tools/cache"
)

func NodeRegister(informer cache.SharedIndexInformer) {
	informer.AddIndexers(cache.Indexers{})
}
