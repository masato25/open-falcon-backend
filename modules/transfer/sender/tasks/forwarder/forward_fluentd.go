package forwarder

import (
	"encoding/json"
	"time"

	"github.com/Cepave/open-falcon-backend/modules/transfer/g"
	cmodel "github.com/open-falcon/common/model"
	nsema "github.com/toolkits/concurrent/semaphore"
	nlist "github.com/toolkits/container/list"
)

// Tsdb定时任务, 将数据通过api发送到tsdb
func Forward2FluentdTask(fqueue *nlist.SafeListLimited, concurrent int) {
	conf := g.Config()
	batch := conf.Fluentd.Batch // 一次发送,最多batch条数据
	sema := nsema.NewSemaphore(concurrent)
	// conn, err := net.Dial("tcp", conf.Fluentd.Address)
	// defer fmt.Println("Forward2FluentdTask is closed")
	// if err != nil {
	// fmt.Println(err.Error())
	// return
	// }
	// defer conn.Close()
	for {
		items := fqueue.PopBackBy(batch)
		itmeSize := len(items)
		if len(items) == 0 {
			time.Sleep(DefaultSendTaskSleepInterval)
			continue
		}
		//  同步Call + 有限并发 进行发送
		sema.Acquire()
		go func(itemList []interface{}) {
			defer sema.Release()
			thisItems := []interface{}{}
			for _, item := range itemList {
				citem := item.(*cmodel.MetaData)
				thisItems = append(thisItems,
					[]interface{}{
						citem.Timestamp,
						map[string]interface{}{
							"message": convert2JsonItem(citem),
						},
					})
			}
			itemstr := MapStructToString([]interface{}{
				"owl.data",
				thisItems,
				map[string]interface{}{
					"option": "optional",
				},
			})
			FluentdCoonPools.Call(itemstr, int64(itmeSize))
		}(items)
	}
}

func MapStructToString(dt []interface{}) string {
	dd, _ := json.Marshal(dt)
	data := string(dd)
	return data
}

// 转化为json string格式
func convert2JsonItem(d *cmodel.MetaData) (data map[string]interface{}) {
	data = map[string]interface{}{
		"endpoint": d.Endpoint,
		"metric":   d.Metric,
		"type":     d.CounterType,
		"tags":     d.Tags,
		"step":     d.Step,
		"time":     d.Timestamp,
		"value":    d.Value,
	}
	return
}
