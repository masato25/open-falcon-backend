package conn_pool

import (
	"fmt"
	"math/rand"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/Cepave/open-falcon-backend/modules/transfer/proc"
	log "github.com/Sirupsen/logrus"
)

type FluentdClient struct {
	sync.RWMutex
	conn    net.Conn
	cli     interface{}
	addr    string
	timeout int
}

// ConnPools Manager
type FluentdConnPools struct {
	FluentdClients []*FluentdClient
}

func (this *FluentdClient) Addr() string {
	return this.addr
}
func (this *FluentdClient) Closed() bool {
	return this.cli == nil
}

func (this *FluentdClient) Close() error {
	err := this.conn.Close()
	if err != nil {
		log.Error(err.Error())
	} else {
		this.cli = nil
	}
	return err
}

func (this *FluentdClient) GetConn() net.Conn {
	return this.conn
}

func (this *FluentdClient) ReConn() error {
	return this.Conn(this.addr, this.timeout)
}

func (this *FluentdClient) Conn(addr string, timeout int) error {
	this.addr = addr
	this.timeout = timeout

	conn, err := net.DialTimeout("tcp", addr, (time.Duration(timeout) * time.Second))
	if err != nil {
		return err
	}
	this.conn = conn
	return err
}

func (this *FluentdConnPools) Call(data string) {
	conn, ok := this.Get()
	if !ok {
		log.Error("get connection error")
	}
	conn.Lock()
	n, err := fmt.Fprintf(conn.GetConn(), data+"\n")
	if err != nil && strings.Contains(err.Error(), "write: broken pipe") {
		err = conn.ReConn()
		if err != nil {
			log.Errorf("trying to reconnect to %v, got error: %v", conn, err.Error())
		} else {
			log.Infof("connect to %v, scuessfuly", conn.Addr())
			//retry
			n, err = fmt.Fprintf(conn.GetConn(), data+"\n")
		}
	} else if err != nil {
		log.Error(err.Error())
	} else {
		log.Debugf("%v items is send", n)
	}
	//failed failed count incrase
	if n == 0 || err != nil {
		proc.SendToFluentdFailCnt.Incr()
	}
	conn.Unlock()
}

func (this *FluentdConnPools) Get() (*FluentdClient, bool) {
	log.Debugf("FluentdConnPools: %v\n", this)
	log.Debugf("fpools: %v\n", len(this.FluentdClients))
	if len(this.FluentdClients) == 0 {
		return nil, false
	}
	ind := rand.Intn(len(this.FluentdClients))
	return this.FluentdClients[ind], true
}

func CreateFluentdConnPools(addrs []string, connTimeout int) *FluentdConnPools {
	pools := FluentdConnPools{}
	for _, addr := range addrs {
		f := FluentdClient{}
		err := f.Conn(addr, connTimeout)
		if err != nil {
			log.Error(err.Error())
			continue
		}
		pools.FluentdClients = append(pools.FluentdClients, &f)
	}
	log.Debugf("CreateFluentdConnPools: %v", pools)
	return &pools
}
