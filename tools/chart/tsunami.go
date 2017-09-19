package main

import (
	"sort"
	"time"

	"github.com/GeoNet/delta/internal/metadb"
)

func (cp ConfigPage) Tsunami(base string) ([]Page, error) {

	db := metadb.NewMetaDB(base)

	var pages []Page

	chans := make(map[string]interface{})
	for _, c := range cp.Options.Channels() {
		chans[c] = true
	}

	var plots []Plot

	for _, net := range cp.Options.Networks {
		var list []Stream

		stns, err := db.NetworkStation(net)
		if err != nil {
			return nil, err
		}

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
				if cp.Options.GetLocation(s.Code) != c.Location {
					continue
				}

				constituents, err := db.GaugeConstituents(s.Code)
				if err != nil {
					continue
				}
				var lag float64
				for _, t := range constituents {
					switch t.Name {
					case "M2":
						lag = t.Lag
					}
				}

				list = append(list, Stream{
					Auto:   "no",
					Colour: "#000000a0",
					Format: "amplitude",
					Rrd:    cp.Rrd(StationChannel(s, c), "amplitude/tsunami/%s.%n/%s.%l-%c.%n.rrd"),
					Style:  "gauge",
					Name: &Name{
						Box:     "#ffffffd0",
						Colour:  "#006400",
						Font:    "LiberationSans Bold Italic 11",
						Pad:     "2",
						String:  s.Name,
						XOffset: "10",
						YOffset: func() string {
							if cp.Options.Detide != 0 {
								return "-15"
							}
							return ""
						}(),
					},
					Tag: &Tag{
						Aligned: "first",
						Colour:  "#006400",
						Font:    "LiberationSans Narrow Bold 8",
						Rotated: "yes",
						String:  s.Code,
						XOffset: "5",
					},
					lag: lag,
				})
			}
		}

		sort.Sort(Streams(list))

		if len(list) > 0 {
			reference := func() int {
				switch {
				case cp.Options.Reference < 0:
					return (len(list) + int(cp.Options.Reference)) % len(list)
				default:
					return int(cp.Options.Reference) % len(list)
				}
			}()
			list[reference].Scalebar = &Scalebar{
				Length:  "1000",
				Stroke:  "2",
				Width:   "5",
				YOffset: "0",
				Scale: Scale{
					String: "one metre",
				},
			}
		}

		plots = append(plots, Plot{
			Id: func() string {
				if cp.Options.Detide != 0 {
					return "detide"
				}
				return "gauge"
			}(),
			Clip:    "100",
			Length:  "36 hours",
			Max:     "1500",
			Min:     "-1500",
			Overlap: "25",
			Width:   "795",
			Border: &Border{
				Bg:       "#90EE90",
				Colour:   "#ffffff",
				Gradient: "#ffffffa0",
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
				String: "New Zealand Tsunami Gauge Network",
			},
			Streams: list,
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
			XTick: &XTick{
				Colour:  "#006400",
				Factor:  "3600",
				Font:    "LiberationSans Bold 12",
				Step:    "21600",
				YOffset: "5",
			},
			YLabel: &YLabel{
				Colour: "#006400",
				Font:   "LiberationSans Italic Bold 12",
				String: "Relative sea level",
			},
		})
		pages = append(pages, Page{
			Id: cp.Id(OptionStream{
				Network: net,
			}, func() string {
				if cp.Options.Detide != 0 {
					return "tsunami-detide-detide-%n"
				}
				return "tsunami-gauge-gauge-%n"
			}()),
			Png: cp.Png(OptionStream{
				Network: net,
			}, func() string {
				if cp.Options.Detide != 0 {
					return "/tsunami/plots/latest/detide.png"
				}
				return "/tsunami/plots/latest/gauge.png"
			}()),
			Plots: plots,
		})
	}

	return pages, nil
}
