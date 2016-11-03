package metadb

import (
	"github.com/GeoNet/delta/meta"
)

func (db *MetaDB) Sites(station string) ([]meta.Site, error) {
	txn := db.Txn(false)
	defer txn.Abort()

	lookup, err := txn.Get("site", "station", station)
	if err != nil {
		return nil, err
	}

	var sites []meta.Site
	for s := lookup.Next(); s != nil; s = lookup.Next() {
		sites = append(sites, s.(meta.Site))
	}

	return sites, nil
}
