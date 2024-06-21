package main

import (
	"fmt"
	"strings"

	"github.com/GeoNet/delta/meta"
)

func (t *Tilde) Coastal(set *meta.Set, network string) error {

	check := make(map[string]interface{})
	for _, n := range strings.Split(network, ",") {
		if s := strings.TrimSpace(n); s != "" {
			check[s] = true
		}
	}

	externals := make(map[string]string)
	for _, n := range set.Networks() {
		externals[n.Code] = n.External
	}

	codes := make(map[string]string)
	for _, s := range set.Stations() {
		codes[s.Code] = s.Network
	}

	var stns []Station

	for _, s := range set.Stations() {
		// must have a network code
		n, ok := codes[s.Code]
		if !ok {
			continue
		}
		if _, ok := check[n]; !ok {
			continue
		}

		// that code must have an external code
		if _, ok := externals[n]; !ok {
			continue
		}

		// look for all sensors at a station
		var sens []Sensor
		for _, c := range set.Sites() {
			if c.Station != s.Code {
				continue
			}

			//this construct avoids -0.0 in xml output
			var height float64
			if d := c.Depth; d != 0.0 {
				height = -d
			}

			sens = append(sens, Sensor{
				Code:  c.Location,
				Start: toTimePtr(c.Start),
				End:   toTimePtr(c.End),

				Latitude:       toFloat(fmt.Sprintf("%.4f", c.Latitude)),
				Longitude:      toFloat(fmt.Sprintf("%.4f", c.Longitude)),
				Elevation:      toFloat(fmt.Sprintf("%.0f", c.Elevation)),
				RelativeHeight: toFloat(fmt.Sprintf("%.0f", height)),
			})
		}

		// add the station to the list
		stns = append(stns, Station{
			Code:        s.Code,
			Description: s.Name,
			Start:       toTimePtr(s.Start),
			End:         toTimePtr(s.End),

			Latitude:  toFloat(fmt.Sprintf("%.4f", s.Latitude)),
			Longitude: toFloat(fmt.Sprintf("%.4f", s.Longitude)),
			Elevation: toFloat(fmt.Sprintf("%.0f", s.Elevation)),

			Sensors: sens,
		})
	}

	// update domains
	t.Domains = append(t.Domains, Domain{
		Name:        "coastal",
		Description: "Coastal Tsunami Gauge Network",
		Stations:    stns,
	})

	return nil
}
