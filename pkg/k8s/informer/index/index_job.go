package index

import (
	"fmt"
	v1 "k8s.io/api/batch/v1"
	"k8s.io/client-go/tools/cache"
)

const (
	JobAnnotations = "JobAnnotations"
)

func jobAnnotations(obj interface{}) ([]string, error) {
	metadata := obj.(*v1.Job)

	var indexKeys []string
	for key, value := range metadata.Spec.Template.Annotations {
		indexKeys = append(indexKeys, key, fmt.Sprintf("%s=%s", key, value))

	}
	return indexKeys, nil
}

func JobRegister(informer cache.SharedIndexInformer) {
	informer.AddIndexers(cache.Indexers{
		JobAnnotations: jobAnnotations,
	})
}
