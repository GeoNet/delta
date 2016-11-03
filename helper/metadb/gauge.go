package metadb

import (
	"github.com/GeoNet/delta/meta"
)

func (db *MetaDB) Gauges() ([]meta.Gauge, error) {
	var gauges []meta.Gauge

	txn := db.Txn(false)
	defer txn.Abort()

	lookup, err := txn.Get("gauge", "id")
	if err != nil {
		return nil, err
	}
	for g := lookup.Next(); g != nil; g = lookup.Next() {
		gauges = append(gauges, g.(meta.Gauge))
	}

	return gauges, nil
}
