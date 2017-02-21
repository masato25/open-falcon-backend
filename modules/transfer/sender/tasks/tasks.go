package tasks

import (
	"time"

	"github.com/Cepave/open-falcon-backend/modules/transfer/g"
	"github.com/Cepave/open-falcon-backend/modules/transfer/queue"
	"github.com/Cepave/open-falcon-backend/modules/transfer/sender/tasks/forwarder"
)

const (
	DefaultSendTaskSleepInterval = time.Millisecond * 50 //默认睡眠间隔为50ms
)

func StartTasks() {
	forwarder.Start()
	cfg := g.Config()
	// init semaphore
	judgeConcurrent := cfg.Judge.MaxConns
	graphConcurrent := cfg.Graph.MaxConns
	tsdbConcurrent := cfg.Tsdb.MaxConns
	influxdbConcurrent := cfg.Influxdb.MaxIdle
	cqueue := queue.GetQ()

	if tsdbConcurrent < 1 {
		tsdbConcurrent = 1
	}

	if judgeConcurrent < 1 {
		judgeConcurrent = 1
	}

	if graphConcurrent < 1 {
		graphConcurrent = 1
	}
	if influxdbConcurrent < 1 {
		influxdbConcurrent = 1
	}

	if cfg.Judge.Enabled {
		// init send go-routines
		for node, _ := range cfg.Judge.Cluster {
			judgeQueue := cqueue.JudgeQueues[node]
			go forwarder.Forward2JudgeTask(judgeQueue, node, judgeConcurrent)
		}
	}

	if cfg.Graph.Enabled {
		for node, nitem := range cfg.Graph.ClusterList {
			for _, addr := range nitem.Addrs {
				graphQueue := cqueue.GraphQueues[node+addr]
				go forwarder.Forward2GraphTask(graphQueue, node, addr, graphConcurrent)
			}
		}
	}

	if cfg.Tsdb.Enabled {
		go forwarder.Forward2TsdbTask(cqueue.TsdbQueue, tsdbConcurrent)
	}

	if cfg.Influxdb.Enabled {
		go forwarder.Forward2InfluxdbTask(cqueue.InfluxdbQueues["default"], influxdbConcurrent)
	}

	if cfg.NqmRest.Enabled {
		go forwarder.Forward2NqmTask(cqueue.NqmIcmpQueue, g.Config().NqmRest.Fping)
		go forwarder.Forward2NqmTask(cqueue.NqmTcpQueue, g.Config().NqmRest.Tcpping)
		go forwarder.Forward2NqmTask(cqueue.NqmTcpconnQueue, g.Config().NqmRest.Tcpconn)
	}

	if cfg.Staging.Enabled {
		go forwarder.Forward2StagingTask(cqueue.StagingQueue)
	}

	if cfg.Fluentd.Enabled {
		go forwarder.Forward2FluentdTask(cqueue.FluentdQueue, cfg.Fluentd.MaxConns)
	}
}
