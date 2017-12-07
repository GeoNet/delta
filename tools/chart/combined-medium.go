package main

import (
	"sort"
	"time"

	"github.com/GeoNet/delta/internal/metadb"
	"github.com/GeoNet/delta/meta"
)

func (cp ConfigPage) CombinedMedium(base string) ([]Page, error) {

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

	for _, l := range cp.Options.Stations {
		s, err := db.Station(l)
		if err != nil {
			return nil, err
		}
		if s != nil {
			stns = append(stns, *s)
		}
	}

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
				Id:     cp.Id(StationChannel(s, c), "drum-volcano-%n-%s-%l-m-combined"),
				Height: "280",
				Width:  "400",
				Png:    cp.Png(StationChannel(s, c), "/volcano/drums/latest/m-%s-seismic-combined.png"),
				Plots: []Plot{
					{
						Id:      "rsam",
						Clip:    "1",
						Height:  "125",
						Length:  "35 days",
						Overlap: "0",
						Width:   "392",
						X:       "4",
						Y:       "14",
						Streams: []Stream{
							{
								Auto:   "yes",
								Colour: "#000000a0",
								Format: "amplitude",
								Rrd:    cp.Rrd(StationChannel(s, c), "/amplitude/rsam/%s.%n/%s.%l-%c.%n.rrd"),
								Style:  "rsam",
								Date: &Date{
									Colour: "25 25 112",
									Font:   "LiberationSans Narrow Bold 7",
								},
								Label: &Label{
									Colour: "25 25 112",
									Font:   "LiberationSans Bold 8",
									String: s.Code + "/" + c.Location + "-" + c.Code + "/" + s.Network,
								},
								Title: &Title{
									Colour: "25 25 112",
									Font:   "LiberationSans Bold Italic 8",
									String: s.Name,
								},
								XGrid: &XGrid{
									Colour: "143 188 143",
									IsDate: "yes",
									Pen:    "1",
									Step:   "172800",
								},
								YGrid: &YGrid{
									Colour: "143 188 143",
									Hint:   "4",
									Pen:    "1",
								},
							},
						},
						Border: &Border{
							Bg:     "#f5f5f5",
							Colour: "137 104 205",
							Pen:    "2",
						},
					},
					{
						Id:      "ssam",
						Clip:    "1",
						Height:  "125",
						Length:  "35 days",
						Overlap: "0",
						Width:   "392",
						X:       "4",
						Y:       "151",
						Streams: []Stream{
							{
								Auto:     "yes",
								Colour:   "#000000a0",
								Equalise: "no",
								Format:   "spectra",
								Map:      "rainbow",
								Reverse:  "no",
								Rrd:      cp.Rrd(StationChannel(s, c), "/spectra/ssam/%s.%n/%s.%l-%c.%n.rrd"),
								Style:    "ssam",
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

	for _, s := range cp.Options.Streams {
		pages = append(pages, Page{
			Id:     cp.Id(s, "drum-volcano-%n-%s-%l-m-combined"),
			Height: "280",
			Width:  "400",
			Png:    cp.Png(s, "/volcano/drums/latest/m-%s-seismic-combined.png"),
			Plots: []Plot{
				{
					Id:      "rsam",
					Clip:    "1",
					Height:  "125",
					Length:  "35 days",
					Overlap: "0",
					Width:   "392",
					X:       "4",
					Y:       "14",
					Streams: []Stream{
						{
							Auto:   "yes",
							Colour: "#000000a0",
							Format: "amplitude",
							Rrd:    cp.Rrd(s, "/amplitude/rsam/%s.%n/%s.%l-%c.%n.rrd"),
							Style:  "rsam",
							Date: &Date{
								Colour: "25 25 112",
								Font:   "LiberationSans Narrow Bold 7",
							},
							Label: &Label{
								Colour: "25 25 112",
								Font:   "LiberationSans Bold 8",
								String: s.Station + "/" + s.Location + "-" + s.Channel + "/" + s.Network,
							},
							Title: &Title{
								Colour: "25 25 112",
								Font:   "LiberationSans Bold Italic 8",
								String: s.Title,
							},
							XGrid: &XGrid{
								Colour: "143 188 143",
								IsDate: "yes",
								Pen:    "1",
								Step:   "172800",
							},
							YGrid: &YGrid{
								Colour: "143 188 143",
								Hint:   "4",
								Pen:    "1",
							},
						},
					},
					Border: &Border{
						Bg:     "#f5f5f5",
						Colour: "137 104 205",
						Pen:    "2",
					},
				},
				{
					Id:      "ssam",
					Clip:    "1",
					Height:  "125",
					Length:  "35 days",
					Overlap: "0",
					Width:   "392",
					X:       "4",
					Y:       "151",
					Streams: []Stream{
						{
							Auto:     "yes",
							Colour:   "#000000a0",
							Equalise: "no",
							Format:   "spectra",
							Map:      "rainbow",
							Reverse:  "no",
							Rrd:      cp.Rrd(s, "/spectra/ssam/%s.%n/%s.%l-%c.%n.rrd"),
							Style:    "ssam",
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

	return pages, nil
}
