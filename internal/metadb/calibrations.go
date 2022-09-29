package metadb

import (
	"path/filepath"
	"sync"
	"time"

	"github.com/GeoNet/delta/meta"
)

type calibrations struct {
	list   meta.CalibrationList
	lookup map[string][]meta.Calibration
	once   sync.Once
}

func (c *calibrations) loadCalibrations(base string) error {
	var err error

	c.once.Do(func() {
		if err = meta.LoadList(filepath.Join(base, "install", "calibrations.csv"), &c.list); err == nil {
			lookup := make(map[string][]meta.Calibration)
			for _, v := range c.list {
				lookup[v.Model] = append(lookup[v.Model], v)
			}
			c.lookup = lookup
		}
	})

	return err
}

func (m *MetaDB) Calibration(model, serial string, comp int, at time.Time) (*meta.Calibration, error) {

	if err := m.loadCalibrations(m.base); err != nil {
		return nil, err
	}

	s, ok := m.calibrations.lookup[model]
	if !ok {
		return nil, nil
	}
	for _, g := range s {
		if g.Serial != serial {
			continue
		}
		if g.Number != comp {
			continue
		}
		if g.Start.After(at) {
			continue
		}
		if g.End.Before(at) {
			continue
		}
		return &g, nil
	}

	return nil, nil
}
