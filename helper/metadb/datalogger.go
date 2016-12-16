package metadb

import (
	"github.com/GeoNet/delta/meta"
)

func (db *MetaDB) PlaceDeployedDataloggers(place string) ([]meta.DeployedDatalogger, error) {
	txn := db.Txn(false)
	defer txn.Abort()

	lookup, err := txn.Get("datalogger", "place", place)
	if err != nil {
		return nil, err
	}

	var dataloggers []meta.DeployedDatalogger
	for d := lookup.Next(); d != nil; d = lookup.Next() {
		dataloggers = append(dataloggers, d.(meta.DeployedDatalogger))
	}

	return dataloggers, nil
}

func (db *MetaDB) PlaceRoleDeployedDataloggers(place, role string) ([]meta.DeployedDatalogger, error) {
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

func (db *MetaDB) ConnectionInstalledSensorDeployedDataloggers(connection meta.Connection, install meta.InstalledSensor) ([]meta.DeployedDatalogger, error) {
	var dataloggers []meta.DeployedDatalogger

	deploys, err := db.PlaceDeployedDataloggers(connection.Place)
	if err != nil || deploys == nil {
		return nil, err
	}

	for _, deploy := range deploys {
		switch {
		case deploy.Role != connection.Role:
			continue
		case deploy.Start.After(connection.End):
			continue
		case deploy.End.Before(connection.Start):
			continue
		case deploy.Start == connection.End:
			continue
		case deploy.Start.After(install.End):
			continue
		case deploy.End.Before(install.Start):
			continue
		case deploy.Start == install.End:
			continue
		case install.End == deploy.Start:
			continue
		}

		dataloggers = append(dataloggers, deploy)
	}

	return dataloggers, nil
}
