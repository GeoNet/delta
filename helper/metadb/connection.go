package metadb

import (
	"github.com/GeoNet/delta/meta"
)

func (db *MetaDB) StationConnections(station string) ([]meta.Connection, error) {
	txn := db.Txn(false)
	defer txn.Abort()

	lookup, err := txn.Get("connection", "station", station)
	if err != nil {
		return nil, err
	}
	var connections []meta.Connection
	for c := lookup.Next(); c != nil; c = lookup.Next() {
		connections = append(connections, c.(meta.Connection))
	}

	return connections, nil
}
func (db *MetaDB) StationLocationConnections(station, location string) ([]meta.Connection, error) {
	txn := db.Txn(false)
	defer txn.Abort()

	lookup, err := txn.Get("connection", "site", station, location)
	if err != nil {
		return nil, err
	}
	var connections []meta.Connection
	for c := lookup.Next(); c != nil; c = lookup.Next() {
		connections = append(connections, c.(meta.Connection))
	}

	return connections, nil
}

func (db *MetaDB) DeployedDataloggerConnections(sensor meta.InstalledSensor, station, location string) ([]meta.DeployedDatalogger, error) {
	var dataloggers []meta.DeployedDatalogger

	connections, err := db.StationLocationConnections(station, location)
	if err != nil {
		return nil, err
	}
	for _, connection := range connections {
		if connection.Start.After(sensor.End) {
			continue
		}
		if connection.End.Before(sensor.Start) {
			continue
		}

		dataloggers, err := db.PlaceRoleDeployedDataloggers(connection.Place, connection.Role)
		if err != nil {
			return nil, err
		}
		for _, datalogger := range dataloggers {
			if connection.Start.After(datalogger.End) {
				continue
			}
			if connection.End.Before(datalogger.Start) {
				continue
			}

			dataloggers = append(dataloggers, datalogger)
		}
	}

	return dataloggers, nil
}

/*
func (db *MetaDB) InstalledSensorConnections(station string) ([]meta.InstalledSensor, error) {
	var sensors []meta.InstalledSensor

	connections, err := db.StationConnections(station)
	if err != nil {
		return nil, err
	}
	for _, connection := range connections {
		streams, err := db.StationLocationStreams(station, connection.Location)
		if err != nil {
			return nil, err
		}
		if streams == nil {
			continue
		}
		deploys, err := db.PlaceDeployedDataloggers(connection.Place)
		if err != nil {
			return nil, err
		}
		if deploys == nil {
			continue
		}
		site, err := db.Site(station, connection.Location)
		if err != nil {
			return nil, err
		}
		if site == nil {
			continue
		}
		installed, err := db.StationInstalledSensors(station)
		if err != nil {
			return nil, err
		}
		for _, sensor := range installed {
			switch {
			case sensor.Location != connection.Location:
				continue
			case sensor.Start.After(connection.End):
				continue
			case sensor.End.Before(connection.Start):
				continue
			case sensor.Start == connection.End:
				continue
			}

			sensors = append(sensors, sensor)
		}
	}

	return sensors, nil
}
*/
