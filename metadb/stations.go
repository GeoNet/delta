package metadb

import (
	"path/filepath"
	"sync"

	"github.com/GeoNet/delta/meta"
)

type stations struct {
	list   meta.StationList
	lookup map[string]meta.Station
	once   sync.Once
}

func (s *stations) loadStations(base string) error {
	var err error

	s.once.Do(func() {
		if err = meta.LoadList(filepath.Join(base, "network", "stations.csv"), &s.list); err == nil {
			lookup := make(map[string]meta.Station)
			for _, v := range s.list {
				lookup[v.Code] = v
			}
			s.lookup = lookup
		}
	})

	return err
}

func (m *MetaDB) Stations() ([]meta.Station, error) {

	if err := m.loadStations(m.base); err != nil {
		return nil, err
	}

	return m.stations.list, nil
}

func (m *MetaDB) Station(code string) (*meta.Station, error) {

	if err := m.loadStations(m.base); err != nil {
		return nil, err
	}

	if s, ok := m.stations.lookup[code]; ok {
		return &s, nil
	}

	return nil, nil
}
