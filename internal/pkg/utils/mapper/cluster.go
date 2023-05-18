package mapper

import (
	corev1 "k8s.io/api/core/v1"
	v1 "neuronet/internal/neuronetserver/dto/v1"
	"neuronet/internal/pkg/cluster"
)

func ClusterResourceMapper(nodesList *corev1.NodeList, podList *corev1.PodList) *v1.ResourceReply {
	nodeInfo := cluster.NewNodeInfo(nodesList, podList)
	return &v1.ResourceReply{
		Nodes:  nodeInfo.NodeDescribeList(),
		Remain: nodeInfo.NodesRemain(),
		Total:  nodeInfo.NodesTotal(),
	}
}
