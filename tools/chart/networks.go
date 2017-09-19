package main

import (
	"strconv"
	"time"

	"github.com/GeoNet/delta/internal/metadb"
	"github.com/GeoNet/delta/meta"
)

func (cp ConfigPage) Networks(base string) ([]Page, error) {

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

	excludes := make(map[string]interface{})
	for _, l := range cp.Options.Excludes {
		excludes[l] = true
	}

	for _, net := range cp.Options.Networks {

		lookup, err := db.Network(net)
		if err != nil {
			return nil, err
		}

		var stns []meta.Station
		list, err := db.NetworkStation(net)
		if err != nil {
			return nil, err
		}
		for _, s := range list {
			if _, ok := excludes[s.Code]; ok {
				continue
			}
			stns = append(stns, s)
		}

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
					Rrd:    cp.Rrd(StationChannel(s, c), "amplitude/rsam/%s.%n/%s.%l-%c.%n.rrd"),
					Style:  "rsam",
					Tags: []Tag{
						Tag{
							Aligned: "right",
							Box:     "#ffffffd0",
							Colour:  "#006400",
							Font:    "LiberationSans Bold Italic 10",
							String:  s.Name,
							XOffset: "-720",
						},
						Tag{
							Aligned: "left",
							Colour:  "#006400",
							Font:    "LiberationSans Bold 10",
							Rotated: "no",
							String:  s.Code,
						},
					},
				})
			}
		}

		pages = append(pages, Page{
			Id:     cp.Id(OptionStream{Network: net}, "rsam-volcano-%n"),
			Height: strconv.Itoa(len(streams)*40 + 70),
			Png:    cp.Png(OptionStream{Network: net}, "/volcano/%n/rsam.png"),
			Plots: []Plot{
				Plot{
					Id:      "rsam",
					Clip:    "1",
					Height:  strconv.Itoa(len(streams) * 40),
					Length:  "25 days",
					Overlap: "0",
					Streams: streams,
					Title: &Title{
						Colour: "#006400",
						Font:   "LiberationSans Bold Italic 14",
					},
					XGrid: &XGrid{
						Colour: "#ffffff70",
						IsDate: "yes",
						Step:   "172800",
					},
					XLabel: &XLabel{
						Colour: "#006400",
						Font:   "LiberationSans Italic Bold 12",
						String: "Day of month",
					},
					XTick: &XTick{
						Colour:  "#006400",
						Factor:  "86400",
						Font:    "LiberationSans Bold 12",
						IsDate:  "yes",
						Step:    "172800",
						YOffset: "5",
						String:  "%d",
					},
					Border: &Border{
						Bg:       "#AFEEEE",
						Colour:   "#ffffff",
						Gradient: "#ffffffa0",
					},
					Copyright: &Copyright{
						Font: "LiberationSans Narrow 9",
					},
					Dates: []Date{
						Date{
							Font: "LiberationSans Narrow Bold 14",
						},
						Date{
							Aligned: "bottom",
							Font:    "LiberationSans Narrow Bold 14",
							String:  "%b %d, %Y",
							Time:    "start",
						},
						Date{
							Aligned: "bottom",
							Font:    "LiberationSans Narrow Bold 14",
							String:  "%b %d, %Y",
							Time:    "end",
						},
					},
					Label: &Label{
						Colour: "#006400",
						Font:   "LiberationSans Bold Italic 14",
						String: lookup.Description,
					},
					Name: &Name{
						Colour: "#006400",
						Font:   "LiberationSans Bold Italic 14",
					},
				},
			},
		})
	}

	for _, net := range cp.Options.Networks {

		lookup, err := db.Network(net)
		if err != nil {
			return nil, err
		}

		var stns []meta.Station
		list, err := db.NetworkStation(net)
		if err != nil {
			return nil, err
		}
		for _, s := range list {
			if _, ok := excludes[s.Code]; ok {
				continue
			}
			stns = append(stns, s)
		}

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
					Format: "spectra",
					High:   "-80",
					Low:    "-180",
					Map:    "polar",
					Rrd:    cp.Rrd(StationChannel(s, c), "spectra/ssam/%s.%n/%s.%l-%c.%n.rrd"),
					Style:  "ssam",
					Tags: []Tag{
						Tag{
							Aligned: "right",
							Box:     "#ffffffd0",
							Colour:  "#006400",
							Font:    "LiberationSans Bold Italic 10",
							String:  s.Name,
							XOffset: "-720",
						},
						Tag{
							Aligned: "left",
							Colour:  "#006400",
							Font:    "LiberationSans Bold 10",
							Rotated: "no",
							String:  s.Code,
						},
					},
				})
			}
		}

		pages = append(pages, Page{
			Id:     cp.Id(OptionStream{Network: net}, "ssam-volcano-%n"),
			Height: strconv.Itoa(len(streams)*40 + 70),
			Png:    cp.Png(OptionStream{Network: net}, "/volcano/%n/ssam.png"),
			Plots: []Plot{
				Plot{
					Id:      "ssam",
					Clip:    "1",
					Height:  strconv.Itoa(len(streams) * 40),
					Length:  "25 days",
					Overlap: "-5",
					Streams: streams,
					Title: &Title{
						Colour: "#006400",
						Font:   "LiberationSans Bold Italic 14",
					},
					XLabel: &XLabel{
						Colour: "#006400",
						Font:   "LiberationSans Italic Bold 12",
						String: "Day of month",
					},
					XTick: &XTick{
						Colour:  "#006400",
						Factor:  "86400",
						Font:    "LiberationSans Bold 12",
						IsDate:  "yes",
						Step:    "172800",
						YOffset: "5",
						String:  "%d",
					},
					Border: &Border{
						Bg:       "#AFEEEE",
						Colour:   "#ffffff",
						Gradient: "#ffffffa0",
					},
					Copyright: &Copyright{
						Font: "LiberationSans Narrow 9",
					},
					Dates: []Date{
						Date{
							Font: "LiberationSans Narrow Bold 14",
						},
						Date{
							Aligned: "bottom",
							Font:    "LiberationSans Narrow Bold 14",
							String:  "%b %d, %Y",
							Time:    "start",
						},
						Date{
							Aligned: "bottom",
							Font:    "LiberationSans Narrow Bold 14",
							String:  "%b %d, %Y",
							Time:    "end",
						},
					},
					Label: &Label{
						Colour: "#006400",
						Font:   "LiberationSans Bold Italic 14",
						String: lookup.Description,
					},
					Name: &Name{
						Colour: "#006400",
						Font:   "LiberationSans Bold Italic 14",
					},
				},
			},
		})
	}

	for _, net := range cp.Options.Networks {

		lookup, err := db.Network(net)
		if err != nil {
			return nil, err
		}

		var stns []meta.Station
		list, err := db.NetworkStation(net)
		if err != nil {
			return nil, err
		}
		for _, s := range list {
			if _, ok := excludes[s.Code]; ok {
				continue
			}
			stns = append(stns, s)
		}

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
					Format: "spectra",
					Rrd:    cp.Rrd(StationChannel(s, c), "spectra/ssam/%s.%n/%s.%l-%c.%n.rrd"),
					Style:  "tremor",
					Tags: []Tag{
						Tag{
							Aligned: "right",
							Box:     "#ffffffd0",
							Colour:  "#006400",
							Font:    "LiberationSans Bold Italic 10",
							String:  s.Name,
							XOffset: "-720",
						},
						Tag{
							Aligned: "left",
							Colour:  "#006400",
							Font:    "LiberationSans Bold 10",
							Rotated: "no",
							String:  s.Code,
						},
					},
				})
			}
		}

		pages = append(pages, Page{
			Id:     cp.Id(OptionStream{Network: net}, "tremor-volcano-%n"),
			Height: strconv.Itoa(len(streams)*40 + 70),
			Png:    cp.Png(OptionStream{Network: net}, "/volcano/%n/tremor.png"),
			Plots: []Plot{
				Plot{
					Id:      "tremor",
					Clip:    "1",
					Height:  strconv.Itoa(len(streams) * 40),
					Length:  "25 days",
					Overlap: "0",
					Streams: streams,
					Title: &Title{
						Colour: "#006400",
						Font:   "LiberationSans Bold Italic 14",
					},
					XGrid: &XGrid{
						Colour: "#ffffff70",
						IsDate: "yes",
						Step:   "172800",
					},
					XLabel: &XLabel{
						Colour: "#006400",
						Font:   "LiberationSans Italic Bold 12",
						String: "Day of month",
					},
					XTick: &XTick{
						Colour:  "#006400",
						Factor:  "86400",
						Font:    "LiberationSans Bold 12",
						IsDate:  "yes",
						Step:    "172800",
						YOffset: "5",
						String:  "%d",
					},
					Border: &Border{
						Bg:       "#AFEEEE",
						Colour:   "#ffffff",
						Gradient: "#ffffffa0",
					},
					Copyright: &Copyright{
						Font: "LiberationSans Narrow 9",
					},
					Dates: []Date{
						Date{
							Font: "LiberationSans Narrow Bold 14",
						},
						Date{
							Aligned: "bottom",
							Font:    "LiberationSans Narrow Bold 14",
							String:  "%b %d, %Y",
							Time:    "start",
						},
						Date{
							Aligned: "bottom",
							Font:    "LiberationSans Narrow Bold 14",
							String:  "%b %d, %Y",
							Time:    "end",
						},
					},
					Label: &Label{
						Colour: "#006400",
						Font:   "LiberationSans Bold Italic 14",
						String: lookup.Description,
					},
					Name: &Name{
						Colour: "#006400",
						Font:   "LiberationSans Bold Italic 14",
					},
				},
			},
		})
	}

	return pages, nil
}
