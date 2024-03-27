package main

import (
	"sort"

	"github.com/GeoNet/delta/meta"
)

func (n *Network) ManualCollection(set *meta.Set, network, label string) error {

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
					Type:        label,
					Description: feature.Description,

					StartDate: feature.Start,
					EndDate:   feature.End,
				})
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

				Sensors: sensors,
			})
		}

		sort.Slice(sites, func(i, j int) bool {
			return sites[i].Less(sites[j])
		})

		n.Samples = append(n.Samples, Station{
			Code:        sample.Code,
			Network:     net.External,
			Name:        sample.Name,
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

	return nil
}
