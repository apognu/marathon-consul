package consul

import (
	"fmt"
	"strings"

	"github.com/apognu/marathon-consul/marathon"
	"github.com/apognu/marathon-consul/util"
	"github.com/hashicorp/consul/api"
)

type Service struct {
	ID      string
	Name    string
	Address string
	Port    []int
	Tags    []string
}

func NewService(task marathon.Task) *Service {
	tags := make([]string, 0)
	for k, v := range task.Labels {
		tags = append(tags, fmt.Sprintf("%s=%s", k, v))
	}

	return &Service{
		ID:      task.ID,
		Name:    normalizeServiceName(task.AppID),
		Address: task.Host,
		Port:    task.Ports,
		Tags:    tags,
	}
}

func (svc *Service) Register() {
	if len(svc.Port) == 1 {
		catalog.Register(&api.CatalogRegistration{
			Node:    svc.Address,
			Address: svc.Address,
			Service: &api.AgentService{
				ID:      normalizeServiceID(svc.ID),
				Service: svc.Name,
				Address: svc.Address,
				Port:    svc.Port[0],
				Tags:    svc.Tags,
			},
		}, nil)
	}

	for p, port := range svc.Port {
		catalog.Register(&api.CatalogRegistration{
			Node:    svc.Address,
			Address: svc.Address,
			Service: &api.AgentService{
				ID:      normalizeServiceIDWithPort(svc.ID, p),
				Service: normalizeServiceNameWithPort(svc.Name, p),
				Address: svc.Address,
				Port:    port,
				Tags:    svc.Tags,
			},
		}, nil)
	}
}

func DeregisterService(node, addr, id string) {
	// agent.ServiceDeregister(id)
	fmt.Printf("deregistering %s\n", id)
	catalog.Deregister(&api.CatalogDeregistration{
		Node:      node,
		Address:   addr,
		ServiceID: id,
	}, nil)
}

func normalizeServiceName(name string) string {
	name = strings.Join(util.ReverseStringArray(strings.Split(name[1:], "/")), "/")

	name = strings.Map(func(r rune) rune {
		if strings.IndexRune("abcdefghijklmnopqrstuvwxyz0123456789/", r) >= 0 {
			return r
		}
		return -1
	}, strings.ToLower(name))

	return strings.Replace(name, "/", "-", -1)
}

func normalizeServiceNameWithPort(name string, port int) string {
	return fmt.Sprintf("port%d-%s", port, name)
}

func normalizeServiceID(id string) string {
	return fmt.Sprintf("marathon-consul:%s", id)
}

func normalizeServiceIDWithPort(id string, port int) string {
	return fmt.Sprintf("marathon-consul:%s:port:%d", id, port)
}
