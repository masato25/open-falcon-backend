package pools

import (
	"log"

	"github.com/Cepave/open-falcon-backend/modules/transfer/g"
	cpool "github.com/Cepave/open-falcon-backend/modules/transfer/sender/conn_pool"
	nset "github.com/toolkits/container/set"
)

var pools = Pools{}

type Pools struct {
	JudgeConnPools        *cpool.SafeRpcConnPools
	GraphConnPools        *cpool.SafeRpcConnPools
	FluentdConnPools      *cpool.FluentdConnPools
	InfluxdbConnPools     *cpool.InfluxdbConnPools
	StagingConnPoolHelper *cpool.StagingConnPoolHelper
	TsdbConnPoolHelper    *cpool.TsdbConnPoolHelper
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
		pools.GraphConnPools = cpool.CreateSafeRpcConnPools(cfg.Graph.MaxConns, cfg.Graph.MaxIdle,
			cfg.Graph.ConnTimeout, cfg.Graph.CallTimeout, graphInstances.ToSlice())
	}

	if cfg.Tsdb.Enabled {
		pools.TsdbConnPoolHelper = cpool.NewTsdbConnPoolHelper(cfg.Tsdb.Address, cfg.Tsdb.MaxConns, cfg.Tsdb.MaxIdle, cfg.Tsdb.ConnTimeout, cfg.Tsdb.CallTimeout)
	}

	if cfg.Influxdb.Enabled {
		influxdbInstances := make([]cpool.InfluxdbConnection, 1)
		dsn, err := InfuxdbParseDSN(cfg.Influxdb.Address)
		if err != nil {
			log.Print("syntax of influxdb address is wrong")
		} else {
			influxdbInstances[0] = *dsn
			pools.InfluxdbConnPools = cpool.CreateInfluxdbConnPools(cfg.Influxdb.MaxConns, cfg.Influxdb.MaxIdle,
				cfg.Influxdb.ConnTimeout, cfg.Influxdb.CallTimeout, influxdbInstances)
		}
	}

	if cfg.Staging.Enabled {
		pools.StagingConnPoolHelper = cpool.NewStagingConnPoolHelper(cfg.Staging.Address, cfg.Staging.MaxConns, cfg.Staging.MaxIdle, cfg.Staging.ConnTimeout, cfg.Staging.CallTimeout)
	}

	if cfg.Fluentd.Enabled {
		pools.FluentdConnPools = cpool.CreateFluentdConnPools(cfg.Fluentd.Address, cfg.Fluentd.ConnTimeout)
	}
}

func (this *Pools) DestroyConnPools() {
	this.JudgeConnPools.Destroy()
	this.GraphConnPools.Destroy()
	this.TsdbConnPoolHelper.Destroy()
	this.InfluxdbConnPools.Destroy()
	this.StagingConnPoolHelper.Destroy()
	this.FluentdConnPools.Destroy()
}
