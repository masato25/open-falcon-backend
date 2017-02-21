package socket

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/Cepave/open-falcon-backend/modules/transfer/g"
	"github.com/Cepave/open-falcon-backend/modules/transfer/proc"
	"github.com/Cepave/open-falcon-backend/modules/transfer/queue"
	cmodel "github.com/open-falcon/common/model"
)

func socketTelnetHandle(conn net.Conn) {
	defer conn.Close()

	items := []*cmodel.MetaData{}
	buf := bufio.NewReader(conn)

	cfg := g.Config()
	timeout := time.Duration(cfg.Socket.Timeout) * time.Second
	que := queue.GetQ()
	for {
		conn.SetReadDeadline(time.Now().Add(timeout))
		line, err := buf.ReadString('\n')
		if err != nil {
			break
		}

		line = strings.Trim(line, "\n")

		if line == "quit" {
			break
		}

		if line == "" {
			continue
		}

		t := strings.Fields(line)
		if len(t) < 2 {
			continue
		}

		cmd := t[0]

		if cmd != "update" {
			continue
		}

		item, err := convertLine2MetaData(t[1:])
		if err != nil {
			continue
		}

		items = append(items, item)
	}

	// statistics
	proc.SocketRecvCnt.IncrBy(int64(len(items)))
	proc.RecvCnt.IncrBy(int64(len(items)))

	// demultiplexing
	nqmFpingItems, nqmTcppingItems, nqmTcpconnItems, genericItems := queue.Demultiplex(items)

	if cfg.Graph.Enabled {
		que.Push2GraphSendQueue(genericItems)
	}

	if cfg.Judge.Enabled {
		que.Push2JudgeSendQueue(genericItems)
	}

	if cfg.Influxdb.Enabled {
		que.Push2InfluxdbSendQueue(genericItems)
	}

	if cfg.NqmRest.Enabled {
		que.Push2NqmIcmpSendQueue(nqmFpingItems)
		que.Push2NqmTcpSendQueue(nqmTcppingItems)
		que.Push2NqmTcpconnSendQueue(nqmTcpconnItems)
	}

	return

}

// example: endpoint counter timestamp value [type] [step]
// default type is DERIVE, default step is 60s
func convertLine2MetaData(fields []string) (item *cmodel.MetaData, err error) {
	if len(fields) != 4 && len(fields) != 5 && len(fields) != 6 {
		err = fmt.Errorf("not_enough_fileds")
		return
	}

	endpoint, metric := fields[0], fields[1]
	ts, err := strconv.ParseInt(fields[2], 10, 64)
	if err != nil {
		return
	}

	v, err := strconv.ParseFloat(fields[3], 64)
	if err != nil {
		return
	}

	type_ := g.COUNTER
	if len(fields) >= 5 {
		type_ = fields[4]
	}

	if type_ != g.DERIVE && type_ != g.GAUGE && type_ != g.COUNTER {
		err = fmt.Errorf("invalid_counter_type")
		return
	}

	var step int64 = g.DEFAULT_STEP
	if len(fields) == 6 {
		dst_args := strings.Split(fields[5], ":")
		if len(dst_args) == 1 {
			step, err = strconv.ParseInt(dst_args[0], 10, 64)
			if err != nil {
				return
			}
		} else if len(dst_args) == 4 {
			// for backend-compatible
			// heartbeat:min:max:step
			step, err = strconv.ParseInt(dst_args[3], 10, 64)
			if err != nil {
				return
			}
		} else {
			err = fmt.Errorf("invalid_counter_step")
			return
		}
	}

	item = &cmodel.MetaData{
		Metric:      metric,
		Endpoint:    endpoint,
		Timestamp:   ts,
		Step:        step,
		Value:       v,
		CounterType: type_,
		Tags:        make(map[string]string),
	}

	return item, nil
}
