package forwarder

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/Cepave/open-falcon-backend/modules/transfer/g"
	"github.com/Cepave/open-falcon-backend/modules/transfer/proc"
	nsema "github.com/toolkits/concurrent/semaphore"
	"github.com/toolkits/container/list"
)

func ForwardNqmItems(nqmItem interface{}, nqmUrl string, s *nsema.Semaphore) {
	defer s.Release()

	jsonItem, jsonErr := json.Marshal(nqmItem)
	if jsonErr != nil {
		log.Errorf("Error on serialization for nqm item(ICMP, TCP, or TCPCONN): %v", jsonErr)
		proc.SendToNqmIcmpFailCnt.IncrBy(1)
		return
	}

	log.Debugf("[ Cassandra ] JSON data to %s: %s", nqmUrl, string(jsonItem))
	postReq, err := http.NewRequest("POST", nqmUrl, bytes.NewBuffer(jsonItem))

	postReq.Header.Set("Content-Type", "application/json; charset=UTF-8")
	postReq.Header.Set("Connection", "close")
	httpClient := &http.Client{}
	postResp, err := httpClient.Do(postReq)
	if err != nil {
		log.Errorln("[ Cassandra ] Error on push:", err)
		proc.SendToNqmIcmpFailCnt.IncrBy(1)
		return
	}
	defer postResp.Body.Close()
	proc.SendToNqmIcmpCnt.IncrBy(1)
}

func Forward2NqmTask(Q *list.SafeListLimited, apiUrl string) {
	batch := g.Config().NqmRest.Batch // 一次发送,最多batch条数据
	concurrent := g.Config().NqmRest.MaxConns
	sema := nsema.NewSemaphore(concurrent)

	for {
		items := Q.PopBackBy(batch)
		if len(items) == 0 {
			time.Sleep(DefaultSendTaskSleepInterval)
			continue
		}

		for _, v := range items {
			sema.Acquire()
			go ForwardNqmItems(v, apiUrl, sema)
		}
	}
}
