package http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Cepave/open-falcon-backend/modules/transfer/sender/pools"
)

func configDebugHttpRoutes() {
	// conn pools
	http.HandleFunc("/debug/connpool/", func(w http.ResponseWriter, r *http.Request) {
		urlParam := r.URL.Path[len("/debug/connpool/"):]
		args := strings.Split(urlParam, "/")
		mpool := pools.GetPools()
		argsLen := len(args)
		if argsLen < 1 {
			w.Write([]byte(fmt.Sprintf("bad args\n")))
			return
		}

		var result string
		receiver := args[0]
		switch receiver {
		case "judge":
			result = strings.Join(mpool.JudgeConnPools.Proc(), "\n")
		case "graph":
			result = strings.Join(mpool.GraphConnPools.Proc(), "\n")
		default:
			result = fmt.Sprintf("bad args, module not exist\n")
		}
		w.Write([]byte(result))
	})
}
