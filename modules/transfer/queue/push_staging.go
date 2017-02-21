package queue

import (
	"github.com/Cepave/open-falcon-backend/modules/transfer/proc"
	cmodel "github.com/open-falcon/common/model"
)

// Push data from endpoint in filters to Staging
func (this *Queues) Push2StagingSendQueue(items []*cmodel.MetricValue) {
	for _, item := range items {
		isSuccess := this.StagingQueue.PushFront(item)

		if !isSuccess {
			proc.SendToStagingDropCnt.Incr()
		}
	}
}
