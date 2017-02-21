package sender

import (
	"strings"
	"time"

	"github.com/Cepave/open-falcon-backend/modules/transfer/proc"
	"github.com/Cepave/open-falcon-backend/modules/transfer/queue"
	"github.com/Cepave/open-falcon-backend/modules/transfer/sender/pools"
	log "github.com/Sirupsen/logrus"
	"github.com/toolkits/container/list"
)

const (
	DefaultProcCronPeriod = time.Duration(5) * time.Second    //ProcCron的周期,默认1s
	DefaultLogCronPeriod  = time.Duration(3600) * time.Second //LogCron的周期,默认300s
)

// send_cron程序入口
func startSenderCron() {
	go startProcCron()
	go startLogCron()
}

func startProcCron() {
	for {
		time.Sleep(DefaultProcCronPeriod)
		refreshSendingCacheSize()
	}
}

func startLogCron() {
	for {
		time.Sleep(DefaultLogCronPeriod)
		logConnPoolsProc()
	}
}

func refreshSendingCacheSize() {
	cqueue := queue.GetQ()
	proc.JudgeQueuesCnt.SetCnt(calcSendCacheSize(cqueue.JudgeQueues))
	proc.GraphQueuesCnt.SetCnt(calcSendCacheSize(cqueue.GraphQueues))
	proc.InfluxdbQueuesCnt.SetCnt(calcSendCacheSize(cqueue.InfluxdbQueues))
}
func calcSendCacheSize(mapList map[string]*list.SafeListLimited) int64 {
	var cnt int64 = 0
	for _, list := range mapList {
		if list != nil {
			cnt += int64(list.Len())
		}
	}
	return cnt
}

func logConnPoolsProc() {
	pool := pools.GetPools()
	log.Printf("connPools proc: \n%v", strings.Join(pool.GraphConnPools.Proc(), "\n"))
}
