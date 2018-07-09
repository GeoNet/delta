package metadb

import (
	"path/filepath"
	"sync"

	"github.com/GeoNet/delta/meta"
)

type constituents struct {
	list   meta.ConstituentList
	gauges map[string][]meta.Constituent
	once   sync.Once
}

func (c *constituents) loadConstituents(base string) error {
	var err error

	c.once.Do(func() {
		if err = meta.LoadList(filepath.Join(base, "environment", "constituents.csv"), &c.list); err == nil {
			gauges := make(map[string][]meta.Constituent)
			for _, v := range c.list {
				if _, ok := gauges[v.Gauge]; !ok {
					gauges[v.Gauge] = []meta.Constituent{}
				}
				gauges[v.Gauge] = append(gauges[v.Gauge], v)
			}
			c.gauges = gauges
		}
	})

	return err
}

func (m *MetaDB) GaugeConstituents(gauge string) ([]meta.Constituent, error) {
	if err := m.loadConstituents(m.base); err != nil {
		return nil, err
	}

	if c, ok := m.constituents.gauges[gauge]; ok {
		return c, nil
	}

	return nil, nil
}
