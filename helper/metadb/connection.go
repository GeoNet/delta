package metadb

import (
	"github.com/GeoNet/delta/meta"
)

func (db *MetaDB) Connections(station, location string) ([]meta.Connection, error) {
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
