package metadb

import (
	"path/filepath"
	"sync"

	"github.com/GeoNet/delta/meta"
)

type sensors struct {
	list      meta.InstalledSensorList
	stations  map[string][]meta.InstalledSensor
	locations map[string]map[string][]meta.InstalledSensor
	once      sync.Once
}

func (s *sensors) loadInstalledSensors(base string) error {
	var err error

	s.once.Do(func() {
		if err = meta.LoadList(filepath.Join(base, "install", "sensors.csv"), &s.list); err == nil {
			stations := make(map[string][]meta.InstalledSensor)
			for _, v := range s.list {
				if _, ok := stations[v.Station]; !ok {
					stations[v.Station] = []meta.InstalledSensor{}
				}
				stations[v.Station] = append(stations[v.Station], v)
			}
			s.stations = stations

			locations := make(map[string]map[string][]meta.InstalledSensor)
			for _, v := range s.list {
				if _, ok := locations[v.Station]; !ok {
					locations[v.Station] = make(map[string][]meta.InstalledSensor)
				}
				if _, ok := locations[v.Station][v.Location]; !ok {
					locations[v.Station][v.Location] = []meta.InstalledSensor{}
				}
				locations[v.Station][v.Location] = append(locations[v.Station][v.Location], v)
			}
			s.locations = locations
		}
	})

	return err
}

func (m *MetaDB) StationInstalledSensors(sta string) ([]meta.InstalledSensor, error) {
	if err := m.loadInstalledSensors(m.base); err != nil {
		return nil, err
	}

	if s, ok := m.sensors.stations[sta]; ok {
		return s, nil
	}

	return nil, nil
}

func (m *MetaDB) StationLocationInstalledSensors(sta, loc string) ([]meta.InstalledSensor, error) {
	if err := m.loadInstalledSensors(m.base); err != nil {
		return nil, err
	}

	if s, ok := m.sensors.locations[sta]; ok {
		if l, ok := s[loc]; ok {
			return l, nil
		}
	}

	return nil, nil
}
