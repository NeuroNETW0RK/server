package index

import (
	"fmt"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/client-go/tools/cache"
)

const (
	NodeByName      = "NodeByName"
	NodeByGroupName = "NodeByGroupName"
)

func nodeByName(obj interface{}) ([]string, error) {
	o, ok := obj.(*v1.Node)
	if !ok {
		return []string{""}, nil
	}
	return []string{o.GetName()}, nil
}

func nodeByGroupName(obj interface{}) ([]string, error) {
	metadata, err := meta.Accessor(obj)
	if err != nil {
		// ...
		return []string{""}, nil
	}

	var indexKeys []string
	for key, value := range metadata.GetLabels() {
		indexKeys = append(indexKeys, key, fmt.Sprintf("%s=%s", key, value))
	}
	return indexKeys, nil

}

func NodeRegister(informer cache.SharedIndexInformer) {
	informer.AddIndexers(cache.Indexers{
		NodeByName:      nodeByName,
		NodeByGroupName: nodeByGroupName,
	})
}
