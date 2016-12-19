package metadb

import (
	"github.com/GeoNet/delta/meta"
)

func (db *MetaDB) StationInstalledRecorders(station string) ([]meta.InstalledRecorder, error) {
	txn := db.Txn(false)
	defer txn.Abort()

	lookup, err := txn.Get("recorder", "station", station)
	if err != nil {
		return nil, err
	}

	var recorders []meta.InstalledRecorder
	for s := lookup.Next(); s != nil; s = lookup.Next() {
		recorders = append(recorders, s.(meta.InstalledRecorder))
	}

	return recorders, nil
}

func (db *MetaDB) StationLocationInstalledRecorders(station, location string) ([]meta.InstalledRecorder, error) {
	txn := db.Txn(false)
	defer txn.Abort()

	lookup, err := txn.Get("recorder", "site", station, location)
	if err != nil {
		return nil, err
	}

	var recorders []meta.InstalledRecorder
	for s := lookup.Next(); s != nil; s = lookup.Next() {
		recorders = append(recorders, s.(meta.InstalledRecorder))
	}

	return recorders, nil
}
