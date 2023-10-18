package main

import (
	"sort"

	"github.com/GeoNet/delta/meta"
)

// Constituent values need by the de-tiding algorithm
type Constituent struct {
	Name      string  `json:"name"`
	Number    int     `json:"number"`
	Amplitude float64 `json:"amplitude"`
	Lag       float64 `json:"lag"`
}

// Detide parameters need by the de-tiding algorithm
type Detide struct {
	TimeZone     float64       `json:"timezone"`
	Latitude     float64       `json:"latitude"`
	Constituents []Constituent `json:"constituents,omitempty"`
}

// NewDetide builds the de-tiding parameters for a given site deployment time and constituent list.
func NewDetide(set *meta.Set, site meta.Site) *Detide {

	for _, g := range set.Gauges() {
		if g.Code != site.Station {
			continue
		}
		if g.Start.After(site.End) {
			continue
		}
		if g.End.Before(site.Start) {
			continue
		}

		var c []Constituent
		for _, v := range set.Constituents() {
			if v.Gauge != site.Station {
				continue
			}
			if v.Start.After(site.End) {
				continue
			}
			if v.End.Before(site.Start) {
				continue
			}

			c = append(c, Constituent{
				Name:      v.Name,
				Number:    v.Number,
				Amplitude: v.Amplitude,
				Lag:       v.Lag,
			})
		}

		sort.Slice(c, func(i, j int) bool {
			return c[i].Number < c[j].Number
		})

		return &Detide{
			TimeZone:     g.TimeZone,
			Latitude:     g.Latitude,
			Constituents: c,
		}
	}

	return nil
}
