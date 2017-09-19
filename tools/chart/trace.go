package main

import (
	"sort"
	"time"

	"github.com/GeoNet/delta/internal/metadb"
	"github.com/GeoNet/delta/meta"
)

func (cp ConfigPage) Trace(base string) ([]Page, error) {

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

	var stns []meta.Station
	for _, net := range cp.Options.Networks {
		list, err := db.NetworkStation(net)
		if err != nil {
			return nil, err
		}
		for _, s := range list {
			stns = append(stns, s)
		}
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
				Auto:   "yes",
				Colour: "#000000a0",
				Format: "amplitude",
				Rrd:    cp.Rrd(StationChannel(s, c), "amplitude/drum/%s.%n/%s.%l-%c.%n.rrd"),
				Style:  "trace",
				Tags: []Tag{
					Tag{
						Aligned: "right",
						Box:     "#ffffffe0",
						Colour:  "#000000",
						Font:    "LiberationSans Bold Italic 7",
						String:  s.Name,
						XOffset: "-720",
					},
					Tag{
						Aligned: "left",
						Colour:  "#000000",
						Font:    "LiberationSans Bold 7",
						Rotated: "no",
						String:  s.Code,
					},
				},
			})
		}
	}

	pages = append(pages, Page{
		Height: "740",
		Width:  "800",
		Png:    cp.Png(OptionStream{}, "/earthquake/drums/latest/national-%a-trace.png"),
		Plots: []Plot{
			Plot{
				Id:      "drum",
				Clip:    "2",
				Height:  "670",
				Length:  "4 hours",
				Overlap: "0",
				Streams: streams,
				Title: &Title{
					Colour: "#000000",
					Font:   "LiberationSans Bold Italic 14",
				},
				XGrid: &XGrid{
					Colour: "#ffffff70",
					Step:   "1200",
				},
				XLabel: &XLabel{
					Colour: "#000000",
					Font:   "LiberationSans Italic Bold 12",
					String: "Minutes before current timestamp",
				},
				XTick: &XTick{
					Colour:  "#000000",
					Factor:  "60",
					Font:    "LiberationSans Bold 12",
					Step:    "1200",
					YOffset: "5",
				},
				Border: &Border{
					Bg:     "#ffffff",
					Colour: "#ffffff",
				},
				Copyright: &Copyright{
					Font: "LiberationSans Narrow 9",
				},
				Date: &Date{
					Font: "LiberationSans Narrow Bold 14",
				},
				Label: &Label{
					Colour: "#000000",
					Font:   "LiberationSans Bold Italic 14",
					String: "New Zealand National Seismograph Network",
				},
				Name: &Name{
					Colour: "#000000",
					Font:   "LiberationSans Bold Italic 14",
				},
			},
		},
	})

	return pages, nil
}
