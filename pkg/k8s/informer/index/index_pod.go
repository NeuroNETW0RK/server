package index

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
)

const (
	PodByNodeName = "PodByNodeName"
	PodByName     = "PodByName"
)

func podByNodeName(obj interface{}) ([]string, error) {
	o, ok := obj.(*v1.Pod)
	if !ok {
		return []string{""}, nil
	}
	return []string{o.Spec.NodeName}, nil
}

func podByName(obj interface{}) ([]string, error) {
	o, ok := obj.(*v1.Pod)
	if !ok {
		return []string{""}, nil
	}
	return []string{o.Name}, nil
}

func PodRegister(informer cache.SharedIndexInformer) {
	informer.AddIndexers(cache.Indexers{
		PodByNodeName: podByNodeName,
		PodByName:     podByName,
	})
}
