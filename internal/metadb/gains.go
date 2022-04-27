package metadb

import (
	"path/filepath"
	"sync"
	"time"

	"github.com/GeoNet/delta/meta"
)

type gains struct {
	list   meta.GainList
	lookup map[string][]meta.Gain
	once   sync.Once
}

func (g *gains) loadGains(base string) error {
	var err error

	g.once.Do(func() {
		if err = meta.LoadList(filepath.Join(base, "install", "gains.csv"), &g.list); err == nil {
			lookup := make(map[string][]meta.Gain)
			for _, v := range g.list {
				lookup[v.Station] = append(lookup[v.Station], v.Gains()...)
			}
			g.lookup = lookup
		}
	})

	return err
}

func (m *MetaDB) Gain(sta, loc, cha string, at time.Time) (*meta.Gain, error) {

	if err := m.loadGains(m.base); err != nil {
		return nil, err
	}

	s, ok := m.gains.lookup[sta]
	if !ok {
		return nil, nil
	}
	for _, g := range s {
		if g.Location != loc {
			continue
		}
		if g.Channel != cha {
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
