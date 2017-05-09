package main

import (
	"sort"
	"time"

	"github.com/GeoNet/delta/meta"
	"github.com/GeoNet/delta/metadb"
)

func (cp ConfigPage) Depth(base string) ([]Page, error) {

	if cp.Options.Thumb != "" {
		return nil, nil
	}

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

			pages = append(pages, Page{
				Id:     cp.Id(StationChannel(s, c), "drum-volcano-%n-%s-%l-drum"),
				Height: "500",
				Width:  "760",
				Png:    cp.Png(StationChannel(s, c), "/volcano/%n/%s/%l/drum.png"),
				Plots: []Plot{
					Plot{
						Id:      "rsam",
						Clip:    "1",
						Height:  "440",
						Length:  "35 days",
						Overlap: "0",
						Width:   "700",
						X:       "40",
						Streams: []Stream{
							Stream{
								Auto:   "no",
								Colour: "#000000a0",
								Format: "amplitude",
								Gain:   "5.0",
								Rrd:    cp.Rrd(StationChannel(s, c), "/amplitude/depth/%s.%n/%s.%l-%c.%n.rrd"),
								Style:  "rsam",
								Date: &Date{
									Colour: "25 25 112",
									Font:   "LiberationSans Narrow Bold 12",
								},
								XTick: &XTick{
									Colour:  "25 25 112",
									Factor:  "86400",
									Font:    "LiberationSans Bold 11",
									IsDate:  "yes",
									Step:    "172800",
									String:  "%d",
									YOffset: "5",
								},
								Label: &Label{
									Colour: "25 25 112",
									Font:   "LiberationSans Bold 12",
									String: s.Code + "/" + c.Location + "-" + c.Code + "/" + s.Network,
								},
								Copyright: &Copyright{
									Aligned: "top",
									Colour:  "25 25 112",
									Font:    "LiberationSans Bold 11",
								},
								Title: &Title{
									Colour: "25 25 112",
									Font:   "LiberationSans Bold Italic 14",
									String: s.Name,
								},
								XGrid: &XGrid{
									Colour: "143 188 143",
									IsDate: "yes",
									Pen:    "0.5",
									Step:   "172800",
								},
								YGrid: &YGrid{
									Colour: "143 188 143",
									Hint:   "10",
									Pen:    "0.5",
									Step:   "5",
								},
								XLabel: &XLabel{
									Colour: "25 25 112",
									Font:   "LiberationSans Bold Italic 11",
									String: "Day of month",
								},
								YLabel: &YLabel{
									Colour: "25 25 112",
									Font:   "LiberationSans Bold Italic 11",
									Hint:   "5",
									Step:   "10",
									String: "Depth (m)",
								},
								YTick: &YTick{
									Colour:  "25 25 112",
									Font:    "LiberationSans Bold 11",
									Hint:    "10",
									Rotated: "yes",
									String:  "%g ",
								},
							},
						},
						Border: &Border{
							Bg:     "#f5f5f5",
							Colour: "137 104 205",
							Pen:    "2",
						},
					},
				},
			})
		}
	}

	return pages, nil
}
