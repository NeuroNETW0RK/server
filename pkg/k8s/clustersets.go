package k8s

import (
	"sync"
)

var clusterSets IClusterSet

func GetClusterSets() IClusterSet {
	return clusterSets
}

type IClusterSet interface {
	Add(clusterName string, clientSets *ClientSets)
	Update(clusterName string, clientSets *ClientSets)
	Delete(clusterName string)
	Get(clusterName string) *ClientSets
	List() map[string]*ClientSets
}

var _ IClusterSet = (*ClusterSet)(nil)

func NewClusterSet() {
	clusterSets = &ClusterSet{clientSets: map[string]*ClientSets{}}
}

type ClusterSet struct {
	lock       sync.RWMutex
	clientSets map[string]*ClientSets
}

func (c *ClusterSet) Add(clusterName string, clientSets *ClientSets) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.clientSets[clusterName] = clientSets
}

func (c *ClusterSet) Update(clusterName string, clientSets *ClientSets) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.clientSets[clusterName] = clientSets
}

func (c *ClusterSet) Delete(clusterName string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if _, ok := c.clientSets[clusterName]; ok {
		delete(c.clientSets, clusterName)
	}
}

func (c *ClusterSet) Get(clusterName string) *ClientSets {
	c.lock.Lock()
	defer c.lock.Unlock()

	item, exists := c.clientSets[clusterName]
	if !exists {
		return &ClientSets{
			K8sClient:      nil,
			MetricsClient:  nil,
			InformerClient: nil,
		}
	}
	return item
}

func (c *ClusterSet) List() map[string]*ClientSets {
	c.lock.Lock()
	defer c.lock.Unlock()

	return c.clientSets
}
