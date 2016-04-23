package util

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/consul/api"
)

var Config = &configOptions{}

type configOptions struct {
	Master          []string
	MarathonPath    string
	ConsulAddress   string
	Interval        time.Duration
	Healthcheck     bool
	HealthcheckPort int

	ConsulConfig *api.Config
}

func ParseFlags() {
	master := flag.String("master", "127.0.0.1:2181", "comma-separated list of zookeeper servers")
	flag.StringVar(&Config.MarathonPath, "marathon", "/marathon", "zookeeper path to marathon")
	flag.DurationVar(&Config.Interval, "interval", 10*time.Second, "interval in seconds for marathon checks")
	flag.StringVar(&Config.ConsulAddress, "consul", "127.0.0.1:8500", "address to a consul agent")
	flag.BoolVar(&Config.Healthcheck, "healthcheck", false, "should we expose an healthcheck endpoint")
	flag.IntVar(&Config.HealthcheckPort, "healthcheck-port", 8080, "port on which to expose the healthcheck endpoint")

	flag.Parse()

	Config.Master = strings.Split(*master, ",")
	Config.MarathonPath = fmt.Sprintf("%s/leader", Config.MarathonPath)

	Config.ConsulConfig = api.DefaultConfig()
	Config.ConsulConfig.Address = Config.ConsulAddress
}
