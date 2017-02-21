package queue

import (
	"github.com/Cepave/open-falcon-backend/modules/transfer/common"
	"github.com/Cepave/open-falcon-backend/modules/transfer/proc"
	cmodel "github.com/open-falcon/common/model"
)

// Push data to 3rd-party database
func (this Queues) Push2InfluxdbSendQueue(items []*cmodel.MetaData) {
	for _, item := range items {
		// align ts
		step := int(item.Step)
		if step < MinStep {
			step = MinStep
		}
		ts := common.AlignTs(item.Timestamp, int64(step))

		influxdbItem := &cmodel.JudgeItem{
			Endpoint:  item.Endpoint,
			Metric:    item.Metric,
			Value:     item.Value,
			Timestamp: ts,
			JudgeType: item.CounterType,
			Tags:      item.Tags,
		}
		Q := this.InfluxdbQueues["default"]
		isSuccess := Q.PushFront(influxdbItem)

		// statistics
		if !isSuccess {
			proc.SendToInfluxdbDropCnt.Incr()
		}
	}
}
