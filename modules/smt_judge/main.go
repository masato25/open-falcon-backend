package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Cepave/f2e-api_2/config"
	"github.com/Cepave/open-falcon-backend/modules/f2e-api/graph"
	"github.com/Cepave/open-falcon-backend/modules/smt_judge/read_graph"
	"github.com/spf13/viper"
)

var (
	Version   = "<UNDEFINED>"
	GitCommit = "<UNDEFINED>"
)

func initGraph() {
	graph.Start(viper.GetStringMapString("graphs.cluster"))
}

func main() {
	cfgTmp := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("v", false, "show version")
	help := flag.Bool("h", false, "help")
	flag.Parse()
	cfg := *cfgTmp
	if *version {
		fmt.Printf("version %s, build %s\n", Version, GitCommit)
		os.Exit(0)
	}

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	viper.AddConfigPath(".")
	viper.AddConfigPath("/")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("./api/config")
	cfg = strings.Replace(cfg, ".json", "", 1)
	viper.SetConfigName(cfg)

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	err = config.InitLog(viper.GetString("log_level"))
	if err != nil {
		log.Fatal(err)
	}
	initGraph()
	read_graph.Read()
	read_graph.ReadHttp()
}
