package main

import (
	"regexp"
	"sort"
	"strings"

	"github.com/GeoNet/delta/meta"
)

func (n *Network) InstalledSensors(set *meta.Set, match *regexp.Regexp, network, prefix string) error {

	for _, stn := range set.Stations() {
		net, ok := set.Network(stn.Network)
		if !ok {
			continue
		}
		if net.Code != network {
			continue
		}

		var sites []Site
		for _, site := range set.Sites() {
			if site.Station != stn.Code {
				continue
			}

			if !match.MatchString(site.Location) {
				continue
			}

			sensors := make(map[Sensor][]string)

			for _, c := range set.Collections(site) {
				label := strings.TrimSpace(prefix + " " + c.Component.Type)

				sensor := Sensor{
					Type:  label,
					Make:  c.InstalledSensor.Make,
					Model: c.InstalledSensor.Model,

					Azimuth: c.InstalledSensor.Azimuth,
					Method:  c.InstalledSensor.Method,
					Dip:     c.InstalledSensor.Dip,

					Vertical: c.InstalledSensor.Vertical,
					North:    c.InstalledSensor.North,
					East:     c.InstalledSensor.East,

					StartDate: c.InstalledSensor.Start,
					EndDate:   c.InstalledSensor.End,
				}

				sensors[sensor] = append(sensors[sensor], c.Code())
			}

			var list []Sensor
			for sensor, chans := range sensors {

				dedupe := make(map[string]interface{})
				for _, c := range chans {
					dedupe[c] = true
				}
				var channels []string
				for k := range dedupe {
					channels = append(channels, k)
				}
				sort.Strings(channels)
				sensor.Channels = strings.Join(channels, ",")
				list = append(list, sensor)
			}

			sort.Slice(list, func(i, j int) bool {
				return list[i].Less(list[j])
			})

			sites = append(sites, Site{
				Code: site.Location,

				Latitude:  site.Latitude,
				Longitude: site.Longitude,
				Elevation: site.Elevation,
				Depth:     site.Depth,
				Datum:     site.Datum,
				Survey:    site.Survey,

				StartDate: site.Start,
				EndDate:   site.End,

				Sensors: list,
			})
		}

		sort.Slice(sites, func(i, j int) bool {
			return sites[i].Less(sites[j])
		})

		n.Stations = append(n.Stations, Station{
			Code:        stn.Code,
			Network:     net.External,
			Name:        stn.Name,
			Description: net.Description,

			Latitude:  stn.Latitude,
			Longitude: stn.Longitude,
			Elevation: stn.Elevation,
			Depth:     stn.Depth,
			Datum:     stn.Datum,

			StartDate: stn.Start,
			EndDate:   stn.End,

			Sites: sites,
		})
	}

	return nil
}

func (s Settings) InstalledSensors(set *meta.Set, name string, match *regexp.Regexp, networks string, list ...string) (Group, bool) {

	nets := make(map[string]interface{})
	for _, n := range strings.Split(networks, ",") {
		if n = strings.TrimSpace(n); n != "" {
			nets[n] = true
		}
	}

	types := make(map[string]interface{})
	for _, t := range list {
		types[strings.TrimSpace(t)] = true
	}

	var stations []Station
	for _, stn := range set.Stations() {
		net, ok := set.Network(stn.Network)
		if !ok {
			continue
		}
		if _, ok := nets[net.Code]; !ok {
			continue
		}

		var sites []Site
		for _, site := range set.Sites() {
			if site.Station != stn.Code {
				continue
			}

			if !match.MatchString(site.Location) {
				continue
			}

			sensors := make(map[Sensor][]string)

			for _, c := range set.Collections(site) {

				if _, ok := types[c.Component.Type]; !ok && len(types) > 0 {
					continue
				}

				sensor := Sensor{
					Make:  c.InstalledSensor.Make,
					Model: c.InstalledSensor.Model,
					Type:  c.Component.Type,

					Azimuth: c.InstalledSensor.Azimuth,
					Method:  c.InstalledSensor.Method,
					Dip:     c.InstalledSensor.Dip,

					Vertical: c.InstalledSensor.Vertical,
					North:    c.InstalledSensor.North,
					East:     c.InstalledSensor.East,

					StartDate: c.InstalledSensor.Start,
					EndDate:   c.InstalledSensor.End,
				}

				sensors[sensor] = append(sensors[sensor], c.Code())
			}

			var list []Sensor
			for sensor, chans := range sensors {
				dedupe := make(map[string]interface{})
				for _, c := range chans {
					dedupe[c] = true
				}
				var channels []string
				for k := range dedupe {
					channels = append(channels, k)
				}
				sort.Strings(channels)
				sensor.Channels = strings.Join(channels, ",")
				list = append(list, sensor)
			}

			// need at least one sensor
			if !(len(list) > 0) {
				continue
			}

			sort.Slice(list, func(i, j int) bool {
				return list[i].Less(list[j])
			})

			sites = append(sites, Site{
				Code: site.Location,

				Latitude:  site.Latitude,
				Longitude: site.Longitude,
				Elevation: site.Elevation,
				Depth:     site.Depth,
				Datum:     site.Datum,
				Survey:    site.Survey,

				StartDate: site.Start,
				EndDate:   site.End,

				Sensors: list,
			})
		}

		// need at least one site
		if !(len(sites) > 0) {
			continue
		}

		sort.Slice(sites, func(i, j int) bool {
			return sites[i].Less(sites[j])
		})

		stations = append(stations, Station{
			Code:        stn.Code,
			Name:        stn.Name,
			Network:     stn.Network,
			External:    net.External,
			Description: net.Description,

			Latitude:  stn.Latitude,
			Longitude: stn.Longitude,
			Elevation: stn.Elevation,
			Depth:     stn.Depth,
			Datum:     stn.Datum,

			StartDate: stn.Start,
			EndDate:   stn.End,

			Sites: sites,
		})
	}

	return Group{Name: name, Stations: stations}, true
}
