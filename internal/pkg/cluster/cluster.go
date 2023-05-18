package cluster

import (
	corev1 "k8s.io/api/core/v1"
	"strings"
)

var _ INode = (*Node)(nil)

type INode interface {
	NodeDescribeList() map[string]*NodeDescribe
	NodesTotal() ResourceList
	NodesRemain() ResourceList
}

func NewNodeInfo(nodesList *corev1.NodeList, podList *corev1.PodList) *Node {
	return &Node{
		nodesList: nodesList,
		podList:   podList,
	}
}

type Node struct {
	nodesList *corev1.NodeList
	podList   *corev1.PodList
}

func (n *Node) NodeDescribeList() map[string]*NodeDescribe {
	nodesInfoMap := make(map[string]*NodeDescribe)
	for _, node := range n.nodesList.Items {
		podList := n.getNodePodsList(node.Name, n.podList)
		nodesInfoMap[node.Name] = &NodeDescribe{
			Name:              node.Name,
			Role:              n.getNodeRoles(&node),
			Labels:            node.GetLabels(),
			Annotations:       node.GetAnnotations(),
			Taints:            node.Spec.Taints,
			Unschedulable:     node.Spec.Unschedulable,
			Conditions:        node.Status.Conditions,
			Address:           node.Status.Addresses,
			Capacity:          n.parseK8sResourceList(node.Status.Capacity),
			Allocatable:       n.parseK8sResourceList(node.Status.Allocatable),
			SystemInfo:        node.Status.NodeInfo,
			PodCIDR:           node.Spec.PodCIDR,
			PodCIDRs:          node.Spec.PodCIDRs,
			NonTerminatedPods: n.getNonTerminatedPods(*podList),
		}
		nodesInfoMap[node.Name].AllocatedResources = n.getAllocatedResources(nodesInfoMap[node.Name].NonTerminatedPods)
	}
	return nodesInfoMap
}

func (n *Node) NodesTotal() ResourceList {
	resource := make(ResourceList)
	for _, node := range n.nodesList.Items {
		for key, value := range n.parseK8sResourceList(node.Status.Allocatable) {
			resource[key] += value
		}
	}
	return resource
}

func (n *Node) NodesRemain() ResourceList {
	requestsResource := make(ResourceList)
	total := n.NodesTotal()
	allocatedResources := n.getAllocatedResources(n.getNonTerminatedPods(*n.podList))
	for key, value := range total {
		if l, ok := allocatedResources.Request[key]; ok {
			requestsResource[key] = value - l
		}
	}
	return requestsResource
}

func (n *Node) getNodeRoles(node *corev1.Node) string {
	var roles string
	labels := node.GetLabels()
	for key := range labels {
		if strings.Contains(key, "node-role") {
			subStr := strings.Split(key, "/")
			roles = subStr[len(subStr)-1]
		}
	}
	return roles
}

func (n *Node) getNonTerminatedPods(pods corev1.PodList) []NonTerminatedPod {
	var (
		nonTerminatedPods []NonTerminatedPod
	)

	for _, pod := range pods.Items {
		if pod.Status.Phase != corev1.PodRunning {
			continue
		}
		limitsResourceList := make(ResourceList)
		requestResourceList := make(ResourceList)
		for _, container := range pod.Spec.Containers {
			limitsResourceList = n.parseK8sResourceList(container.Resources.Limits)
			requestResourceList = n.parseK8sResourceList(container.Resources.Requests)
		}
		nonTerminatedPods = append(nonTerminatedPods, NonTerminatedPod{
			NameSpace: pod.Namespace,
			Name:      pod.Name,
			Limits:    limitsResourceList,
			Request:   requestResourceList,
		})
	}
	return nonTerminatedPods
}

func (n *Node) getAllocatedResources(nonTerminatedPods []NonTerminatedPod) AllocatedResources {
	var (
		limitsResourceList  = make(ResourceList)
		requestResourceList = make(ResourceList)
		allocatedResources  AllocatedResources
	)

	for _, pod := range nonTerminatedPods {
		for key, value := range pod.Limits {
			limitsResourceList[key] += value
		}
		for key, value := range pod.Request {
			requestResourceList[key] += value
		}
	}
	allocatedResources.Limits = limitsResourceList
	allocatedResources.Request = requestResourceList
	return allocatedResources
}

func (n *Node) getNodePodsList(nodeName string, allPodsList *corev1.PodList) *corev1.PodList {
	var nodePodsList []corev1.Pod
	for _, item := range allPodsList.Items {
		if item.Spec.NodeName == nodeName {
			nodePodsList = append(nodePodsList, item)
		}
	}
	return &corev1.PodList{Items: nodePodsList}
}

func (n *Node) parseK8sResourceList(list corev1.ResourceList) ResourceList {
	resourceList := make(ResourceList)
	for key, value := range list {
		switch key {
		case corev1.ResourceCPU:
			resourceList[string(key)] += value.MilliValue()
		case corev1.ResourceMemory:
			resourceList[string(key)] += value.Value()
		case "nvidia.com/gpu":
			resourceList[string(key)] += value.Value()
		default:
			resourceList[string(key)] += value.Value()
		}
	}
	return resourceList
}
