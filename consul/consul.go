package consul

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/apognu/marathon-consul/marathon"
	"github.com/apognu/marathon-consul/util"
	"github.com/hashicorp/consul/api"
)

type ApiServices map[string]interface{}

var (
	agent   *api.Agent
	catalog *api.Catalog
)

func Register(tasks []marathon.Task) {
	if agent == nil {
		client, err := api.NewClient(util.Config.ConsulConfig)
		if err != nil {
			logrus.Error("could not connect to consul")
			return
		}

		agent = client.Agent()
		catalog = client.Catalog()
	}

	for _, task := range tasks {
		NewService(task).Register()
	}

	deregisterServices(tasks)
}

func deregisterServices(tasks []marathon.Task) {
	services, _, _ := catalog.Services(nil)

	for svc, _ := range services {
		if svc == "consul" {
			continue
		}

		catalogTasks, _, _ := catalog.Service(svc, "", nil)

	DeregisterLoop:
		for _, ct := range catalogTasks {
			for _, t := range tasks {
				if len(t.Ports) == 1 {
					if normalizeServiceID(t.ID) == ct.ServiceID {
						fmt.Println(t.ID, "is still here")
						continue DeregisterLoop
					}
				}

				for p, _ := range t.Ports {
					if normalizeServiceIDWithPort(t.ID, p) == ct.ServiceID {
						fmt.Println(t.ID, "is still here")
						continue DeregisterLoop
					}
				}
			}

			DeregisterService(ct.Node, ct.Address, ct.ServiceID)
		}
	}
}
