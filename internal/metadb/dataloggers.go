package metadb

import (
	"path/filepath"
	"sync"

	"github.com/GeoNet/delta/meta"
)

type dataloggers struct {
	list  meta.DeployedDataloggerList
	roles map[string]map[string][]meta.DeployedDatalogger
	once  sync.Once
}

func (d *dataloggers) loadDeployedDataloggers(base string) error {
	var err error

	d.once.Do(func() {
		if err = meta.LoadList(filepath.Join(base, "install", "dataloggers.csv"), &d.list); err == nil {
			roles := make(map[string]map[string][]meta.DeployedDatalogger)
			for _, v := range d.list {
				if _, ok := roles[v.Place]; !ok {
					roles[v.Place] = make(map[string][]meta.DeployedDatalogger)
				}
				if _, ok := roles[v.Place][v.Role]; !ok {
					roles[v.Place][v.Role] = []meta.DeployedDatalogger{}
				}
				roles[v.Place][v.Role] = append(roles[v.Place][v.Role], v)
			}
			d.roles = roles
		}
	})

	return err
}

func (m *MetaDB) PlaceRoleDeployedDataloggers(place, role string) ([]meta.DeployedDatalogger, error) {

	if err := m.loadDeployedDataloggers(m.base); err != nil {
		return nil, err
	}

	if p, ok := m.dataloggers.roles[place]; ok {
		if r, ok := p[role]; ok {
			return r, nil
		}
	}

	return nil, nil
}

func (m *MetaDB) ConnectionInstalledSensorDeployedDataloggers(con meta.Connection, sen meta.InstalledSensor) ([]meta.DeployedDatalogger, error) {

	if con.Station != sen.Station || con.Location != sen.Location {
		return nil, nil
	}

	if err := m.loadDeployedDataloggers(m.base); err != nil {
		return nil, err
	}

	var loggers []meta.DeployedDatalogger

	if s, ok := m.dataloggers.roles[con.Place]; ok {
		if r, ok := s[con.Role]; ok {
			for _, d := range r {
				switch {
				case d.Start.After(con.End):
				case d.End.Before(con.Start):
				case d.Start.After(sen.End):
				case d.End.Before(sen.Start):
				default:
					loggers = append(loggers, d)
				}
			}
		}
	}

	return loggers, nil
}

func (m *MetaDB) DeployedDataloggerConnections(sensor meta.InstalledSensor, station, location string) ([]meta.DeployedDatalogger, error) {

	var dataloggers []meta.DeployedDatalogger

	connections, err := m.StationLocationConnections(station, location)
	if err != nil {
		return nil, err
	}
	for _, connection := range connections {
		switch {
		case connection.Start.After(sensor.End):
		case connection.End.Before(sensor.Start):
		default:
			dataloggers, err := m.PlaceRoleDeployedDataloggers(connection.Place, connection.Role)
			if err != nil {
				return nil, err
			}
			for _, datalogger := range dataloggers {
				switch {
				case connection.Start.After(datalogger.End):
				case connection.End.Before(datalogger.Start):
				default:
					dataloggers = append(dataloggers, datalogger)
				}
			}
		}
	}

	return dataloggers, nil
}
