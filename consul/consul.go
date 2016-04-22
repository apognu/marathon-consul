package consul

import (
	"github.com/Sirupsen/logrus"
	"github.com/apognu/marathon-consul/marathon"
	"github.com/hashicorp/consul/api"
)

type ApiServices map[string]interface{}

var (
	client, _ = api.NewClient(api.DefaultConfig())
	agent     = client.Agent()
)

func Register(tasks []marathon.Task) {
	logrus.Info("registering services")

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
				if normalizeServiceID(t.ID, p) == id {
					continue ServiceLoop
				}
			}
		}

		DeregisterService(id)
	}
}
