package consul

import (
	"github.com/Sirupsen/logrus"
	"github.com/apognu/marathon-consul/marathon"
	"github.com/apognu/marathon-consul/util"
	"github.com/hashicorp/consul/api"
)

type ApiServices map[string]interface{}

var (
	agent *api.Agent
)

func Register(tasks []marathon.Task) {
	if agent == nil {
		client, err := api.NewClient(util.Config.ConsulConfig)
		if err != nil {
			logrus.Error("could not connect to consul")
			return
		}

		agent = client.Agent()
	}

	for _, task := range tasks {
		NewService(task).Register()
	}

	deregisterServices(tasks)
}

func deregisterServices(tasks []marathon.Task) {
	svc, err := agent.Services()
	if err != nil {
		return
	}

ServiceLoop:
	for id, _ := range svc {
		if id == "consul" {
			continue
		}

		for _, t := range tasks {
			for p, _ := range t.Ports {
				if len(t.Ports) == 1 && normalizeServiceID(t.ID) == id {
					continue ServiceLoop
				}

				if normalizeServiceIDWithPort(t.ID, p) == id {
					continue ServiceLoop
				}
			}
		}

		DeregisterService(id)
	}
}
