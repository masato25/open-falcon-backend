package queue

import (
	"github.com/Cepave/open-falcon-backend/modules/transfer/g"
	"github.com/Cepave/open-falcon-backend/modules/transfer/rings"
	nlist "github.com/toolkits/container/list"
)

const (
	DefaultSendQueueMaxSize = 102400 //10.24w
)

var cqueue = Queues{}
var MinStep = 30

var QueueTypes = []string{
	"Judge",
	"Graph",
	"Influxdb",
	"Tsdb",
	"NqmIcmp",
	"NqmTcp",
	"NqmTcpconn",
	"Staging",
	"Fluentd",
}

type Queues struct {
	Rings           rings.Rings
	JudgeQueues     map[string]*nlist.SafeListLimited
	GraphQueues     map[string]*nlist.SafeListLimited
	InfluxdbQueues  map[string]*nlist.SafeListLimited
	TsdbQueue       *nlist.SafeListLimited
	NqmIcmpQueue    *nlist.SafeListLimited
	NqmTcpQueue     *nlist.SafeListLimited
	NqmTcpconnQueue *nlist.SafeListLimited
	StagingQueue    *nlist.SafeListLimited
	FluentdQueue    *nlist.SafeListLimited
}

func GetQ() Queues {
	return cqueue
}

func Start() {
	cfg := g.Config()
	if MinStep < 1 {
		MinStep = 30 //默认30s
	}
	cqueue.Rings = rings.GetRings()
	if cfg.Judge.Enabled {
		for node, _ := range cfg.Judge.Cluster {
			Q := nlist.NewSafeListLimited(DefaultSendQueueMaxSize)
			cqueue.JudgeQueues[node] = Q
		}
	}

	if cfg.Graph.Enabled {
		for node, nitem := range cfg.Graph.ClusterList {
			for _, addr := range nitem.Addrs {
				Q := nlist.NewSafeListLimited(DefaultSendQueueMaxSize)
				cqueue.GraphQueues[node+addr] = Q
			}
		}
	}

	if cfg.Tsdb.Enabled {
		cqueue.TsdbQueue = nlist.NewSafeListLimited(DefaultSendQueueMaxSize)
	}

	if cfg.Influxdb.Enabled {
		Q := nlist.NewSafeListLimited(DefaultSendQueueMaxSize)
		cqueue.InfluxdbQueues["default"] = Q
	}

	if cfg.NqmRest.Enabled {
		cqueue.NqmIcmpQueue = nlist.NewSafeListLimited(DefaultSendQueueMaxSize)
		cqueue.NqmTcpQueue = nlist.NewSafeListLimited(DefaultSendQueueMaxSize)
		cqueue.NqmTcpconnQueue = nlist.NewSafeListLimited(DefaultSendQueueMaxSize)
	}

	if cfg.Staging.Enabled {
		cqueue.StagingQueue = nlist.NewSafeListLimited(DefaultSendQueueMaxSize)
	}

	if cfg.Fluentd.Enabled {
		cqueue.FluentdQueue = nlist.NewSafeListLimited(DefaultSendQueueMaxSize)
	}
}
