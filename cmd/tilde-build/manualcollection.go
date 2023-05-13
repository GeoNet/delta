package main

import (
	"fmt"

	"github.com/GeoNet/delta/meta"
)

func (t *Tilde) ManualCollection(set *meta.Set, networks ...string) error {

	var stns []Station
	for _, s := range set.Samples() {
		var sens []Sensor
		for _, c := range set.Sites() {
			if c.Station != s.Code {
				continue
			}

			// find list of installed sensors
			var installs []meta.InstalledSensor
			for _, v := range set.InstalledSensors() {
				if v.Station != c.Station {
					continue
				}
				if v.Location != c.Location {
					continue
				}
				installs = append(installs, v)
			}

			// skip if there are installed sensors
			if len(installs) > 0 {
				continue
			}

			sens = append(sens, Sensor{
				Code:  c.Location,
				Start: toTimePtr(c.Start),
				End:   toTimePtr(c.End),

				Latitude:       toFloat(fmt.Sprintf("%.4f", c.Latitude)),
				Longitude:      toFloat(fmt.Sprintf("%.4f", c.Longitude)),
				Elevation:      toFloat(fmt.Sprintf("%.0f", c.Elevation)),
				RelativeHeight: toFloat(fmt.Sprintf("%.0f", -c.Depth)),
			})

		}

		if len(sens) == 0 {
			continue
		}

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
		Name:        "manualcollection",
		Description: "Manually Collected Samples",
		Stations:    stns,
	})

	return nil
}
