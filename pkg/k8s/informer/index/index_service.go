package index

import (
	"k8s.io/client-go/tools/cache"
)

func ServiceRegister(informer cache.SharedIndexInformer) {
	informer.AddIndexers(cache.Indexers{})
}
