package index

import (
	"fmt"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/client-go/tools/cache"
)

const (
	StatefulsetAnnotations = "StatefulsetAnnotations"
)

func statefulsetAnnotations(obj interface{}) ([]string, error) {
	metadata := obj.(*v1.StatefulSet)

	var indexKeys []string
	for key, value := range metadata.Spec.Template.Annotations {
		indexKeys = append(indexKeys, key, fmt.Sprintf("%s=%s", key, value))

	}
	return indexKeys, nil
}

func StatefulsetRegister(informer cache.SharedIndexInformer) {
	informer.AddIndexers(cache.Indexers{
		StatefulsetAnnotations: statefulsetAnnotations,
	})
}
