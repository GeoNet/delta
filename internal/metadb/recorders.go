package metadb

import (
	"path/filepath"
	"sync"

	"github.com/GeoNet/delta/meta"
)

type recorders struct {
	list     meta.InstalledRecorderList
	stations map[string][]meta.InstalledRecorder
	once     sync.Once
}

func (r *recorders) loadInstalledRecorders(base string) error {
	var err error

	r.once.Do(func() {
		if err = meta.LoadList(filepath.Join(base, "install", "recorders.csv"), &r.list); err == nil {
			stations := make(map[string][]meta.InstalledRecorder)
			for _, v := range r.list {
				if _, ok := stations[v.Station]; !ok {
					stations[v.Station] = []meta.InstalledRecorder{}
				}
				stations[v.Station] = append(stations[v.Station], v)
			}
			r.stations = stations
		}
	})

	return err
}

func (m *MetaDB) StationInstalledRecorders(sta string) ([]meta.InstalledRecorder, error) {
	if err := m.loadInstalledRecorders(m.base); err != nil {
		return nil, err
	}

	if r, ok := m.recorders.stations[sta]; ok {
		return r, nil
	}

	return nil, nil
}
