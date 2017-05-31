package read_graph

import (
	"fmt"
	"log"

	cmodel "github.com/Cepave/open-falcon-backend/common/model"
	cutils "github.com/Cepave/open-falcon-backend/common/utils"
	"github.com/Cepave/open-falcon-backend/modules/f2e-api/graph"
	"github.com/masato25/resty"
	"github.com/montanaflynn/stats"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

var ct = cmodel.GraphLastParam{
	Endpoint: "enp01",
	Counter:  "counter01",
}

func Read() {
	r, e := graph.Last(ct)
	if e != nil {
		log.Printf("%v", e)
	} else {
		log.Printf("result: %v", r)
	}
}

func ReadHttp() {
	pkey := cutils.PK2(ct.Endpoint, ct.Counter)
	p, err := graph.GraphNodeRing.GetNode(pkey)
	log.Printf("key: %v , err: %v", p, err)
	hostAddr := viper.GetString(fmt.Sprintf("graphs.clusterHttp.%s", p))
	log.Printf("host addr: %s", hostAddr)
	rt := resty.New()
	rbody, e := rt.R().Get(fmt.Sprintf("http://%s/history/%s/%s", hostAddr, ct.Endpoint, ct.Counter))
	if e != nil {
		log.Printf("%v", e)
	}
	// fmt.Printf("%v", rbody.String())
	result := gjson.Get(rbody.String(), "data.#.value")
	var rts []float64
	for _, r := range result.Array() {
		fmt.Printf("value: %v", r.Float())
		rts = append(rts, r.Float())
	}
	o, _ := stats.Max(rts)
	o2, _ := stats.Mean(rts)
	fmt.Println("max value", o)
	fmt.Println("mean value", o2)

	// []*cmodel.GraphItem
}
