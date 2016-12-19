package metadb

import (
	"github.com/GeoNet/delta/meta"
)

func (db *MetaDB) StationInstalledSensors(station string) ([]meta.InstalledSensor, error) {
	txn := db.Txn(false)
	defer txn.Abort()

	lookup, err := txn.Get("sensor", "station", station)
	if err != nil {
		return nil, err
	}

	var sensors []meta.InstalledSensor
	for s := lookup.Next(); s != nil; s = lookup.Next() {
		sensors = append(sensors, s.(meta.InstalledSensor))
	}

	return sensors, nil
}

func (db *MetaDB) StationLocationInstalledSensors(station, location string) ([]meta.InstalledSensor, error) {
	txn := db.Txn(false)
	defer txn.Abort()

	lookup, err := txn.Get("sensor", "site", station, location)
	if err != nil {
		return nil, err
	}

	var sensors []meta.InstalledSensor
	for s := lookup.Next(); s != nil; s = lookup.Next() {
		sensors = append(sensors, s.(meta.InstalledSensor))
	}

	return sensors, nil
}

func (db *MetaDB) ConnectionInstalledSensors(connection meta.Connection) ([]meta.InstalledSensor, error) {
	var sensors []meta.InstalledSensor

	streams, err := db.StationLocationStreams(connection.Station, connection.Location)
	if err != nil || streams == nil {
		return nil, err
	}

	installs, err := db.StationLocationInstalledSensors(connection.Station, connection.Location)
	if err != nil || installs == nil {
		return nil, err
	}

	for _, install := range installs {
		switch {
		case install.Start.After(connection.End):
			continue
		case install.End.Before(connection.Start):
			continue
		case install.Start == connection.End:
			continue
		}

		sensors = append(sensors, install)
	}

	return sensors, nil
}
