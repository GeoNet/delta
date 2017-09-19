package main

import (
	"sort"
	"time"

	"github.com/GeoNet/delta/internal/metadb"
	"github.com/GeoNet/delta/meta"
)

func (cp ConfigPage) DrumSmall(base string) ([]Page, error) {

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
	for _, s := range cp.Options.Stations {
		list, err := db.Station(s)
		if err != nil {
			return nil, err
		}
		if list != nil {
			stns = append(stns, *list)
		}
	}

	sort.Sort(Stations(stns))

	var locations []string
	for _, l := range cp.Options.Locations {
		locations = append(locations, l)
	}
	if cp.Options.Reversed != "" {
		sort.Sort(sort.Reverse(sort.StringSlice(locations)))
	} else {
		sort.Strings(locations)
	}

	for _, s := range stns {
		chas, err := db.Channels(s.Code)
		if err != nil {
			return nil, err
		}

		for _, l := range locations {
			for _, x := range cp.Options.Channels() {
				for _, c := range chas {
					if c.Code != x {
						continue
					}
					if c.Location != l {
						continue
					}

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
						Id:     cp.Id(StationChannel(s, c), "%t-%n-%s-%l-s-drum"),
						Height: "263",
						Width:  "400",
						Png:    cp.Png(StationChannel(s, c), "/earthquake/drums/latest/%s-seismic-drum.png"),
						Plots: []Plot{
							Plot{
								Id:      cp.Id(StationChannel(s, c), "%t-%n-%s-%l-s-drum"),
								Clip:    "5",
								Height:  "242",
								Length:  "30 minutes",
								Overlap: "0",
								Width:   "400",
								X:       "0",
								Y:       "21",
								Missing: "Not operational",
								Streams: []Stream{
									Stream{
										Alt:    "137 104 205",
										Colour: "25 25 112",
										Format: "amplitude",
										Font:   "LiberationSans Bold 10",
										Auto:   cp.Auto(StationChannel(s, c)),
										Gain:   cp.Gain(StationChannel(s, c)),
										Row:    "4",
										Rrd:    cp.Rrd(StationChannel(s, c), "amplitude/drum/%s.%n/%s.%l-%c.%n.rrd"),
										Style:  "drum",
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
										XGrids: []XGrid{
											XGrid{
												Colour: "143 188 143",
												IsDate: "yes",
												Pen:    "0.5",
												Step:   "180",
											},
											XGrid{
												Colour: "32 178 170",
												IsDate: "yes",
												Pen:    "0.5",
												Step:   "1800",
											},
										},
									},
								},
								Border: &Border{
									Bg:     "#ffffff",
									Colour: "137 104 205",
									Pen:    "2",
								},
							},
						},
					})
				}
			}
		}
	}

	for _, s := range cp.Options.Streams {

		pages = append(pages, Page{
			Id:     cp.Id(s, "%t-%n-%s-%l-s-drum"),
			Height: "263",
			Width:  "400",
			Png:    cp.Png(s, "/earthquake/drums/latest/%s-seismic-drum.png"),
			Plots: []Plot{
				Plot{
					Id:      cp.Id(s, "%t-%n-%s-%l-s-drum"),
					Clip:    "5",
					Height:  "242",
					Length:  "30 minutes",
					Overlap: "0",
					Width:   "400",
					X:       "0",
					Y:       "21",
					Missing: "Not operational",
					Streams: []Stream{
						Stream{
							Alt:    "137 104 205",
							Colour: "25 25 112",
							Font:   "LiberationSans Bold 10",
							Format: "amplitude",
							Auto:   cp.Auto(s),
							Gain:   cp.Gain(s),
							Row:    "4",
							Rrd:    cp.Rrd(s, "amplitude/drum/%s.%n/%s.%l-%c.%n.rrd"),
							Style:  "drum",
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
							XGrids: []XGrid{
								XGrid{
									Colour: "143 188 143",
									IsDate: "yes",
									Pen:    "0.5",
									Step:   "180",
								},
								XGrid{
									Colour: "32 178 170",
									IsDate: "yes",
									Pen:    "0.5",
									Step:   "1800",
								},
							},
						},
					},
					Border: &Border{
						Bg:     "#ffffff",
						Colour: "137 104 205",
						Pen:    "2",
					},
				},
			},
		})
	}

	return pages, nil
}
