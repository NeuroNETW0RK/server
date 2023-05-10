package index

import (
	"k8s.io/client-go/tools/cache"
)

func PodRegister(informer cache.SharedIndexInformer) {
	informer.AddIndexers(cache.Indexers{})
}
