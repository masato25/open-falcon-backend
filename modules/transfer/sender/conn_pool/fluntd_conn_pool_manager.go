package conn_pool

import (
	"fmt"
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/Cepave/open-falcon-backend/modules/transfer/proc"
	log "github.com/Sirupsen/logrus"
)

type FluentdClient struct {
	sync.RWMutex
	conn      net.Conn
	connected bool
	cli       interface{}
	addr      string
	timeout   int
}

// ConnPools Manager
type FluentdConnPools struct {
	FluentdClients []*FluentdClient
	Size           int
}

func (this *FluentdClient) Addr() string {
	return this.addr
}

func (this *FluentdClient) Closed() bool {
	return this.connected == false
}

func (this *FluentdClient) Close() error {
	err := this.conn.Close()
	if err != nil {
		log.Error(err.Error())
	} else {
		this.connected = false
		this.cli = nil
	}
	return err
}

func (this *FluentdClient) GetConn() net.Conn {
	return this.conn
}

func (this *FluentdClient) ReConn() bool {
	err := this.Connect2(this.addr, this.timeout)
	if err != nil {
		log.Errorf("trying to reconnect to %v, got error: %v", this.Addr(), err.Error())
		return false
	} else {
		log.Infof("connect to %v, scuessfuly", this.Addr())
		this.connected = true
		return true
	}
}

func (this *FluentdClient) Connect2(addr string, timeout int) error {
	this.addr = addr
	this.timeout = timeout

	conn, err := net.DialTimeout("tcp", addr, (time.Duration(timeout) * time.Second))
	if err != nil {
		return err
	} else {
		this.connected = true
	}
	this.conn = conn
	return err
}

func sendData(fc *FluentdClient, data string) bool {
	if fc.Closed() {
		fc.ReConn()
		return false
	}
	conn := fc.GetConn()
	n, err := fmt.Fprintf(conn, data+"\n")
	if err != nil {
		fc.ReConn()
		log.Errorf("send itme to fluentd got error: %v", err.Error())
		return false
	} else {
		log.Debugf("%v item is send to flunetd:%v", n, fc.Addr())
	}
	return true
}

func (this *FluentdConnPools) Call(data string, count int64) {
	conn := this.Get()
	ok := sendData(conn, data+"\n")
	if !ok {
		rand.Seed(time.Now().UnixNano())
		for _, i := range rand.Perm(len(this.FluentdClients)) {
			conn = this.Get(i)
			if sendOk := sendData(conn, data+"\n"); sendOk {
				proc.SendToFluentdCnt.IncrBy(count)
				return
			}
		}
		log.Error("all connection is gone, item lost.")
		proc.SendToFluentdDropCnt.IncrBy(count)
		return
	}
	proc.SendToFluentdCnt.IncrBy(count)
	return
}

func (this *FluentdConnPools) Get(inds ...int) *FluentdClient {
	log.Debugf("FluentdConnPools: %v\n", this)
	log.Debugf("fpools: %v\n", this.Size)
	if inds == nil {
		ind := rand.Intn(this.Size)
		return this.FluentdClients[ind]
	} else {
		log.Debugf("inds: %v", inds)
		ind := inds[0]
		return this.FluentdClients[ind]
	}
}

func CreateFluentdConnPools(addrs []string, connTimeout int) *FluentdConnPools {
	pools := FluentdConnPools{}
	counter := 0
	for _, addr := range addrs {
		f := FluentdClient{}
		err := f.Connect2(addr, connTimeout)
		if err != nil {
			log.Error(err.Error())
		}
		pools.FluentdClients = append(pools.FluentdClients, &f)
		counter += 1
	}
	pools.Size = counter
	log.Debugf("CreateFluentdConnPools: %v", pools)
	return &pools
}

func (this *FluentdConnPools) Destroy() {
	this.FluentdClients = []*FluentdClient{}
}
