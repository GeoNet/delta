package respdb

import (
	"github.com/GeoNet/delta/resp"
)

func (db *RespDB) Sensor(model string) (*resp.SensorModel, error) {
	txn := db.Txn(false)
	defer txn.Abort()

	entry, err := txn.First("sensor", "id", model)
	if err != nil || entry == nil {
		return nil, err
	}

	sensor := resp.SensorModel(entry.(resp.SensorModel))

	return &sensor, nil
}

func (db *RespDB) Datalogger(model string) (*resp.DataloggerModel, error) {
	txn := db.Txn(false)
	defer txn.Abort()

	entry, err := txn.First("datalogger", "id", model)
	if err != nil || entry == nil {
		return nil, err
	}

	datalogger := resp.DataloggerModel(entry.(resp.DataloggerModel))

	return &datalogger, nil
}

/*
func (db *RespDB) SensorModel(model string) (*resp.Sensor, error) {
	txn := db.Txn(false)
	defer txn.Abort()

	lookup, err := txn.Get("sensor", "role", place, role)
	if err != nil {
		return nil, err
	}

	var dataloggers []meta.DeployedDatalogger
	for d := lookup.Next(); d != nil; d = lookup.Next() {
		dataloggers = append(dataloggers, d.(meta.DeployedDatalogger))
	}

	return dataloggers, nil
}
*/

/*
func (db *RespDB) SensorModels(model string) ([]resp.Sensor, error) {
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
*/
