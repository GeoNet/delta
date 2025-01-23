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
		stns = append(stns, Station{
			Code:        m.Code,
			Description: m.Name,
			Start:       toTimePtr(m.Start),
			End:         toTimePtr(m.End),
			Latitude:    toFloat(fmt.Sprintf("%0.4f", m.Latitude)),
			Longitude:   toFloat(fmt.Sprintf("%0.4f", m.Longitude)),
			Elevation:   toFloat(fmt.Sprintf("%0.0f", math.Round(m.Elevation))),
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
