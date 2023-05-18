package clusterresource

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	v1 "neuronet/internal/neuronetserver/dto/v1"
	"neuronet/internal/pkg/utils/mapper"
	"neuronet/pkg/k8s"
	"neuronet/pkg/k8s/informer/index"
	"neuronet/pkg/k8s/meta"
	"neuronet/pkg/log"
)

var _ Service = (*service)(nil)

type Service interface {
	SingleNodes(c context.Context, args *v1.SingleNodeArgs) (*v1.ResourceReply, error)
	GroupNodes(c context.Context, args *v1.GroupNodesArgs) (*v1.ResourceReply, error)
	AllNodes(c context.Context, args *v1.AllNodesArgs) (*v1.ResourceReply, error)
}

type service struct {
}

func NewService() *service {
	return &service{}
}

func (s *service) SingleNodes(c context.Context, args *v1.SingleNodeArgs) (*v1.ResourceReply, error) {
	var podItems []corev1.Pod
	nodeGetOptions := meta.GetOptions{
		ObjectName: args.NodeName,
	}
	node, err := k8s.GetClient().Nodes(args.ClusterName).Get(c, nodeGetOptions)
	if err != nil {
		log.C(c).Warnf("get node info failed: %v", err)
		return nil, err
	}

	podIndexMap := map[string]string{
		index.PodByNodeName: args.NodeName,
	}

	podListOptions := meta.ListOptions{
		ListSelector: meta.ListSelector{
			IndexMap: podIndexMap,
		},
	}
	pods, err := k8s.GetClient().Pods(args.ClusterName).List(c, podListOptions)
	if err != nil {
		log.C(c).Warnf("get pod info by node name failed: %v", err)
		return nil, err
	}

	nodeList := &corev1.NodeList{
		Items: []corev1.Node{*node},
	}

	for _, pod := range pods {
		podItems = append(podItems, *pod)
	}

	podList := &corev1.PodList{
		Items: podItems,
	}

	return mapper.ClusterResourceMapper(nodeList, podList), nil
}

func (s *service) GroupNodes(c context.Context, args *v1.GroupNodesArgs) (*v1.ResourceReply, error) {
	var (
		nodeItems     []corev1.Node
		groupNodeName []string
		groupPod      []corev1.Pod
	)

	nodeListOptions := meta.ListOptions{
		ListSelector: meta.ListSelector{
			Label: args.Label,
		},
	}
	nodes, err := k8s.GetClient().Nodes(args.ClusterName).List(c, nodeListOptions)
	if err != nil {
		log.C(c).Warnf("get node info by label failed: %v", err)
		return nil, err
	}
	for _, item := range nodes {
		nodeItems = append(nodeItems, *item)
		groupNodeName = append(groupNodeName, item.Name)
	}

	pods, err := k8s.GetClient().Pods(args.ClusterName).List(c, meta.ListOptions{})
	if err != nil {
		log.C(c).Warnf("get pod info by label failed: %v", err)
		return nil, err
	}

	for _, item := range pods {
		for _, nodeName := range groupNodeName {
			if item.Spec.NodeName == nodeName {
				groupPod = append(groupPod, *item)
			}
		}
	}
	podList := &corev1.PodList{Items: groupPod}
	nodeList := &corev1.NodeList{Items: nodeItems}

	return mapper.ClusterResourceMapper(nodeList, podList), nil
}

func (s *service) AllNodes(c context.Context, args *v1.AllNodesArgs) (*v1.ResourceReply, error) {
	var (
		podItems  []corev1.Pod
		nodeItems []corev1.Node
	)
	nodes, err := k8s.GetClient().Nodes(args.ClusterName).List(c, meta.ListOptions{})
	if err != nil {
		log.C(c).Warnf("get node info failed: %v", err)
		return nil, err
	}
	pods, err := k8s.GetClient().Pods(args.ClusterName).List(c, meta.ListOptions{})
	if err != nil {
		log.C(c).Warnf("get pod info failed: %v", err)
		return nil, err
	}

	for _, node := range nodes {
		nodeItems = append(nodeItems, *node)
	}
	for _, pod := range pods {
		podItems = append(podItems, *pod)
	}

	nodeList := &corev1.NodeList{
		Items: nodeItems,
	}
	podList := &corev1.PodList{
		Items: podItems,
	}

	return mapper.ClusterResourceMapper(nodeList, podList), nil
}
