package metadb

import (
	"regexp"

	"github.com/GeoNet/delta/meta"
)

func (db *MetaDB) Stations() ([]meta.Station, error) {
	txn := db.Txn(false)
	defer txn.Abort()

	var stations []meta.Station

	lookup, err := txn.Get("station", "id")
	if err != nil {
		return nil, err
	}
	for s := lookup.Next(); s != nil; s = lookup.Next() {
		stations = append(stations, s.(meta.Station))
	}

	return stations, nil
}

func (db *MetaDB) Station(code string) (meta.Station, error) {
	txn := db.Txn(false)
	defer txn.Abort()

	station, err := txn.First("station", "id", code)
	if err != nil {
		return meta.Station{}, err
	}

	return station.(meta.Station), nil
}

func (db *MetaDB) MatchStations(net, sta *regexp.Regexp) ([]meta.Station, error) {
	txn := db.Txn(false)
	defer txn.Abort()

	var stations []meta.Station

	lookup, err := txn.Get("station", "id")
	if err != nil {
		return nil, err
	}
	for s := lookup.Next(); s != nil; s = lookup.Next() {
		if sta == nil || sta.MatchString(s.(meta.Station).Code) {
			n, err := db.Network(s.(meta.Station).Network)
			if err != nil {
				return nil, err
			}
			if net == nil || net.MatchString(n.External) {
				stations = append(stations, s.(meta.Station))
			}
		}
	}

	return stations, nil
}
