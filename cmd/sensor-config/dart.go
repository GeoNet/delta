package main

import (
	"sort"

	"github.com/GeoNet/delta/meta"
)

func (n *Network) Dart(set *meta.Set, network, label string) error {

	net, ok := set.Network(network)
	if !ok {
		return nil
	}

	for _, stn := range set.Stations() {
		if stn.Network != net.Code {
			continue
		}

		var sites []Site
		for _, site := range set.Sites() {
			if site.Station != stn.Code {
				continue
			}

			var sensors []Sensor

			for _, v := range set.InstalledRecorders() {
				if v.Station != site.Station {
					continue
				}
				if v.Location != site.Location {
					continue
				}

				sensors = append(sensors, Sensor{
					Make:  v.InstalledSensor.Make,
					Model: v.InstalledSensor.Model,
					Type:  label,

					StartDate: v.Start,
					EndDate:   v.End,
				})
			}

			sort.Slice(sensors, func(i, j int) bool {
				return sensors[i].Less(sensors[j])
			})

			sites = append(sites, Site{
				Code: site.Location,

				Latitude:  site.Latitude,
				Longitude: site.Longitude,
				Depth:     site.Depth,
				Datum:     site.Datum,
				Survey:    site.Survey,

				StartDate: site.Start,
				EndDate:   site.End,

				Sensors: sensors,
			})
		}

		sort.Slice(sites, func(i, j int) bool {
			return sites[i].Less(sites[j])
		})

		n.Buoys = append(n.Buoys, Station{
			Code:        stn.Code,
			Network:     net.External,
			Name:        stn.Name,
			Description: net.Description,

			Latitude:  stn.Latitude,
			Longitude: stn.Longitude,
			Depth:     stn.Depth,
			Datum:     stn.Datum,

			StartDate: stn.Start,
			EndDate:   stn.End,

			Sites: sites,
		})
	}

	return nil
}
