package forwarder

import (
	"time"

	cpool "github.com/Cepave/open-falcon-backend/modules/transfer/sender/conn_pool"
	"github.com/Cepave/open-falcon-backend/modules/transfer/sender/pools"
)

const (
	DefaultSendTaskSleepInterval = time.Millisecond * 50 //默认睡眠间隔为50ms
)

var (
	GraphConnPools        *cpool.SafeRpcConnPools
	JudgeConnPools        *cpool.SafeRpcConnPools
	TsdbConnPoolHelper    *cpool.TsdbConnPoolHelper
	InfluxdbConnPools     *cpool.InfluxdbConnPools
	StagingConnPoolHelper *cpool.StagingConnPoolHelper
	FluentdCoonPools      *cpool.FluentdConnPools
)

func Start() {
	mypools := pools.GetPools()
	JudgeConnPools = mypools.JudgeConnPools
	GraphConnPools = mypools.GraphConnPools
	FluentdCoonPools = mypools.FluentdConnPools
}
