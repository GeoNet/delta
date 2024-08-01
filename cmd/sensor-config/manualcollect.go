package main

import (
	"sort"

	"github.com/GeoNet/delta/meta"
)

func (s Settings) ManualCollection(set *meta.Set, name, network string) (Group, bool) {

	var samples []Station
	for _, sample := range set.Samples() {

		net, ok := set.Network(sample.Network)
		if !ok {
			continue
		}
		if net.Code != network {
			continue
		}

		var sites []Site
		for _, point := range set.Points() {
			if point.Sample != sample.Code {
				continue
			}

			var sensors []Sensor
			for _, feature := range set.Features() {
				if feature.Station != sample.Code {
					continue
				}
				if feature.Location != point.Location {
					continue
				}

				sensors = append(sensors, Sensor{
					Code:        feature.Sublocation,
					Property:    feature.Property,
					Aspect:      feature.Aspect,
					Description: feature.Description,

					StartDate: feature.Start,
					EndDate:   feature.End,
				})
			}

			if !(len(sensors) > 0) {
				continue
			}

			sort.Slice(sensors, func(i, j int) bool {
				return sensors[i].Less(sensors[j])
			})

			sites = append(sites, Site{
				Code: point.Location,

				Latitude:  point.Latitude,
				Longitude: point.Longitude,
				Elevation: point.Elevation,
				Depth:     point.Depth,
				Datum:     point.Datum,
				Survey:    point.Survey,

				StartDate: point.Start,
				EndDate:   point.End,

				Features: sensors,
			})
		}

		if !(len(sites) > 0) {
			continue
		}

		sort.Slice(sites, func(i, j int) bool {
			return sites[i].Less(sites[j])
		})

		samples = append(samples, Station{
			Code:        sample.Code,
			Name:        sample.Name,
			Network:     sample.Network,
			External:    net.External,
			Description: net.Description,

			Latitude:  sample.Latitude,
			Longitude: sample.Longitude,
			Elevation: sample.Elevation,
			Depth:     sample.Depth,
			Datum:     sample.Datum,

			StartDate: sample.Start,
			EndDate:   sample.End,

			Sites: sites,
		})
	}

	return Group{Name: name, Samples: samples}, true
}
