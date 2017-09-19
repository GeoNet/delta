package main

import (
	"sort"
	"time"

	"github.com/GeoNet/delta/internal/metadb"
	"github.com/GeoNet/delta/meta"
)

func (cp ConfigPage) Gauge(base string) ([]Page, error) {

	db := metadb.NewMetaDB(base)

	var pages []Page

	chans := make(map[string]interface{})
	for _, c := range cp.Options.Channels() {
		chans[c] = true
	}

	locs := make(map[string]interface{})
	for _, l := range cp.Options.Locations {
		locs[l] = true
	}

	for _, net := range cp.Options.Networks {
		var stns []meta.Station
		list, err := db.NetworkStation(net)
		if err != nil {
			return nil, err
		}
		for _, s := range list {
			stns = append(stns, s)
		}

		sort.Sort(Stations(stns))

		var streams []Stream

		for _, s := range stns {
			chas, err := db.Channels(s.Code)
			if err != nil {
				return nil, err
			}

			for _, c := range chas {
				if time.Now().After(c.End) {
					continue
				}

				if _, ok := chans[c.Code]; !ok {
					continue
				}
				if _, ok := locs[c.Location]; !ok {
					continue
				}

				streams = append(streams, Stream{
					Auto:   "no",
					Colour: "#000000a0",
					Format: "amplitude",
					Rrd:    cp.Rrd(StationChannel(s, c), "amplitude/pressure/%s.%n/%s.%l-%c.%n.rrd"),
					Style:  "gauge",
					Tag: &Tag{
						Aligned: "first",
						Colour:  "#006400",
						Font:    "LiberationSans Narrow Bold 12",
						Rotated: "yes",
						String:  s.Code,
						XOffset: "5",
					},
					Name: &Name{
						Box:     "#ffffffd0",
						Colour:  "#006400",
						Font:    "LiberationSans Bold Italic 14",
						Pad:     "3",
						XOffset: "10",
						String:  s.Name,
					},
				})
			}
		}

		pages = append(pages, Page{
			Id:  cp.Id(OptionStream{Network: net}, "pressure-gauge-gauge-%n"),
			Png: cp.Png(OptionStream{Network: net}, "/volcano/%n/gauge.png"),
			Plots: []Plot{
				Plot{
					Id:      "gauge",
					Clip:    "100",
					Length:  "48 hours",
					Overlap: "50",
					Streams: streams,
					Max:     "860",
					Min:     "840",
					Width:   "795",
					Title: &Title{
						Colour: "#006400",
						Font:   "LiberationSans Bold Italic 14",
					},
					XGrid: &XGrid{
						Colour: "#ffffff70",
						Step:   "21600",
					},
					XLabel: &XLabel{
						Colour: "#006400",
						Font:   "LiberationSans Italic Bold 12",
						String: "Hours before current timestamp",
					},
					YLabel: &YLabel{
						Colour: "#006400",
						Font:   "LiberationSans Italic Bold 12",
						String: "Pressure level",
					},
					XTick: &XTick{
						Colour:  "#006400",
						Factor:  "3600",
						Font:    "LiberationSans Bold 12",
						Step:    "21600",
						YOffset: "5",
					},
					Border: &Border{
						Gradient: "#ffffffa0",
						Bg:       "#90EE90",
						Colour:   "#ffffff",
					},
					Copyright: &Copyright{
						Font: "LiberationSans Narrow 9",
					},
					Date: &Date{
						Font: "LiberationSans Narrow Bold 14",
					},
					Label: &Label{
						Colour: "#006400",
						Font:   "LiberationSans Bold Italic 14",
						String: "New Zealand Pressure Gauge Network",
					},
				},
			},
		})
	}

	return pages, nil
}
