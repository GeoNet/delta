package main

import (
	"fmt"
	"math"

	"github.com/GeoNet/delta/meta"
)

func (t *Tilde) Gnss(set *meta.Set) error {
	var stns []Station

	// all marks are GNSS stations
	for _, m := range set.Marks() {

		// GNSS data does not explicitly have sensor information, which is currently
		// problematic with Tilde's expectation of station + sensor,
		// see: https://github.com/GeoNet/tickets/issues/17327.
		// Simply create a nil sensor, which inherits station level information

		sens := []Sensor{
			Sensor{
				Code:      "nil",
				Start:     toTimePtr(m.Start),
				End:       toTimePtr(m.End),
				Latitude:  toFloat(fmt.Sprintf("%0.4f", m.Latitude)),
				Longitude: toFloat(fmt.Sprintf("%0.4f", m.Longitude)),
				Elevation: toFloat(fmt.Sprintf("%0.0f", math.Round(m.Elevation))),
				Datum:     m.Datum,
			},
		}

		stns = append(stns, Station{
			Code:        m.Code,
			Description: m.Name,
			Start:       sens[0].Start,
			End:         sens[0].End,
			Latitude:    sens[0].Latitude,
			Longitude:   sens[0].Longitude,
			Elevation:   sens[0].Elevation,
			Sensors:     sens,
		})
	}

	// update domains
	t.Domains = append(t.Domains, Domain{
		Name:        "gnss",
		Description: "Global Navigation Satellite System",
		Stations:    stns,
	})

	return nil
}
