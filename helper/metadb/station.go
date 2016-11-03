package metadb

import (
	"github.com/GeoNet/delta/meta"
)

func (db *MetaDB) Station(code string) (meta.Station, error) {
	txn := db.Txn(false)
	defer txn.Abort()

	station, err := txn.First("station", "id", code)
	if err != nil {
		return meta.Station{}, err
	}

	return station.(meta.Station), nil
}
