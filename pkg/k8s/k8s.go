package k8s

var (
	_      ICoreV1Store = (*coreV1)(nil)
	client ICoreV1Store
)

func NewCoreV1Store(clusterSet IClusterSet) {

	client = &coreV1{clusterSet}
}

func GetClient() ICoreV1Store {
	return client
}

type ICoreV1Store interface {
	IPod
	IDeployment
	IService
	IStatefulSet
	IJob
	IEvent
	INode
	INamespace
}

type coreV1 struct {
	clusterSets IClusterSet
}

func (c *coreV1) Pods(clusterName string) IPodAction {
	clients := c.clusterSets.Get(clusterName)
	return newPods(clients.K8sClient, clients.MetricsClient, clients.InformerClient)
}

func (c *coreV1) Deployments(clusterName string) IDeploymentAction {
	clients := c.clusterSets.Get(clusterName)
	return newDeployments(clients.K8sClient, clients.InformerClient)
}

func (c *coreV1) Services(clusterName string) IServiceAction {
	clients := c.clusterSets.Get(clusterName)
	return newServices(clients.K8sClient, clients.InformerClient)
}

func (c *coreV1) StatefulSets(clusterName string) IStatefulSetAction {
	clients := c.clusterSets.Get(clusterName)
	return newStatefulSets(clients.K8sClient, clients.InformerClient)
}

func (c *coreV1) Jobs(clusterName string) IJobAction {
	clients := c.clusterSets.Get(clusterName)
	return newJobs(clients.K8sClient, clients.InformerClient)
}

func (c *coreV1) Events(clusterName string) IEventAction {
	clients := c.clusterSets.Get(clusterName)
	return newEvents(clients.K8sClient)
}

func (c *coreV1) Nodes(clusterName string) INodeAction {
	clients := c.clusterSets.Get(clusterName)
	return newNodes(clients.K8sClient, clients.InformerClient)
}

func (c *coreV1) Namespaces(clusterName string) INamespaceAction {
	clients := c.clusterSets.Get(clusterName)
	return newNamespaces(clients.K8sClient, clients.InformerClient)
}
