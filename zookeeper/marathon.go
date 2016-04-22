package zookeeper

import (
	"fmt"
	"sync"
	"time"

	"github.com/apognu/marathon-consul/util"
	"github.com/samuel/go-zookeeper/zk"
)

type marathonDiscovery struct {
	zoo *zk.Conn

	mutex *sync.Mutex
	Nodes []string
}

func newMarathonDiscovery(zoo *zk.Conn) *marathonDiscovery {
	return &marathonDiscovery{
		zoo:   zoo,
		mutex: &sync.Mutex{},
	}
}

func (m *marathonDiscovery) Init() {
	m.discoverNodes()

	go func() {
		for {
			ch, err := m.marathonWatch(util.Config.MarathonPath)
			if err != nil {
				time.Sleep(2 * time.Second)
				continue
			}

			if (<-ch).Type == zk.EventNodeChildrenChanged {
				m.discoverNodes()
			}
		}
	}()
}

func (m *marathonDiscovery) marathonWatch(path string) (<-chan zk.Event, error) {
	_, _, ch, err := m.zoo.ChildrenW(util.Config.MarathonPath)
	return ch, err
}

func (m *marathonDiscovery) discoverNodes() {
	marathon, _, err := m.zoo.Children(util.Config.MarathonPath)
	if err != nil {
		return
	}

	m.mutex.Lock()

	m.Nodes = []string{}
	for _, leaf := range marathon {
		node, _, err := m.zoo.Get(fmt.Sprintf("%s/%s", util.Config.MarathonPath, leaf))
		if err != nil {
			continue
		}

		m.Nodes = append(m.Nodes, string(node))
	}
	m.mutex.Unlock()
}
