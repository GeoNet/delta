package metadb

import (
	"path/filepath"
	"sync"

	"github.com/GeoNet/delta/meta"
)

type sites struct {
	list      meta.SiteList
	stations  map[string][]meta.Site
	locations map[string]map[string]meta.Site
	once      sync.Once
}

func (s *sites) loadSites(base string) error {
	var err error

	s.once.Do(func() {
		if err = meta.LoadList(filepath.Join(base, "network", "sites.csv"), &s.list); err == nil {
			stations := make(map[string][]meta.Site)
			for _, v := range s.list {
				if _, ok := stations[v.Station]; !ok {
					stations[v.Station] = []meta.Site{}
				}
				stations[v.Station] = append(stations[v.Station], v)
			}
			s.stations = stations

			locations := make(map[string]map[string]meta.Site)
			for _, v := range s.list {
				if _, ok := locations[v.Station]; !ok {
					locations[v.Station] = make(map[string]meta.Site)
				}
				locations[v.Station][v.Location] = v
			}
			s.locations = locations
		}
	})

	return err
}

func (m *MetaDB) Sites(sta string) ([]meta.Site, error) {
	if err := m.loadSites(m.base); err != nil {
		return nil, err
	}

	if s, ok := m.sites.stations[sta]; ok {
		return s, nil
	}

	return nil, nil
}

func (m *MetaDB) Site(sta, loc string) (*meta.Site, error) {
	if err := m.loadSites(m.base); err != nil {
		return nil, err
	}

	if s, ok := m.sites.locations[sta]; ok {
		if l, ok := s[loc]; ok {
			return &l, nil
		}
	}

	return nil, nil
}
