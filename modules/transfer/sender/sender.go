package sender

import (
	"log"

	"github.com/Cepave/open-falcon-backend/modules/transfer/sender/pools"
	"github.com/Cepave/open-falcon-backend/modules/transfer/sender/tasks"
)

// 初始化数据发送服务, 在main函数中调用
func Start() {
	pools.StartPools()
	tasks.StartTasks()
	startSenderCron()
	log.Println("send.Start, ok")
}
