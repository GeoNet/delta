package main

import (
	"fmt"
	"strings"

	"github.com/GeoNet/delta/meta"
)

func (t *Tilde) Geomag(set *meta.Set, geomag string, extra ...string) error {

	var stns []Station
	for _, s := range set.Stations() {

		// needs to match the station code
		if s.Network != geomag {
			continue
		}

		// look for all sensors at a station
		var sens []Sensor
		for _, x := range set.Sites() {
			if x.Station != s.Code {
				continue
			}

			var installs []meta.InstalledSensor
			for _, c := range set.Collections(x) {
				installs = append(installs, c.InstalledSensor)
			}

			// skip if there are no installed sensors
			if len(installs) == 0 {
				continue
			}

			sens = append(sens, Sensor{
				Code:  x.Location,
				Start: toTimePtr(x.Start),
				End:   toTimePtr(x.End),

				Latitude:       toFloat(fmt.Sprintf("%.4f", x.Latitude)),
				Longitude:      toFloat(fmt.Sprintf("%.4f", x.Longitude)),
				Elevation:      toFloat(fmt.Sprintf("%.0f", x.Elevation)),
				RelativeHeight: toFloat(fmt.Sprintf("%.0f", -x.Depth)),
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

	for _, e := range extra {
		parts := strings.Split(e, "_")
		if len(parts) != 3 {
			return fmt.Errorf("invalid extra channel: %s", e)
		}
		net, sta, loc := parts[0], parts[1], parts[2]

		s, ok := set.Station(sta)
		if !ok || s.Network != net {
			continue
		}

		// look for all sensors at a station
		var sens []Sensor
		for _, x := range set.Sites() {
			if x.Station != s.Code || x.Location != loc {
				continue
			}

			var installs []meta.InstalledSensor
			for _, c := range set.Collections(x) {
				installs = append(installs, c.InstalledSensor)
			}

			// skip if there are no installed sensors
			if len(installs) == 0 {
				continue
			}

			sens = append(sens, Sensor{
				Code:  x.Location,
				Start: toTimePtr(s.Start),
				End:   toTimePtr(s.End),

				Latitude:       toFloat(fmt.Sprintf("%.4f", x.Latitude)),
				Longitude:      toFloat(fmt.Sprintf("%.4f", x.Longitude)),
				Elevation:      toFloat(fmt.Sprintf("%.0f", x.Elevation)),
				RelativeHeight: toFloat(fmt.Sprintf("%.0f", -x.Depth)),
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
		Name:        "geomag",
		Description: "Geomagnetic Sensors",
		Stations:    stns,
	})

	return nil
}
