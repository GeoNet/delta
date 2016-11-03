package metadb

import (
	"github.com/GeoNet/delta/meta"
)

func (db *MetaDB) DeployedDataloggers(place, role string) ([]meta.DeployedDatalogger, error) {
	txn := db.Txn(false)
	defer txn.Abort()

	lookup, err := txn.Get("datalogger", "role", place, role)
	if err != nil {
		return nil, err
	}

	var dataloggers []meta.DeployedDatalogger
	for d := lookup.Next(); d != nil; d = lookup.Next() {
		dataloggers = append(dataloggers, d.(meta.DeployedDatalogger))
	}

	return dataloggers, nil
}
