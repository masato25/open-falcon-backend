package forwarder

import (
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/Cepave/open-falcon-backend/modules/transfer/g"
	"github.com/Cepave/open-falcon-backend/modules/transfer/proc"
	cmodel "github.com/open-falcon/common/model"
	nsema "github.com/toolkits/concurrent/semaphore"
	nlist "github.com/toolkits/container/list"
)

func Forward2StagingTask(stagingQueue *nlist.SafeListLimited) {
	batch := g.Config().Staging.Batch
	retry := g.Config().Staging.MaxRetry
	concurrent := g.Config().Staging.MaxConns
	sema := nsema.NewSemaphore(concurrent)

	for {
		items := stagingQueue.PopBackBy(batch)
		count := len(items)
		if count == 0 {
			time.Sleep(DefaultSendTaskSleepInterval)
			continue
		}

		stagingItems := make([]*cmodel.MetricValue, count)
		for i := 0; i < count; i++ {
			stagingItems[i] = items[i].(*cmodel.MetricValue)
		}

		//	A synchronous call with limited concurrence
		sema.Acquire()
		go func(stagingItems []*cmodel.MetricValue, count int) {
			defer sema.Release()

			resp := &cmodel.SimpleRpcResponse{}
			var err error
			sendOk := false
			for i := 0; i < retry; i++ {
				err = StagingConnPoolHelper.Call("Transfer.Update", stagingItems, resp)
				if err == nil {
					sendOk = true
					break
				}
				time.Sleep(time.Millisecond * 10)
			}

			// statistics
			if !sendOk {
				log.Printf("send staging fail: %v", err)
				proc.SendToStagingFailCnt.IncrBy(int64(count))
			} else {
				proc.SendToStagingCnt.IncrBy(int64(count))
			}
		}(stagingItems, count)
	}
}
