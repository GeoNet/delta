package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/GeoNet/delta/internal/stationxml/v1.2"

	"github.com/GeoNet/delta"
	"github.com/GeoNet/delta/meta"
)

const ClockDrift = 0.0001

func main() {

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Build a network StationXML file from delta meta & response information\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  %s [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fmt.Fprintf(os.Stderr, "\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")
	}

	var base string
	flag.StringVar(&base, "base", "", "base of delta files on disk")

	var created bool
	flag.BoolVar(&created, "created", false, "add a created value")

	var lookup string
	flag.StringVar(&lookup, "lookup", "", "base for response xml files on disk")

	var source string
	flag.StringVar(&source, "source", "GeoNet", "stationxml source")

	var sender string
	flag.StringVar(&sender, "sender", "WEL(GNS_Test)", "stationxml sender")

	var module string
	flag.StringVar(&module, "module", "Delta", "stationxml module")

	station := regexp.MustCompile("[A-Z0-9]+")
	flag.Func("station", "regexp selection of stations", func(s string) error {
		re, err := regexp.Compile(s)
		if err != nil {
			return err
		}
		station = re
		return nil
	})

	network := regexp.MustCompile("^(AK|CB|CH|CY|EC|FI|GM|HA|HB|KI|NM|NX|NZ|OT|PK|RT|SB|SC|SI|SM|SP|SX|TD|TG|TP|TR|WL|XX)$")
	flag.Func("network", "regexp selection of internal networks", func(s string) error {
		re, err := regexp.Compile(s)
		if err != nil {
			return err
		}
		network = re
		return nil
	})

	external := regexp.MustCompile("^(NZ)$")
	flag.Func("external", "regexp selection of external networks", func(s string) error {
		re, err := regexp.Compile(s)
		if err != nil {
			return err
		}
		external = re
		return nil
	})

	location := regexp.MustCompile("[A-Z0-9]+")
	flag.Func("location", "regexp selection of locations", func(s string) error {
		re, err := regexp.Compile(s)
		if err != nil {
			return err
		}
		location = re
		return nil
	})

	channel := regexp.MustCompile("[A-Z0-9]+")
	flag.Func("channel", "regexp selection of channels", func(s string) error {
		re, err := regexp.Compile(s)
		if err != nil {
			return err
		}
		channel = re
		return nil
	})

	var output string
	flag.StringVar(&output, "output", "", "output xml file, use \"-\" for stdout")

	flag.Parse()

	start := time.Now()

	set, err := delta.NewBase(base)
	if err != nil {
		log.Fatal(err)
	}

	places := meta.PlacenameList(set.Placenames())

	exts := make(map[string][]string)
	for _, n := range set.Networks() {
		if !external.MatchString(n.External) {
			continue
		}
		exts[n.External] = append(exts[n.External], n.Code)
	}

	nets := make(map[string]meta.Network)
	for _, n := range set.Networks() {
		if !network.MatchString(n.Code) {
			continue
		}
		if _, ok := exts[n.External]; !ok {
			continue
		}

		nets[n.Code] = n
	}

	stns := make(map[string]meta.Station)
	for _, s := range set.Stations() {
		if !station.MatchString(s.Code) {
			continue
		}
		if _, ok := nets[s.Network]; !ok {
			continue
		}

		stns[s.Code] = s
	}

	sites := make(map[string][]meta.Site)
	for _, s := range set.Sites() {
		if !location.MatchString(s.Location) {
			continue
		}
		if _, ok := stns[s.Station]; !ok {
			continue
		}

		sites[s.Station] = append(sites[s.Station], s)
	}

	resps, err := NewResponses(lookup, set)
	if err != nil {
		log.Fatal(err)
	}

	var networks []stationxml.NetworkType
	for n, lst := range exts {

		ext, ok := nets[n]
		if !ok {
			continue
		}

		var stations []stationxml.StationType
		for _, n := range lst {

			net, ok := nets[n]
			if !ok {
				continue
			}

			for _, stn := range stns {
				if stn.Network != n {
					continue
				}

				var channels []stationxml.ChannelType
				for _, site := range sites[stn.Code] {
					for _, c := range set.Collections(site) {
						for _, v := range set.Corrections(c) {

							dip, azimuth := c.Dip(v.Polarity), c.Azimuth(v.Polarity)

							var types []stationxml.Type
							switch {
							case c.Stream.Triggered:
								types = append(types, stationxml.Triggered)
							default:
								types = append(types, stationxml.Continuous)
							}
							for _, t := range c.Component.Types {
								switch t {
								case 'G':
									types = append(types, stationxml.Geophysical)
								case 'W':
									types = append(types, stationxml.Weather)
								case 'H':
									types = append(types, stationxml.Health)
								}
							}

							var depth float64
							if c.InstalledSensor.Vertical != 0.0 {
								depth = -c.InstalledSensor.Vertical
							}

							res, err := resps.Response(c, v)
							if err != nil {
								log.Fatal(err)
							}

							channels = append(channels, stationxml.ChannelType{
								BaseNodeType: stationxml.BaseNodeType{
									Code:      c.Code(),
									StartDate: stationxml.DateTime{Time: v.Span.Start},
									EndDate: func() stationxml.DateTime {
										if time.Since(v.Span.End) > 0 {
											return stationxml.DateTime{Time: v.Span.End}
										}
										return stationxml.DateTime{}
									}(),
									RestrictedStatus: func() stationxml.RestrictedStatusType {
										switch net.Restricted {
										case true:
											return stationxml.ClosedRestrictedStatus
										default:
											return stationxml.OpenRestrictedStatus
										}
									}(),
									Comment: []stationxml.CommentType{
										{
											Id: stationxml.CounterType(1),
											Value: func() string {
												switch site.Survey {
												case "Unknown":
													return "Location estimation method is unknown"
												case "External GPS Device":
													return "Location estimated from external GPS measurement"
												case "Internal GPS Clock":
													return "Location estimated from internal GPS clock"
												case "Topographic Map":
													return "Location estimated from topographic map"
												case "Site Survey":
													return "Location estimated from plans and survey to mark"
												default:
													return "Location estimation method is unknown"
												}
											}(),
										},
										{
											Id:    stationxml.CounterType(2),
											Value: "Location is given in " + site.Datum,
										},
										{
											Id:    stationxml.CounterType(3),
											Value: "Sensor orientation not known",
										},
									},
								},
								LocationCode: site.Location,
								Latitude: stationxml.LatitudeType{
									LatitudeBaseType: stationxml.LatitudeBaseType{
										FloatType: stationxml.FloatType{
											Value: site.Latitude,
										},
									},
									Datum: site.Datum,
								},
								Longitude: stationxml.LongitudeType{
									LongitudeBaseType: stationxml.LongitudeBaseType{
										FloatType: stationxml.FloatType{
											Value: site.Longitude,
										},
									},
									Datum: site.Datum,
								},
								Elevation:  stationxml.DistanceType{FloatType: stationxml.FloatType{Value: site.Elevation}},
								Depth:      stationxml.DistanceType{FloatType: stationxml.FloatType{Value: depth}},
								Azimuth:    &stationxml.AzimuthType{FloatType: stationxml.FloatType{Value: azimuth}},
								Dip:        &stationxml.DipType{FloatType: stationxml.FloatType{Value: dip}},
								Type:       types,
								SampleRate: stationxml.SampleRateType{FloatType: stationxml.FloatType{Value: c.Stream.SamplingRate}},
								SampleRateRatio: func() *stationxml.SampleRateRatioType {
									switch f := c.Stream.SamplingRate; {
									case f > 1.0:
										return &stationxml.SampleRateRatioType{
											NumberSamples: int(f),
											NumberSeconds: 1,
										}
									default:
										return &stationxml.SampleRateRatioType{
											NumberSamples: 1,
											NumberSeconds: int(1.0 / f),
										}
									}
								}(),
								ClockDrift: &stationxml.ClockDrift{FloatType: stationxml.FloatType{Value: ClockDrift}},
								Sensor: &stationxml.EquipmentType{
									ResourceId:   "Sensor#" + c.InstalledSensor.Model + ":" + c.InstalledSensor.Serial,
									Type:         c.Component.Type,
									Description:  strings.Fields(c.Component.Model)[0],
									Manufacturer: strings.Fields(c.Component.Make)[0],
									Model:        c.InstalledSensor.Model,
									SerialNumber: c.InstalledSensor.Serial,
									InstallationDate: stationxml.DateTime{
										Time: c.InstalledSensor.Start,
									},
									RemovalDate: func() stationxml.DateTime {
										if time.Since(c.InstalledSensor.End) > 0 {
											return stationxml.DateTime{Time: c.InstalledSensor.End}
										}
										return stationxml.DateTime{}
									}(),
								},

								DataLogger: &stationxml.EquipmentType{
									ResourceId:   "Datalogger#" + c.DeployedDatalogger.Model + ":" + c.DeployedDatalogger.Serial,
									Type:         c.Channel.Type,
									Description:  LegacyDescription(strings.Split(strings.Fields(c.Channel.Model)[0], "/")[0]),
									Manufacturer: strings.Fields(c.Channel.Make)[0],
									Model:        c.DeployedDatalogger.Model,
									SerialNumber: c.DeployedDatalogger.Serial,
									InstallationDate: func() stationxml.DateTime {
										return stationxml.DateTime{Time: c.DeployedDatalogger.Start}
									}(),
									RemovalDate: func() stationxml.DateTime {
										if time.Now().After(c.DeployedDatalogger.End) {
											return stationxml.DateTime{Time: c.DeployedDatalogger.End}
										}
										return stationxml.DateTime{}
									}(),
								},

								Response: res,
							})
						}
					}
				}

				sort.Slice(channels, func(i, j int) bool {
					switch {
					case channels[i].LocationCode < channels[j].LocationCode:
						return true
					case channels[i].LocationCode > channels[j].LocationCode:
						return false
					case channels[i].StartDate.Before(channels[j].StartDate.Time):
						return true
					case channels[i].StartDate.After(channels[j].StartDate.Time):
						return false
					case channels[i].SampleRate.Value > channels[j].SampleRate.Value:
						return true
					case channels[i].SampleRate.Value < channels[j].SampleRate.Value:
						return false
					case channels[i].Code < channels[j].Code:
						return true
					default:
						return false
					}
				})

				stations = append(stations, stationxml.StationType{
					BaseNodeType: stationxml.BaseNodeType{
						Code:        stn.Code,
						Description: net.Description,
						RestrictedStatus: func() stationxml.RestrictedStatusType {
							switch {
							case net.Restricted:
								return stationxml.ClosedRestrictedStatus
							default:
								return stationxml.OpenRestrictedStatus
							}
						}(),
						StartDate: stationxml.DateTime{Time: stn.Start},
						EndDate: func() stationxml.DateTime {
							if time.Since(stn.End) > 0 {
								return stationxml.DateTime{Time: stn.End}
							}
							return stationxml.DateTime{}
						}(),
						Comment: []stationxml.CommentType{
							{
								Id:    stationxml.CounterType(1),
								Value: "Location is given in " + stn.Datum,
							},
						},
					},
					Latitude: stationxml.LatitudeType{LatitudeBaseType: stationxml.LatitudeBaseType{
						FloatType: stationxml.FloatType{
							Value: stn.Latitude,
						}}, Datum: stn.Datum},
					Longitude: stationxml.LongitudeType{LongitudeBaseType: stationxml.LongitudeBaseType{
						FloatType: stationxml.FloatType{
							Value: stn.Longitude,
						}}, Datum: stn.Datum},
					Elevation: stationxml.DistanceType{
						FloatType: stationxml.FloatType{Value: stn.Elevation},
					},
					Site: stationxml.SiteType{
						Name: func() string {
							if stn.Name != "" {
								return stn.Name
							}
							return stn.Code
						}(),
						Description: places.Description(stn.Latitude, stn.Longitude),
					},
					CreationDate: stationxml.DateTime{Time: stn.Start},
					TerminationDate: func() stationxml.DateTime {
						if time.Since(stn.End) > 0 {
							return stationxml.DateTime{Time: stn.End}
						}
						return stationxml.DateTime{}
					}(),
					Channel: channels,
				})
			}
		}

		var start, end time.Time
		for _, s := range stns {
			if start.IsZero() || s.Start.Before(start) {
				start = s.Start
			}
			if end.IsZero() || s.End.After(end) {
				end = s.End
			}
		}

		sort.Slice(stations, func(i, j int) bool {
			return stations[i].BaseNodeType.Code < stations[j].BaseNodeType.Code
		})

		networks = append(networks, stationxml.NetworkType{
			BaseNodeType: stationxml.BaseNodeType{
				Code:        ext.Code,
				Description: ext.Description,
				RestrictedStatus: func() stationxml.RestrictedStatusType {
					switch ext.Restricted {
					case true:
						return stationxml.ClosedRestrictedStatus
					default:
						return stationxml.OpenRestrictedStatus
					}
				}(),
				StartDate: stationxml.DateTime{Time: start},
				EndDate: func() stationxml.DateTime {
					if time.Since(end) > 0 {
						return stationxml.DateTime{Time: end}
					}
					return stationxml.DateTime{}
				}(),
			},
			Station: stations,
		})
	}

	sort.Slice(networks, func(i, j int) bool {
		return networks[i].BaseNodeType.Code < networks[j].BaseNodeType.Code
	})

	root := stationxml.FDSNStationXML{
		RootType: stationxml.RootType{
			SchemaVersion: stationxml.SchemaVersion,
			Source:        source,
			Sender:        sender,
			Module:        module,
			Network:       networks,
			Created: func() *stationxml.DateTime {
				if created {
					n := stationxml.DateTime{Time: time.Now().UTC()}
					return &n
				}
				return nil
			}(),
		},
	}

	// marshal into xml
	res, err := root.MarshalIndent("", "  ")
	if err != nil {
		log.Fatalf("unable to marshal stationxml: %v", err)
	}

	switch output {
	case "", "-":
		fmt.Fprintln(os.Stdout, string(res))
	default:
		file, err := os.Create(output)
		if err != nil {
			log.Fatalf("error: unable to create file %s: %v", output, err)
		}
		defer file.Close()
		if _, err := file.Write(res); err != nil {
			log.Fatalf("error: unable to write file %s: %v", output, err)
		}
	}

	log.Println(time.Since(start))
}
