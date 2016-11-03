package metadb

import (
	"github.com/GeoNet/delta/meta"
)

func (db *MetaDB) Constituents(code string) ([]meta.Constituent, error) {
	txn := db.Txn(false)
	defer txn.Abort()

	var constituents []meta.Constituent
	lookup, err := txn.Get("constituent", "gauge", code)
	if err != nil {
		return nil, err
	}
	for c := lookup.Next(); c != nil; c = lookup.Next() {
		constituents = append(constituents, c.(meta.Constituent))
	}

	return constituents, nil
}
