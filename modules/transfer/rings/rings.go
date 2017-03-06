package rings

import (
	"github.com/Cepave/consistent"
	"github.com/Cepave/open-falcon-backend/modules/transfer/g"
)

type Rings struct {
	JudgeNodeRing *ConsistentHashNodeRing
	GraphNodeRing *ConsistentHashNodeRing
}

var rings = Rings{}

func GetRings() Rings {
	return rings
}

func Start() {
	cfg := g.Config()
	rings.JudgeNodeRing = newConsistentHashNodesRing(cfg.Judge.Replicas, KeysOfMap(cfg.Judge.Cluster))
	rings.GraphNodeRing = newConsistentHashNodesRing(cfg.Graph.Replicas, KeysOfMap(cfg.Graph.Cluster))
}

// TODO 考虑放到公共组件库,或utils库
func KeysOfMap(m map[string]string) []string {
	keys := make([]string, len(m))
	i := 0
	for key, _ := range m {
		keys[i] = key
		i++
	}

	return keys
}

// 一致性哈希环,用于管理服务器节点.
type ConsistentHashNodeRing struct {
	ring *consistent.Consistent
}

func newConsistentHashNodesRing(numberOfReplicas int, nodes []string) *ConsistentHashNodeRing {
	ret := &ConsistentHashNodeRing{ring: consistent.New()}
	ret.SetNumberOfReplicas(numberOfReplicas)
	ret.SetNodes(nodes)
	return ret
}

// 根据pk,获取node节点. chash(pk) -> node
func (this *ConsistentHashNodeRing) GetNode(pk string) (string, error) {
	return this.ring.Get(pk)
}

func (this *ConsistentHashNodeRing) SetNodes(nodes []string) {
	for _, node := range nodes {
		this.ring.Add(node)
	}
}

func (this *ConsistentHashNodeRing) SetNumberOfReplicas(num int) {
	this.ring.NumberOfReplicas = num
}