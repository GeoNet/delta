package main

import (
	"sort"
	"strings"

	"github.com/GeoNet/delta/meta"
)

func (n *Network) Gnss(set *meta.Set, antenna, receiver string) error {

	for _, mark := range set.Marks() {
		net, ok := set.Network(mark.Network)
		if !ok {
			continue
		}

		monument, ok := set.Monument(mark.Code)
		if !ok {
			continue
		}

		var receivers []Sensor
		for _, r := range set.DeployedReceivers() {
			if r.Mark != mark.Code {
				continue
			}

			receivers = append(receivers, Sensor{
				Make:  r.Make,
				Model: r.Model,
				Type:  receiver,

				StartDate: r.Start,
				EndDate:   r.End,
			})
		}

		sort.Slice(receivers, func(i, j int) bool {
			return receivers[i].Less(receivers[j])
		})

		var antennas []Sensor
		for _, a := range set.InstalledAntennas() {
			if a.Mark != mark.Code {
				continue
			}

			antennas = append(antennas, Sensor{
				Make:  a.Make,
				Model: a.Model,
				Type:  antenna,

				Vertical: a.Vertical,
				North:    a.North,
				East:     a.East,
				Azimuth:  a.Azimuth,

				StartDate: a.Start,
				EndDate:   a.End,
			})
		}

		sort.Slice(antennas, func(i, j int) bool {
			return antennas[i].Less(antennas[j])
		})

		n.Marks = append(n.Marks, Mark{
			Code:        mark.Code,
			Network:     net.External,
			Name:        mark.Name,
			DomesNumber: monument.DomesNumber,
			Description: net.Description,

			Latitude:  mark.Latitude,
			Longitude: mark.Longitude,
			Elevation: mark.Elevation,
			Datum:     mark.Datum,

			GroundRelationship: monument.GroundRelationship,

			MarkType:        monument.MarkType,
			MonumentType:    monument.Type,
			FoundationType:  monument.FoundationType,
			FoundationDepth: monument.FoundationDepth,
			Bedrock:         monument.Bedrock,
			Geology:         monument.Geology,

			StartDate: mark.Start,
			EndDate:   mark.End,

			Antennas:  antennas,
			Receivers: receivers,
		})
	}

	return nil
}

func (s Settings) Gnss(set *meta.Set, name, networks string) (Group, bool) {

	nets := make(map[string]interface{})
	for _, n := range strings.Split(networks, ",") {
		if n = strings.TrimSpace(n); n != "" {
			nets[n] = true
		}
	}

	var marks []Mark
	for _, mark := range set.Marks() {
		net, ok := set.Network(mark.Network)
		if !ok {
			continue
		}
		if _, ok := nets[net.Code]; !ok {
			continue
		}

		monument, ok := set.Monument(mark.Code)
		if !ok {
			continue
		}

		var receivers []Sensor
		for _, r := range set.DeployedReceivers() {
			if r.Mark != mark.Code {
				continue
			}

			receivers = append(receivers, Sensor{
				Make:  r.Make,
				Model: r.Model,
				Type:  "GNSS Receiver",

				StartDate: r.Start,
				EndDate:   r.End,
			})
		}

		if !(len(receivers) > 0) {
			continue
		}

		sort.Slice(receivers, func(i, j int) bool {
			return receivers[i].Less(receivers[j])
		})

		var antennas []Sensor
		for _, a := range set.InstalledAntennas() {
			if a.Mark != mark.Code {
				continue
			}

			antennas = append(antennas, Sensor{
				Make:  a.Make,
				Model: a.Model,
				Type:  "GNSS Antenna",

				Vertical: a.Vertical,
				North:    a.North,
				East:     a.East,
				Azimuth:  a.Azimuth,

				StartDate: a.Start,
				EndDate:   a.End,
			})
		}

		if !(len(antennas) > 0) {
			continue
		}

		sort.Slice(antennas, func(i, j int) bool {
			return antennas[i].Less(antennas[j])
		})

		marks = append(marks, Mark{
			Code:        mark.Code,
			Name:        mark.Name,
			DomesNumber: monument.DomesNumber,
			Network:     net.Code,
			Description: net.Description,

			Latitude:  mark.Latitude,
			Longitude: mark.Longitude,
			Elevation: mark.Elevation,
			Datum:     mark.Datum,

			GroundRelationship: monument.GroundRelationship,

			MarkType:        monument.MarkType,
			MonumentType:    monument.Type,
			FoundationType:  monument.FoundationType,
			FoundationDepth: monument.FoundationDepth,
			Bedrock:         monument.Bedrock,
			Geology:         monument.Geology,

			StartDate: mark.Start,
			EndDate:   mark.End,

			Antennas:  antennas,
			Receivers: receivers,
		})
	}

	return Group{Name: name, Marks: marks}, true
}
