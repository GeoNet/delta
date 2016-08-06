package main

import (
	"path/filepath"

	"github.com/GeoNet/delta/meta"
)

func SiteMap(network string) (map[string]map[string]meta.Site, error) {

	siteMap := make(map[string]map[string]meta.Site)

	var locs meta.SiteList
	if err := meta.LoadList(filepath.Join(network, "sites.csv"), &locs); err != nil {
		return nil, err
	}

	for _, l := range locs {
		if _, ok := siteMap[l.StationCode]; !ok {
			siteMap[l.StationCode] = make(map[string]meta.Site)
		}
		siteMap[l.StationCode][l.LocationCode] = l
	}

	return siteMap, nil
}
