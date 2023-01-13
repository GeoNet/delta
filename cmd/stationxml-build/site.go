package main

import (
	"github.com/GeoNet/delta/meta"
)

// toSiteName builds a station name from a meta Station.
func toSiteName(station meta.Station) string {
	if station.Name != "" {
		return station.Name
	}
	return station.Code
}
