package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/apognu/marathon-consul/consul"
	"github.com/apognu/marathon-consul/marathon"
	"github.com/apognu/marathon-consul/util"
	"github.com/apognu/marathon-consul/zookeeper"
)

func main() {
	util.ParseFlags()

	zoo := zookeeper.NewDiscovery(util.Config.Master)
	zoo.Init()

	if util.Config.Healthcheck {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("OK"))
		})

		go http.ListenAndServe(fmt.Sprintf(":%d", util.Config.HealthcheckPort), nil)
	}

	for {
		if len(zoo.Marathon.Nodes) == 0 {
			logrus.Error("could not connect to zookeeper")
			time.Sleep(util.Config.Interval)
			continue
		}

		apps, err := marathon.FetchApps(zoo.Marathon.Nodes)
		if err == nil {
			consul.Register(apps.Tasks)
		}

		time.Sleep(util.Config.Interval)
	}
}
