package metadb

import (
	"path/filepath"
	"sync"

	"github.com/GeoNet/delta/meta"
)

type connections struct {
	list      meta.ConnectionList
	stations  map[string][]meta.Connection
	locations map[string]map[string][]meta.Connection
	once      sync.Once
}

func (c *connections) loadConnections(base string) error {
	var err error

	c.once.Do(func() {
		if err = meta.LoadList(filepath.Join(base, "install", "connections.csv"), &c.list); err == nil {
			stations := make(map[string][]meta.Connection)
			for _, v := range c.list {
				if _, ok := stations[v.Station]; !ok {
					stations[v.Station] = []meta.Connection{}
				}
				stations[v.Station] = append(stations[v.Station], v)
			}
			c.stations = stations

			locations := make(map[string]map[string][]meta.Connection)
			for _, v := range c.list {
				if _, ok := locations[v.Station]; !ok {
					locations[v.Station] = make(map[string][]meta.Connection)
				}
				if _, ok := locations[v.Station][v.Location]; !ok {
					locations[v.Station][v.Location] = []meta.Connection{}
				}
				locations[v.Station][v.Location] = append(locations[v.Station][v.Location], v)
			}
			c.locations = locations
		}
	})

	return err
}

func (m *MetaDB) StationConnections(sta string) ([]meta.Connection, error) {
	if err := m.loadConnections(m.base); err != nil {
		return nil, err
	}

	if c, ok := m.connections.stations[sta]; ok {
		return c, nil
	}

	return nil, nil
}

func (m *MetaDB) StationLocationConnections(sta, loc string) ([]meta.Connection, error) {
	if err := m.loadConnections(m.base); err != nil {
		return nil, err
	}

	if s, ok := m.connections.locations[sta]; ok {
		if l, ok := s[loc]; ok {
			return l, nil
		}
	}

	return nil, nil
}
