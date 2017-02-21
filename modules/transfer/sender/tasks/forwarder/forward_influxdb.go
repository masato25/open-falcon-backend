package forwarder

import (
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/Cepave/open-falcon-backend/modules/transfer/g"
	"github.com/Cepave/open-falcon-backend/modules/transfer/proc"
	"github.com/Cepave/open-falcon-backend/modules/transfer/sender/pools"
	cmodel "github.com/open-falcon/common/model"
	nsema "github.com/toolkits/concurrent/semaphore"
	"github.com/toolkits/container/list"
)

// Influxdb schedule
func Forward2InfluxdbTask(Q *list.SafeListLimited, concurrent int) {
	cfg := g.Config().Influxdb
	batch := cfg.Batch // 一次发送,最多batch条数据
	conn, err := pools.InfuxdbParseDSN(cfg.Address)
	if err != nil {
		log.Print("syntax of influxdb address is wrong")
		return
	}
	addr := conn.Address

	sema := nsema.NewSemaphore(concurrent)

	for {
		items := Q.PopBackBy(batch)
		count := len(items)
		if count == 0 {
			time.Sleep(DefaultSendTaskSleepInterval)
			continue
		}

		influxdbItems := make([]*cmodel.JudgeItem, count)
		for i := 0; i < count; i++ {
			influxdbItems[i] = items[i].(*cmodel.JudgeItem)
		}

		//	同步Call + 有限并发 进行发送
		sema.Acquire()
		go func(addr string, influxdbItems []*cmodel.JudgeItem, count int) {
			defer sema.Release()

			var err error
			sendOk := false
			for i := 0; i < 3; i++ { //最多重试3次
				err = InfluxdbConnPools.Call(addr, influxdbItems)
				if err == nil {
					sendOk = true
					break
				}
				time.Sleep(time.Millisecond * 10)
			}

			// statistics
			if !sendOk {
				log.Printf("send influxdb %s fail: %v", addr, err)
				proc.SendToInfluxdbFailCnt.IncrBy(int64(count))
			} else {
				proc.SendToInfluxdbCnt.IncrBy(int64(count))
			}
		}(addr, influxdbItems, count)
	}
}
