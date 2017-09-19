package metadb

import (
	"path/filepath"
	"sync"

	"github.com/GeoNet/delta/meta"
)

type gauges struct {
	list meta.GaugeList
	once sync.Once
}

func (g *gauges) loadGauges(base string) error {
	var err error

	g.once.Do(func() {
		err = meta.LoadList(filepath.Join(base, "network", "gauges.csv"), &g.list)
	})

	return err
}

func (m *MetaDB) Gauges() ([]meta.Gauge, error) {
	if err := m.loadGauges(m.base); err != nil {
		return nil, err
	}
	return m.gauges.list, nil
}
