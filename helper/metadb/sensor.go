package metadb

import (
	"github.com/GeoNet/delta/meta"
)

func (db *MetaDB) InstalledSensors(station, location string) ([]meta.InstalledSensor, error) {
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
