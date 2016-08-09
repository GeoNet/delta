package main

import (
	"path/filepath"

	"github.com/GeoNet/delta/meta"
)

func ConnectionMap(install string) (map[string]meta.ConnectionList, error) {

	connectionMap := make(map[string]meta.ConnectionList)

	var cons meta.ConnectionList
	if err := meta.LoadList(filepath.Join(install, "connections.csv"), &cons); err != nil {
		return nil, err
	}

	for _, c := range cons {
		if _, ok := connectionMap[c.StationCode]; ok {
			connectionMap[c.StationCode] = append(connectionMap[c.StationCode], c)
		} else {
			connectionMap[c.StationCode] = meta.ConnectionList{c}
		}
	}

	var recs meta.InstalledRecorderList
	if err := meta.LoadList(filepath.Join(install, "recorders.csv"), &recs); err != nil {
		return nil, err
	}
	for _, r := range recs {
		c := meta.Connection{
			StationCode:  r.StationCode,
			LocationCode: r.LocationCode,
			Span: meta.Span{
				Start: r.Start,
				End:   r.End,
			},
			Place: r.StationCode,
			Role:  r.LocationCode,
		}
		if _, ok := connectionMap[c.StationCode]; ok {
			connectionMap[c.StationCode] = append(connectionMap[c.StationCode], c)
		} else {
			connectionMap[c.StationCode] = meta.ConnectionList{c}
		}
	}

	return connectionMap, nil
}
