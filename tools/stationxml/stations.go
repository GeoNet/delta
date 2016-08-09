package main

import (
	"path/filepath"

	"github.com/GeoNet/delta/meta"
)

func StationMap(network string) (map[string]meta.Station, error) {

	stationMap := make(map[string]meta.Station)

	// load station details
	var s meta.StationList
	if err := meta.LoadList(filepath.Join(network, "stations.csv"), &s); err != nil {
		return nil, err
	}

	for _, v := range s {
		stationMap[v.Code] = v
	}

	return stationMap, nil
}
