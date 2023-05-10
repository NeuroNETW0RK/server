package index

import (
	"fmt"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/client-go/tools/cache"
)

const (
	DeploymentAnnotations = "DeploymentAnnotations"
)

func deploymentAnnotations(obj interface{}) ([]string, error) {
	metadata := obj.(*v1.Deployment)

	var indexKeys []string
	for key, value := range metadata.Spec.Template.Annotations {
		indexKeys = append(indexKeys, key, fmt.Sprintf("%s=%s", key, value))

	}
	return indexKeys, nil
}

func DeploymentRegister(informer cache.SharedIndexInformer) {
	informer.AddIndexers(cache.Indexers{
		DeploymentAnnotations: deploymentAnnotations,
	})
}
