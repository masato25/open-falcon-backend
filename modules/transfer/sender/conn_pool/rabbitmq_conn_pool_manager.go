package conn_pool

import (
	"fmt"
	"net"
	"net/rpc/jsonrpc"
	"time"

	cmodel "github.com/open-falcon/common/model"
)

func newRabbitMqConnPool(address string, maxConns int, maxIdle int, connTimeout int) *ConnPool {
	connectionTimeout := time.Duration(connTimeout) * time.Millisecond
	p := NewConnPool("rabbitmq", address, maxConns, maxIdle)

	p.New = func(connName string) (NConn, error) {
		_, err := net.ResolveTCPAddr("tcp", p.Address)
		if err != nil {
			return nil, err
		}

		conn, err := net.DialTimeout("tcp", p.Address, connectionTimeout)
		if err != nil {
			return nil, err
		}

		return RpcClient{cli: jsonrpc.NewClient(conn), name: connName}, nil
	}

	return p
}

type RabbitMqConnPoolHelper struct {
	p           *ConnPool
	maxConns    int
	maxIdle     int
	connTimeout int
	callTimeout int
	address     string
}

func NewRabbitMqConnPoolHelper(address string, maxConns, maxIdle, connTimeout, callTimeout int) *RabbitMqConnPoolHelper {
	return &RabbitMqConnPoolHelper{
		p:           newRabbitMqConnPool(address, maxConns, maxIdle, connTimeout),
		maxConns:    maxConns,
		maxIdle:     maxIdle,
		connTimeout: connTimeout,
		callTimeout: callTimeout,
		address:     address,
	}
}

// A synchronous call; return if completed or time-out
func (this RabbitMqConnPoolHelper) Call(items []*cmodel.MetricValue) (err error) {

	for _, item := range items {
		fmt.Printf("got item: %v", item)
	}

	// Write the batch
	return
}

func (this *RabbitMqConnPoolHelper) Destroy() {
	if this.p != nil {
		this.p.Destroy()
	}
}
