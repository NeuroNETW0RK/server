package cluster

import (
	corrdinv1 "k8s.io/api/coordination/v1"
	corev1 "k8s.io/api/core/v1"
)

type ResourceList map[string]int64

type NonTerminatedPod struct {
	NameSpace string       `json:"namespace"`
	Name      string       `json:"name"`
	Limits    ResourceList `json:"limits"`
	Request   ResourceList `json:"request"`
}

type AllocatedResources struct {
	Limits  ResourceList `json:"limits"`
	Request ResourceList `json:"request"`
}

type NodeDescribe struct {
	Name        string            `json:"name"`
	Role        string            `json:"role"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`

	Taints             []corev1.Taint         `json:"taints"`
	Unschedulable      bool                   `json:"unschedulable"`
	Lease              corrdinv1.Lease        `json:"lease,omitempty"`
	Conditions         []corev1.NodeCondition `json:"conditions"`
	Address            []corev1.NodeAddress   `json:"address"`
	Capacity           ResourceList           `json:"capacity"`
	Allocatable        ResourceList           `json:"allocatable"`
	SystemInfo         corev1.NodeSystemInfo  `json:"systemInfo"`
	PodCIDR            string                 `json:"podCIDR"`
	PodCIDRs           []string               `json:"podCIDRs"`
	NonTerminatedPods  []NonTerminatedPod     `json:"nonTerminatedPods"`
	AllocatedResources AllocatedResources     `json:"allocatedResources"`
}
