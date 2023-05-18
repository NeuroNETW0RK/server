package v1

import "neuronet/internal/pkg/cluster"

type SingleNodeArgs struct {
	ClusterName string `uri:"cluster_name" json:"cluster_name"`
	NodeName    string `uri:"node_name" json:"node_name"`
}

type GroupNodesArgs struct {
	ClusterName string `uri:"cluster_name" json:"cluster_name"`
	Label       string `uri:"label" json:"label"`
}

type AllNodesArgs struct {
	ClusterName string `uri:"cluster_name" json:"cluster_name"`
}

type ResourceReply struct {
	Nodes  map[string]*cluster.NodeDescribe `json:"nodes"`
	Remain cluster.ResourceList             `json:"remain"`
	Total  cluster.ResourceList             `json:"total"`
}
