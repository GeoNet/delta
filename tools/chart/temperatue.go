package main

import (
	"strings"
)

func (cp ConfigPage) Temperature(base string) ([]Page, error) {

	var pages []Page

	for _, s := range cp.Options.Streams {
		pages = append(pages, Page{
			Id:     cp.Id(s, "drum-volcano-%n-%s-%l-drum"),
			Height: "505",
			Width:  "770",
			Png:    cp.Png(s, "/volcano/%n/%s/%l/drum.png"),
			Plots: []Plot{
				{
					Id:      "rsam",
					Clip:    "1",
					Height:  "440",
					Length:  "35 days",
					Overlap: "0",
					Width:   "700",
					X:       "40",
					Streams: []Stream{
						{
							Auto:   "false",
							Gain:   "50.0",
							Colour: "#000000a0",
							Format: "amplitude",
							Rrd:    cp.Rrd(s, "amplitude/temperature/%s.%n/%s.%l-%c.%n.rrd"),
							Style:  "rsam",
							Copyright: &Copyright{
								Aligned: "top",
								Colour:  "25 25 112",
								Font:    "LiberationSans Bold 11",
							},
							Date: &Date{
								Colour: "25 25 112",
								Font:   "LiberationSans Narrow Bold 12",
							},
							Title: &Title{
								Colour: "25 25 112",
								Font:   "LiberationSans Bold Italic 14",
								String: s.Title,
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
								String: "Temperature (&#xb0;C)",
							},
							XTick: &XTick{
								Colour:  "25 25 112",
								Factor:  "86400",
								IsDate:  "yes",
								Font:    "LiberationSans Bold 11",
								Step:    "172800",
								YOffset: "5",
								String:  "%d",
							},
							YTick: &YTick{
								Colour:  "25 25 112",
								Font:    "LiberationSans Bold 11",
								Hint:    "10",
								Rotated: "yes",
								String:  "%g ",
							},
							Label: &Label{
								Colour: "25 25 112",
								Font:   "LiberationSans Bold 12",
								String: strings.Join([]string{s.Station, strings.Join([]string{s.Location, s.Channel}, "-"), s.Network}, "/"),
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
