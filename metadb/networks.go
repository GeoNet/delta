package metadb

import (
	"path/filepath"
	"sync"

	"github.com/GeoNet/delta/meta"
)

type networks struct {
	list   meta.NetworkList
	lookup map[string]meta.Network
	once   sync.Once
}

func (n *networks) loadNetworks(base string) error {
	var err error

	n.once.Do(func() {
		if err = meta.LoadList(filepath.Join(base, "network", "networks.csv"), &n.list); err == nil {
			lookup := make(map[string]meta.Network)
			for _, v := range n.list {
				lookup[v.Code] = v
			}
			n.lookup = lookup
		}
	})

	return err
}

func (m *MetaDB) Network(code string) (*meta.Network, error) {

	if err := m.loadNetworks(m.base); err != nil {
		return nil, err
	}

	if n, ok := m.networks.lookup[code]; ok {
		return &n, nil
	}

	return nil, nil
}
