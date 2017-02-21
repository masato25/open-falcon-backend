package queue

import (
	"github.com/Cepave/open-falcon-backend/modules/transfer/proc"
	cmodel "github.com/open-falcon/common/model"
)

// 将原始数据入到flentd发送缓存队列
func (this *Queues) Push2FluentdSendQueue(items []*cmodel.MetaData) {
	for _, item := range items {
		isSuccess := this.FluentdQueue.PushFront(item)
		if !isSuccess {
			proc.SendToFluentdDropCnt.Incr()
		}
	}
}
