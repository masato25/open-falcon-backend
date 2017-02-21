package conn_pool

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFluentdQueue(t *testing.T) {
	Convey("Create A connection", t, func() {
		pool := CreateFluentdConnPools([]string{"127.0.0.1:24224", "127.0.0.1:24224"}, 1000)
		ShouldNotEqual(pool, nil)
		pool.Get()
	})

}
