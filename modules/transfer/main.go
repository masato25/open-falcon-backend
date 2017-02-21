package main

import (
	"fmt"
	"os"

	"github.com/Cepave/open-falcon-backend/common/logruslog"
	"github.com/Cepave/open-falcon-backend/common/vipercfg"
	"github.com/Cepave/open-falcon-backend/modules/transfer/g"
	"github.com/Cepave/open-falcon-backend/modules/transfer/http"
	"github.com/Cepave/open-falcon-backend/modules/transfer/proc"
	"github.com/Cepave/open-falcon-backend/modules/transfer/queue"
	"github.com/Cepave/open-falcon-backend/modules/transfer/receiver"
	"github.com/Cepave/open-falcon-backend/modules/transfer/rings"
	"github.com/Cepave/open-falcon-backend/modules/transfer/sender"
)

func main() {
	vipercfg.Parse()
	vipercfg.Bind()

	if vipercfg.Config().GetBool("version") {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}
	if vipercfg.Config().GetBool("vg") {
		fmt.Println(g.VERSION, g.COMMIT)
		os.Exit(0)
	}

	// global config
	vipercfg.Load()
	g.ParseConfig(vipercfg.Config().GetString("config"))
	logruslog.Init()

	// proc
	proc.Start()

	// rings (hash gen)
	rings.Start()

	// queue
	queue.Start()

	// receive
	receiver.Start()

	// sender
	sender.Start()

	// http
	http.Start()

	select {}
}
