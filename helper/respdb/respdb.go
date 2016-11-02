package respdb

import (
	"github.com/GeoNet/delta/resp"
	"github.com/hashicorp/go-memdb"
)

type RespDB struct {
	*memdb.MemDB
}

func NewSchema() *memdb.DBSchema {
	return &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"datalogger": &memdb.TableSchema{
				Name: "datalogger",
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Name"},
					},
					"type": &memdb.IndexSchema{
						Name:         "type",
						AllowMissing: true,
						Indexer:      &memdb.StringFieldIndex{Field: "Type"},
					},
				},
			},
			"sensor": &memdb.TableSchema{
				Name: "sensor",
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Name"},
					},
					"type": &memdb.IndexSchema{
						Name:         "type",
						AllowMissing: true,
						Indexer:      &memdb.StringFieldIndex{Field: "Type"},
					},
				},
			},
			"response": &memdb.TableSchema{
				Name: "response",
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Name"},
					},
				},
			},
		},
	}
}

func NewRespDB(base string) (*RespDB, error) {

	db, err := memdb.NewMemDB(NewSchema())
	if err != nil {
		return nil, err
	}

	txn := db.Txn(true)
	for _, v := range resp.DataloggerModels {
		if err := txn.Insert("datalogger", v); err != nil {
			return nil, err
		}
	}
	for _, v := range resp.SensorModels {
		if err := txn.Insert("sensor", v); err != nil {
			return nil, err
		}
	}
	for _, v := range resp.Responses {
		if err := txn.Insert("response", v); err != nil {
			return nil, err
		}
	}
	txn.Commit()

	return &RespDB{MemDB: db}, nil
}
