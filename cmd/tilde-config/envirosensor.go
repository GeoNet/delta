package main

import (
	"fmt"

	"github.com/GeoNet/delta/meta"
)

func (t *Tilde) EnviroSensor(set *meta.Set, enviro string) error {

	valid := make(map[string]interface{})
	for _, f := range set.Features() {
		valid[f.Station] = true
	}

	var stns []Station
	for _, s := range set.Stations() {
		if _, ok := valid[s.Code]; !ok {
			continue
		}

		if s.Network != enviro {
			continue
		}

		// look for all sensors at a station
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

			// skip if there are no installed sensors
			if len(installs) == 0 {
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
		Name:        "envirosensor",
		Description: "Environmental Sensors",
		Stations:    stns,
	})

	return nil
}
