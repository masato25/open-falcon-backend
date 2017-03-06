package queue

import (
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/Cepave/open-falcon-backend/modules/transfer/proc"
	cmodel "github.com/open-falcon/common/model"
)

// Push metrics from fping to the queue for RESTful API
func (this *Queues) Push2NqmIcmpSendQueue(pingItems []*cmodel.MetaData) {
	for _, item := range pingItems {
		item, err := convert2NqmPingItem(item)
		if err != nil {
			log.Println("NqmPing converting error:", err)
			continue
		}
		isSuccess := this.NqmIcmpQueue.PushFront(item)

		if !isSuccess {
			proc.SendToNqmIcmpDropCnt.Incr()
		}
	}
}

// Push metrics from tcpping to the queue for RESTful API
func (this *Queues) Push2NqmTcpSendQueue(pingItems []*cmodel.MetaData) {
	for _, item := range pingItems {
		item, err := convert2NqmPingItem(item)
		if err != nil {
			log.Println("NqmPing converting error:", err)
			continue
		}
		isSuccess := this.NqmTcpQueue.PushFront(item)

		if !isSuccess {
			proc.SendToNqmTcpDropCnt.Incr()
		}
	}
}

// Push metrics from tcpconn to the queue for RESTful API
func (this *Queues) Push2NqmTcpconnSendQueue(connItems []*cmodel.MetaData) {
	for _, item := range connItems {
		nqmitem, err := convert2NqmConnItem(item)
		if err != nil {
			log.Println("NqmConn converting error:", err)
			continue
		}
		isSuccess := this.NqmTcpconnQueue.PushFront(nqmitem)

		if !isSuccess {
			proc.SendToNqmTcpconnDropCnt.Incr()
		}
	}
}

func Demultiplex(items []*cmodel.MetaData) ([]*cmodel.MetaData, []*cmodel.MetaData, []*cmodel.MetaData, []*cmodel.MetaData) {
	nqmFpings := []*cmodel.MetaData{}
	nqmTcppings := []*cmodel.MetaData{}
	nqmTcpconns := []*cmodel.MetaData{}
	generics := []*cmodel.MetaData{}

	for _, item := range items {
		switch item.Metric {
		case "nqm-fping":
			nqmFpings = append(nqmFpings, item)
		case "nqm-tcpping":
			nqmTcppings = append(nqmTcppings, item)
		case "nqm-tcpconn":
			nqmTcpconns = append(nqmTcpconns, item)
		default:
			generics = append(generics, item)
		}
	}

	return nqmFpings, nqmTcppings, nqmTcpconns, generics
}

func convert2NqmPingItem(d *cmodel.MetaData) (*nqmPingItem, error) {
	var t nqmPingItem
	agent, err := convert2NqmEndpoint(d, "agent")
	if err != nil {
		return &t, err
	}
	target, err := convert2NqmEndpoint(d, "target")
	if err != nil {
		return &t, err
	}
	metrics, err := Convert2NqmMetrics(d)
	if err != nil {
		return &t, err
	}

	t = nqmPingItem{
		Timestamp: d.Timestamp,
		Agent:     *agent,
		Target:    *target,
		Metrics:   *metrics,
	}

	return &t, nil
}

func convert2NqmConnItem(d *cmodel.MetaData) (*nqmConnItem, error) {
	var t nqmConnItem
	var tt float32
	t.Timestamp = d.Timestamp
	agent, err := convert2NqmEndpoint(d, "agent")
	if err != nil {
		return &t, err
	}
	target, err := convert2NqmEndpoint(d, "target")
	if err != nil {
		return &t, err
	}
	if err := strToFloat32(&tt, "time", d.Tags); err != nil {
		return nil, err
	}
	t = nqmConnItem{
		Timestamp: d.Timestamp,
		Agent:     *agent,
		Target:    *target,
		TotalTime: tt,
	}

	return &t, nil
}

func strToFloat32(out *float32, index string, dict map[string]string) error {
	var err error
	var ff float64
	if v, ok := dict[index]; ok {
		ff, err = strconv.ParseFloat(v, 32)
		if err != nil {
			return err
		}
		*out = float32(ff)
	}
	return nil
}

func strToInt32(out *int32, index string, dict map[string]string) error {
	var err error
	var ii int64
	if v, ok := dict[index]; ok {
		ii, err = strconv.ParseInt(v, 10, 32)
		if err != nil {
			return err
		}
		*out = int32(ii)
	}
	return nil
}

func strToInt32Slc(out *[]int32, index string, dict map[string]string) error {
	if v, ok := dict[index]; ok {
		if v == "" {
			return nil
		}
		strSlc := strings.Split(v, "-")
		for _, str := range strSlc {
			i, err := strconv.Atoi(str)
			if err != nil {
				return err
			}
			*out = append(*out, int32(i))
		}
	}
	return nil
}

func strToInt16(out *int16, index string, dict map[string]string) error {
	var err error
	var ii int64
	if v, ok := dict[index]; ok {
		ii, err = strconv.ParseInt(v, 10, 16)
		if err != nil {
			return err
		}
		*out = int16(ii)
	}
	return nil
}

func convert2NqmEndpoint(d *cmodel.MetaData, endType string) (*nqmEndpoint, error) {
	t := nqmEndpoint{
		Id:          -1,
		IspId:       -1,
		ProvinceId:  -1,
		CityId:      -1,
		NameTagId:   -1,
		GroupTagIds: []int32{},
	}

	if err := strToInt32(&t.Id, endType+"-id", d.Tags); err != nil {
		return nil, err
	}
	if err := strToInt16(&t.IspId, endType+"-isp-id", d.Tags); err != nil {
		return nil, err
	}
	if err := strToInt16(&t.ProvinceId, endType+"-province-id", d.Tags); err != nil {
		return nil, err
	}
	if err := strToInt16(&t.CityId, endType+"-city-id", d.Tags); err != nil {
		return nil, err
	}
	if err := strToInt32(&t.NameTagId, endType+"-name-tag-id", d.Tags); err != nil {
		return nil, err
	}
	if err := strToInt32Slc(&t.GroupTagIds, endType+"-group-tag-ids", d.Tags); err != nil {
		return nil, err
	}

	return &t, nil
}

// 轉化成 nqmMetrc 格式
func Convert2NqmMetrics(d *cmodel.MetaData) (*nqmMetrics, error) {
	t := nqmMetrics{
		Rttmin:      -1,
		Rttavg:      -1,
		Rttmax:      -1,
		Rttmdev:     -1,
		Rttmedian:   -1,
		Pkttransmit: -1,
		Pktreceive:  -1,
	}
	var ff float32
	if err := strToFloat32(&ff, "rttmin", d.Tags); err != nil {
		return nil, err
	}
	t.Rttmin = int32(ff)
	if err := strToFloat32(&ff, "rttmax", d.Tags); err != nil {
		return nil, err
	}
	t.Rttmax = int32(ff)

	if err := strToFloat32(&t.Rttavg, "rttavg", d.Tags); err != nil {
		return nil, err
	}
	if err := strToFloat32(&t.Rttmdev, "rttmdev", d.Tags); err != nil {
		return nil, err
	}
	if err := strToFloat32(&t.Rttmedian, "rttmedian", d.Tags); err != nil {
		return nil, err
	}
	if err := strToInt32(&t.Pkttransmit, "pkttransmit", d.Tags); err != nil {
		return nil, err
	}
	if err := strToInt32(&t.Pktreceive, "pktreceive", d.Tags); err != nil {
		return nil, err
	}

	return &t, nil
}
