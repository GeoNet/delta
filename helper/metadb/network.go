package metadb

import (
	"github.com/GeoNet/delta/meta"
)

func (db *MetaDB) Network(code string) (*meta.Network, error) {
	txn := db.Txn(false)
	defer txn.Abort()

	check, err := txn.First("network", "id", code)
	if err != nil {
		return nil, err
	}

	network := check.(meta.Network)

	return &network, nil
}
