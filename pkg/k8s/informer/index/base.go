package index

import (
	"fmt"
	"k8s.io/apimachinery/pkg/api/meta"
)

const (
	IndexByLabels = "IndexByLabels"
)

func Label(obj interface{}) ([]string, error) {
	metadata, err := meta.Accessor(obj)
	if err != nil {
		return []string{""}, err
	}

	var indexKeys []string
	for key, value := range metadata.GetLabels() {
		indexKeys = append(indexKeys, key, fmt.Sprintf("%s=%s", key, value))
	}
	return indexKeys, nil

}
