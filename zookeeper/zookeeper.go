package zookeeper

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/samuel/go-zookeeper/zk"
)

type discovery struct {
	Servers []string
	Conn    *zk.Conn

	Marathon *marathonDiscovery
}

func NewDiscovery(servers []string) *discovery {
	return &discovery{
		Servers: servers,
	}
}

func (z *discovery) Init() {
	var err error
	z.Conn, _, err = zk.Connect(z.Servers, 10*time.Second)
	if err != nil {
		logrus.Fatal(err)
	}

	z.Marathon = newMarathonDiscovery(z.Conn)
	z.Marathon.Init()
}
