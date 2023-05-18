package cluster

import (
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
	"testing"
)

func getK8sClient() kubernetes.Interface {
	kubeconfigPath := filepath.Join(
		os.Getenv("HOME"), ".kube", "config",
	)

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	return clientset
}

func TestName(t *testing.T) {
	client := getK8sClient()
	nodeList, err := client.CoreV1().Nodes().Get(context.Background(), "cpu", metav1.GetOptions{})
	if err != nil {
		return
	}
	node := &v1.NodeList{Items: []v1.Node{*nodeList}}
	options := metav1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.nodeName=%s", "cpu"),
	}
	podList, err := client.CoreV1().Pods("").List(context.Background(), options)
	if err != nil {
		return
	}

	nodeInfo := NewNodeInfo(node, podList)
	list := nodeInfo.NodesRemain()
	if err != nil {
		return
	}
	fmt.Println("node remain: ", list)

	total := nodeInfo.NodesTotal()
	if err != nil {
		return
	}
	fmt.Println("node total: ", total)
}
