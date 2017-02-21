package pools

import (
	"github.com/Cepave/open-falcon-backend/modules/transfer/g"
	cpool "github.com/Cepave/open-falcon-backend/modules/transfer/sender/conn_pool"
	nset "github.com/toolkits/container/set"
)

var pools = Pools{}

type Pools struct {
	JudgeConnPools   *cpool.SafeRpcConnPools
	GraphConnPools   *cpool.SafeRpcConnPools
	FluentdConnPools *cpool.FluentdConnPools
}

func GetPools() Pools {
	return pools
}

func StartPools() {
	cfg := g.Config()

	if cfg.Judge.Enabled {
		//init judge pools
		judgeInstances := nset.NewStringSet()
		for _, instance := range cfg.Judge.Cluster {
			judgeInstances.Add(instance)
		}
		pools.JudgeConnPools = cpool.CreateSafeRpcConnPools(cfg.Judge.MaxConns, cfg.Judge.MaxIdle,
			cfg.Judge.ConnTimeout, cfg.Judge.CallTimeout, judgeInstances.ToSlice())
	}

	if cfg.Graph.Enabled {
		// init graph's pool
		graphInstances := nset.NewSafeSet()
		for _, nitem := range cfg.Graph.ClusterList {
			for _, addr := range nitem.Addrs {
				graphInstances.Add(addr)
			}
		}
	}

	if cfg.Fluentd.Enabled {
		pools.FluentdConnPools = cpool.CreateFluentdConnPools(cfg.Fluentd.Address, cfg.Fluentd.ConnTimeout)
	}
}

func (this *Pools) DestroyConnPools() {
	this.JudgeConnPools.Destroy()
	this.GraphConnPools.Destroy()
	// TsdbConnPoolHelper.Destroy()
	// InfluxdbConnPools.Destroy()
	// StagingConnPoolHelper.Destroy()
}
