package main

import (
	"sort"
	"time"

	"github.com/GeoNet/delta/internal/metadb"
	"github.com/GeoNet/delta/meta"
)

func (cp ConfigPage) Combined(base string) ([]Page, error) {

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
		stns = append(stns, list...)
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
				Id:     cp.Id(StationChannel(s, c), "drum-volcano-%n-%s-%l-combined"),
				Height: "525",
				Width:  "750",
				Png:    cp.Png(StationChannel(s, c), "/volcano/drums/latest/%s-seismic-combined.png"),
				Plots: []Plot{
					{
						Id:     "dummy",
						Height: "220",
						Length: "35 days",
						Width:  "645",
						X:      "40",
						Y:      "260",
						Streams: []Stream{
							{
								Rrd: cp.Rrd(StationChannel(s, c), ""),
							},
						},
						Date: &Date{
							Aligned: "bottom",
							Colour:  "25 25 112",
							Font:    "LiberationSans Narrow Bold 11",
							String:  "%b %d, %Y",
							Time:    "end",
						},
					},
					{
						Id:      "rsam",
						Clip:    "1",
						Height:  "220",
						Length:  "35 days",
						Overlap: "0",
						Width:   "645",
						X:       "40",
						Y:       "20",
						Streams: []Stream{
							{
								Auto:   "yes",
								Colour: "#000000a0",
								Format: "amplitude",
								Rrd:    cp.Rrd(StationChannel(s, c), "/amplitude/rsam/%s.%n/%s.%l-%c.%n.rrd"),
								Style:  "rsam",
								Date: &Date{
									Colour: "25 25 112",
									Font:   "LiberationSans Narrow Bold 12",
								},
								Label: &Label{
									Colour: "25 25 112",
									Font:   "LiberationSans Bold 12",
									String: s.Code + "/" + c.Location + "-" + c.Code + "/" + s.Network,
								},
								Scale: &Scale{
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
									Hint:   "4",
									Pen:    "0.5",
								},
								YLabel: &YLabel{
									Colour: "25 25 112",
									Font:   "LiberationSans Bold Italic 11",
									String: "Amplitude Level",
								},
								YTick: &YTick{
									Colour:  "25 25 112",
									Font:    "LiberationSans Bold 11",
									Hint:    "4",
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
					{
						Id:      "ssam",
						Clip:    "1",
						Height:  "220",
						Length:  "35 days",
						Overlap: "0",
						Width:   "645",
						X:       "40",
						Y:       "260",
						Streams: []Stream{
							{
								Auto:     "yes",
								Colour:   "#000000a0",
								Equalise: "no",
								Format:   "spectra",
								Map:      "rainbow",
								Reverse:  "no",
								Rrd:      cp.Rrd(StationChannel(s, c), "/spectra/ssam/%s.%n/%s.%l-%c.%n.rrd"),
								ColourMap: &ColourMap{
									Border: &Border{
										Colour: "137 104 205",
										Pen:    "1",
									},
									Title: &Title{
										Colour: "25 25 112",
										Font:   "LiberationSans Narrow Bold 12",
										String: "dB",
									},
									YGrid: &YGrid{
										Colour: "#ffffff",
										Pen:    "2",
										Step:   "20",
									},
									YTick: &YTick{
										Aligned: "right",
										Colour:  "25 25 112",
										Font:    "LiberationSans Narrow Bold 11",
										Factor:  "1.0",
										Rotated: "no",
										Step:    "20",
										String:  "%2.0f",
										XOffset: "5",
									},
								},
								Style: "ssam",
								Date: &Date{
									Aligned: "bottom",
									Colour:  "25 25 112",
									Font:    "LiberationSans Narrow Bold 11",
									String:  "%b %d, %Y",
									Time:    "start",
								},
								XLabel: &XLabel{
									Colour: "25 25 112",
									Font:   "LiberationSans Bold Italic 11",
									String: "Day of month",
								},
								YLabel: &YLabel{
									Colour: "25 25 112",
									Font:   "LiberationSans Bold Italic 11",
									String: "Frequency (Hz)",
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
								YTick: &YTick{
									Colour:  "25 25 112",
									Font:    "LiberationSans Bold 11",
									Step:    "5",
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

	for _, s := range cp.Options.Streams {
		lookup, err := db.Station(s.Station)
		if err != nil {
			return nil, err
		}

		pages = append(pages, Page{
			Id:     cp.Id(s, "drum-volcano-%n-%s-%l-combined"),
			Height: "525",
			Width:  "750",
			Png:    cp.Png(s, "/volcano/drums/latest/%s-seismic-combined.png"),
			Plots: []Plot{
				{
					Id:     "dummy",
					Height: "220",
					Length: "35 days",
					Width:  "645",
					X:      "40",
					Y:      "260",
					Streams: []Stream{
						{
							Rrd: cp.Rrd(s, ""),
						},
					},
					Date: &Date{
						Aligned: "bottom",
						Colour:  "25 25 112",
						Font:    "LiberationSans Narrow Bold 11",
						String:  "%b %d, %Y",
						Time:    "end",
					},
				},
				{
					Id:      "rsam",
					Clip:    "1",
					Height:  "220",
					Length:  "35 days",
					Overlap: "0",
					Width:   "645",
					X:       "40",
					Y:       "20",
					Streams: []Stream{
						{
							Auto:   "yes",
							Colour: "#000000a0",
							Format: "amplitude",
							Rrd:    cp.Rrd(s, "/amplitude/rsam/%s.%n/%s.%l-%c.%n.rrd"),
							Style:  "rsam",
							Date: &Date{
								Colour: "25 25 112",
								Font:   "LiberationSans Narrow Bold 12",
							},
							Label: &Label{
								Colour: "25 25 112",
								Font:   "LiberationSans Bold 12",
								String: s.Station + "/" + s.Location + "-" + s.Channel + "/" + s.Network,
							},
							Scale: &Scale{
								Aligned: "top",
								Colour:  "25 25 112",
								Font:    "LiberationSans Bold 11",
							},
							Title: &Title{
								Colour: "25 25 112",
								Font:   "LiberationSans Bold Italic 14",
								String: lookup.Name,
							},
							XGrid: &XGrid{
								Colour: "143 188 143",
								IsDate: "yes",
								Pen:    "0.5",
								Step:   "172800",
							},
							YGrid: &YGrid{
								Colour: "143 188 143",
								Hint:   "4",
								Pen:    "0.5",
							},
							YLabel: &YLabel{
								Colour: "25 25 112",
								Font:   "LiberationSans Bold Italic 11",
								String: "Amplitude Level",
							},
							YTick: &YTick{
								Colour:  "25 25 112",
								Font:    "LiberationSans Bold 11",
								Hint:    "4",
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
				{
					Id:      "ssam",
					Clip:    "1",
					Height:  "220",
					Length:  "35 days",
					Overlap: "0",
					Width:   "645",
					X:       "40",
					Y:       "260",
					Streams: []Stream{
						{
							Auto:     "yes",
							Colour:   "#000000a0",
							Equalise: "no",
							Format:   "spectra",
							Map:      "rainbow",
							Reverse:  "no",
							Rrd:      cp.Rrd(s, "/spectra/ssam/%s.%n/%s.%l-%c.%n.rrd"),
							ColourMap: &ColourMap{
								Border: &Border{
									Colour: "137 104 205",
									Pen:    "1",
								},
								Title: &Title{
									Colour: "25 25 112",
									Font:   "LiberationSans Narrow Bold 12",
									String: "dB",
								},
								YGrid: &YGrid{
									Colour: "#ffffff",
									Pen:    "2",
									Step:   "20",
								},
								YTick: &YTick{
									Aligned: "right",
									Colour:  "25 25 112",
									Font:    "LiberationSans Narrow Bold 11",
									Factor:  "1.0",
									Rotated: "no",
									Step:    "20",
									String:  "%2.0f",
									XOffset: "5",
								},
							},
							Style: "ssam",
							Date: &Date{
								Aligned: "bottom",
								Colour:  "25 25 112",
								Font:    "LiberationSans Narrow Bold 11",
								String:  "%b %d, %Y",
								Time:    "start",
							},
							XLabel: &XLabel{
								Colour: "25 25 112",
								Font:   "LiberationSans Bold Italic 11",
								String: "Day of month",
							},
							YLabel: &YLabel{
								Colour: "25 25 112",
								Font:   "LiberationSans Bold Italic 11",
								String: "Frequency (Hz)",
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
							YTick: &YTick{
								Colour:  "25 25 112",
								Font:    "LiberationSans Bold 11",
								Step:    "5",
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

	return pages, nil
}
