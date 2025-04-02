package main

import (
	"fmt"
	"slices"

	"github.com/GeoNet/delta/meta"
)

func (t *Tilde) ScanDOAS(set *meta.Set, network string) error {

	var stns []Station
	for _, m := range set.Mounts() {

		if m.Network != network {
			continue
		}

		// look for all sensors at a station
		var sens []Sensor
		for _, f := range set.Features() {
			if f.Station != m.Code {
				continue
			}

			sensorFoundAlready := slices.ContainsFunc(sens, func(sensor Sensor) bool {
				return sensor.Code == f.Location
			})
			if sensorFoundAlready {
				continue
			}

			sens = append(sens, Sensor{
				Code:  f.Location,
				Start: toTimePtr(f.Start),
				End:   toTimePtr(f.End),

				Latitude:       toFloat(fmt.Sprintf("%.4f", m.Latitude)),
				Longitude:      toFloat(fmt.Sprintf("%.4f", m.Longitude)),
				Elevation:      toFloat(fmt.Sprintf("%.0f", m.Elevation)),
				RelativeHeight: toFloat(fmt.Sprintf("%.0f", -m.Depth)),
			})
		}

		// add the station to the list
		stns = append(stns, Station{
			Code:        m.Code,
			Description: m.Name,
			Start:       toTimePtr(m.Start),
			End:         toTimePtr(m.End),

			Latitude:  toFloat(fmt.Sprintf("%.4f", m.Latitude)),
			Longitude: toFloat(fmt.Sprintf("%.4f", m.Longitude)),
			Elevation: toFloat(fmt.Sprintf("%.0f", m.Elevation)),

			Sensors: sens,
		})
	}

	// update domains
	t.Domains = append(t.Domains, Domain{
		Name:        "scandoas",
		Description: "Continuous sulphur dioxide gas emission rates",
		Stations:    stns,
	})

	return nil
}
