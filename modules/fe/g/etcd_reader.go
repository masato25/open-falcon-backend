package g

import (
	"time"

	"github.com/Cepave/open-falcon-backend/common/vipercfg"
	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"
)

func GetConf() (resp string, err error) {
	envKey := vipercfg.Config().GetString("ConfEnv")
	cfg := client.Config{
		Endpoints: []string{"http://127.0.0.1:2379"},
		Transport: client.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}
	c, err := client.New(cfg)
	kapi := client.NewKeysAPI(c)
	client_resp, err := kapi.Get(context.Background(), "/"+envKey, nil)
	resp = client_resp.Node.Value
	return
}
