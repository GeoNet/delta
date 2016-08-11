package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/GeoNet/delta/meta"
	"github.com/ozym/fdsn/stationxml"
)

type Channels []stationxml.Channel

func (c Channels) Len() int           { return len(c) }
func (c Channels) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c Channels) Less(i, j int) bool { return c[i].StartDate.Time.Before(c[j].StartDate.Time) }

func main() {

	var verbose bool
	flag.BoolVar(&verbose, "verbose", false, "make noise")

	var source string
	flag.StringVar(&source, "source", "GeoNet", "stationxml source")

	var sender string
	flag.StringVar(&sender, "sender", "WEL(GNS_Test)", "stationxml sender")

	var module string
	flag.StringVar(&module, "module", "Delta", "stationxml module")

	var output string
	flag.StringVar(&output, "output", "-", "output xml file")

	var network string
	flag.StringVar(&network, "network", "../../network", "base network directory")

	var install string
	flag.StringVar(&install, "install", "../../install", "base installs directory")

	var stationRegexp string
	flag.StringVar(&stationRegexp, "stations", "[A-Z0-9]+", "regex selection of stations")

	var networkRegexp string
	flag.StringVar(&networkRegexp, "networks", "[A-Z0-9]+", "regex selection of networks")

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

	flag.Parse()

	// which stations to process
	stationMatch, err := regexp.Compile(stationRegexp)
	if err != nil {
		log.Fatalf("unable to compile station regexp %s: %v", stationRegexp, err)
	}

	// which networks to process
	networkMatch, err := regexp.Compile(networkRegexp)
	if err != nil {
		log.Fatalf("unable to compile network regexp %s: %v", networkRegexp, err)
	}

	// load network map from data files
	networkMap, err := NetworkMap(network)
	if err != nil {
		log.Fatalf("unable to load network map: %v", err)
	}

	// load station map from data files
	stationMap, err := StationMap(network)
	if err != nil {
		log.Fatalf("unable to load station map: %v", err)
	}

	// sorted station ids
	var keys []string
	for k, _ := range stationMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// load connection map from data files
	connectionMap, err := ConnectionMap(install)
	if err != nil {
		log.Fatalf("unable to load connection map: %v", err)
	}

	// load network map from data files
	streamMap, err := StreamMap(install)
	if err != nil {
		log.Fatalf("unable to load stream map: %v", err)
	}

	// load site map from data files
	siteMap, err := SiteMap(network)
	if err != nil {
		log.Fatalf("unable to load site map: %v", err)
	}

	// load deployed datalogger map from data files
	dataloggerDeploys, err := DataloggerDeploys(install)
	if err != nil {
		log.Fatalf("unable to load deployed datalogger list: %v", err)
	}

	// load installed sensor map from data files
	sensorInstalls, err := SensorInstalls(install)
	if err != nil {
		log.Fatalf("unable to load installed sensor list: %v", err)
	}

	// load response details
	resmap := ResponseMap()

	// load sensor component details
	components := Components()

	var networks []stationxml.Network

	stas := make(map[string][]stationxml.Station)
	for _, sta := range keys {
		station := stationMap[sta]

		if !stationMatch.MatchString(sta) {
			continue
		}

		if verbose {
			log.Printf("checking station: %s", sta)
		}

		if _, ok := siteMap[sta]; !ok {
			log.Printf("skipping station: %s [no site map entry]", sta)
			continue
		}
		net, ok := networkMap[station.Network]
		if !ok {
			log.Fatalf("error: unknown network %s", station.Network)
			os.Exit(-1)
		}

		if _, ok := connectionMap[sta]; !ok {
			log.Printf("skipping station: %s [no connection map entry]", sta)
			continue
		}

		if _, ok := sensorInstalls[sta]; !ok {
			log.Printf("skipping station: %s [no sensor installs entry]", sta)
			continue
		}

		var chans []stationxml.Channel

		if _, ok := streamMap[sta]; ok {
			for _, conn := range connectionMap[sta] {
				if verbose {
					log.Printf("checking station channel: %s %s", sta, conn.LocationCode)
				}

				streams, ok := streamMap[sta][conn.LocationCode]
				if !ok {
					log.Printf("skipping station channel: %s %s [no stream map]", sta, conn.LocationCode)
					continue
				}

				if _, ok := dataloggerDeploys[conn.Place]; !ok {
					log.Printf("skipping station channel: %s %s [no deployed datalogger]", sta, conn.Place)
					continue
				}
				l, ok := siteMap[sta][conn.LocationCode]
				if !ok {
					log.Printf("skipping station channel: %s %s [no site map]", sta, conn.LocationCode)
					continue
				}

				for _, s := range sensorInstalls[sta] {
					if verbose {
						log.Printf("checking station sensor install: %s %s", sta, s.LocationCode)
					}
					switch {
					case s.LocationCode != conn.LocationCode:
						continue
					case s.Start.After(conn.End):
						continue
					case s.End.Before(conn.Start):
						continue
					case s.Start == conn.End:
						continue
					}
					for _, d := range dataloggerDeploys[conn.Place] {
						if verbose {
							log.Printf("checking station datalogger deploys: %s %s %s", sta, conn.Place, d.Role)
						}
						switch {
						case d.Role != conn.Role:
							continue
						case d.Start.After(conn.End):
							continue
						case d.End.Before(conn.Start):
							continue
						case d.Start == conn.End:
							continue
						case d.Start.After(s.End):
							continue
						case d.End.Before(s.Start):
							continue
						case d.Start == s.End:
							continue
						case s.End == d.Start:
							continue
						}

						on := conn.Start
						if s.Start.After(on) {
							on = s.Start
						}
						if d.Start.After(on) {
							on = d.Start
						}
						off := conn.End
						if s.End.Before(off) {
							off = s.End
						}
						if d.End.Before(off) {
							off = d.End
						}

						if verbose {
							log.Printf("checking station datalogger response: %s %s", sta, d.Model)
						}
						if _, ok := resmap[d.Model]; !ok {
							log.Printf("skipping station datalogger response: %s \"%s\" no resmap", sta, d.Model)
							continue
						}
						if _, ok := resmap[d.Model][s.Model]; !ok {
							log.Printf("skipping station datalogger response: %s \"%s\" \"%s\" no resmap", sta, d.Model, s.Model)
							continue
						}
						if verbose {
							log.Printf("checking station response: %s %s %s", sta, d.Model, s.Model)
						}
						for _, r := range resmap[d.Model][s.Model] {
							var stream *meta.Stream
							for _, s := range streams {
								if s.SamplingRate != r.Datalogger.SampleRate {
									continue
								}
								if s.End.Before(on) {
									continue
								}
								if s.Start.After(on) {
									break
								}
								stream = &s
							}
							if stream == nil {
								continue
							}

							var types []stationxml.Type
							for _, t := range r.Type {
								switch t {
								case 'c', 'C':
									types = append(types, stationxml.TypeContinuous)
								case 't', 'T':
									types = append(types, stationxml.TypeTriggered)
								case 'g', 'G':
									types = append(types, stationxml.TypeGeophysical)
								case 'w', 'W':
									types = append(types, stationxml.TypeWeather)
								}
							}

							if _, ok := components[s.Model]; !ok {
								log.Printf("no components found, skipping: %s", s.Model)
								continue
							}

							labels := r.Channels
							if stream.Axial {
								labels = strings.Replace(labels, "N", "1", -1)
								labels = strings.Replace(labels, "E", "2", -1)
							} else {
								labels = strings.Replace(labels, "1", "N", -1)
								labels = strings.Replace(labels, "2", "E", -1)
							}

							freq := r.Datalogger.Frequency
							for p, cha := range labels {
								comp, ok := components[s.Model][p]
								if !ok {
									continue
								}

								dip := comp.Dip
								azimuth := s.Azimuth + comp.Azimuth

								// only rotate horizontal components
								if dip == 0.0 {
									if r.Sensor.Reversed {
										azimuth += 180.0
									}
									if r.Datalogger.Reversed {
										azimuth += 180.0
									}
									if stream.Reversed {
										azimuth += 180.0
									}
									// avoid negative zero
									dip = 0.0
								} else {
									if r.Sensor.Reversed {
										dip *= -1.0
									}
									if r.Datalogger.Reversed {
										dip *= -1.0
									}
									if stream.Reversed {
										dip *= -1.0
									}
								}

								// bring into positive range
								for azimuth < 0.0 {
									azimuth += 360.0
								}
								for azimuth > 360.0 {
									azimuth -= 360.0
								}

								tag := fmt.Sprintf("%s.%s.%s%c", sta, l.LocationCode, r.Label, cha)

								var stages []stationxml.ResponseStage
								for _, s := range append(r.Sensor.Stages, r.Datalogger.Stages...) {
									if s.StageSet == nil {
										continue
									}
									stages = append(stages, s.StageSet.ResponseStage(Stage{
										responseStage: s,
										count:         len(stages) + 1,
										id:            s.Filter,
										name:          fmt.Sprintf("%s.%04d.%03d.stage_%d", tag, on.Year(), on.YearDay(), len(stages)+1),
										frequency:     freq,
									}))
								}

								chans = append(chans, stationxml.Channel{
									BaseNode: stationxml.BaseNode{
										Code:      r.Label + string(cha),
										StartDate: &stationxml.DateTime{on},
										EndDate:   &stationxml.DateTime{off},
										RestrictedStatus: func() stationxml.RestrictedStatus {
											switch net.Restricted {
											case true:
												return stationxml.StatusClosed
											default:
												return stationxml.StatusOpen
											}
										}(),
										Comments: []stationxml.Comment{
											stationxml.Comment{
												Id: 1,
												Value: func() string {
													switch l.Survey {
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
											stationxml.Comment{
												Id:    2,
												Value: "Location is given in " + l.Datum,
											},
											stationxml.Comment{
												Id:    3,
												Value: "Sensor orientation not known",
											},
										},
									},
									LocationCode: l.LocationCode,
									Latitude: stationxml.Latitude{
										LatitudeBase: stationxml.LatitudeBase{
											Float: stationxml.Float{
												Value: l.Latitude,
											},
										},
										Datum: l.Datum,
									},
									Longitude: stationxml.Longitude{
										LongitudeBase: stationxml.LongitudeBase{
											Float: stationxml.Float{
												Value: l.Longitude,
											},
										},
										Datum: l.Datum,
									},
									Elevation: stationxml.Distance{Float: stationxml.Float{Value: l.Elevation}},
									Depth:     stationxml.Distance{Float: stationxml.Float{Value: -s.Height}},
									Azimuth:   &stationxml.Azimuth{Float: stationxml.Float{Value: azimuth}},
									Dip:       &stationxml.Dip{Float: stationxml.Float{Value: dip}},
									Types:     types,
									SampleRateGroup: stationxml.SampleRateGroup{
										SampleRate: stationxml.SampleRate{Float: stationxml.Float{Value: r.SampleRate}},
										SampleRateRatio: func() *stationxml.SampleRateRatio {
											if r.SampleRate > 1.0 {
												return &stationxml.SampleRateRatio{
													NumberSamples: int32(r.SampleRate),
													NumberSeconds: 1,
												}
											} else {
												return &stationxml.SampleRateRatio{
													NumberSamples: 1,
													NumberSeconds: int32(1.0 / r.SampleRate),
												}
											}
										}(),
									},
									StorageFormat: r.StorageFormat,
									ClockDrift:    &stationxml.ClockDrift{Float: stationxml.Float{Value: r.ClockDrift}},
									Sensor: &stationxml.Equipment{
										ResourceId: "Sensor#" + s.Model + ":" + s.Serial,
										Type: func() string {
											if t, ok := SensorModels[s.Model]; ok {
												return t.Type
											}
											return ""
										}(),
										Description: func() string {
											if t, ok := SensorModels[s.Model]; ok {
												return t.Description
											}
											return ""
										}(),
										Manufacturer: func() string {
											if t, ok := SensorModels[s.Model]; ok {
												return t.Manufacturer
											}
											return ""
										}(),
										Vendor: func() string {
											if t, ok := SensorModels[s.Model]; ok {
												return t.Vendor
											}
											return ""
										}(),
										Model:        s.Model,
										SerialNumber: s.Serial,
										InstallationDate: func() *stationxml.DateTime {
											return &stationxml.DateTime{s.Start}
										}(),
										RemovalDate: func() *stationxml.DateTime {
											if time.Now().After(s.End) {
												return &stationxml.DateTime{s.End}
											}
											return nil
										}(),
									},

									DataLogger: &stationxml.Equipment{
										ResourceId: "Datalogger#" + d.Model + ":" + d.Serial,
										Type: func() string {
											if t, ok := DataloggerModels[d.Model]; ok {
												return t.Type
											}
											return ""
										}(),
										Description: func() string {
											if t, ok := DataloggerModels[d.Model]; ok {
												return t.Description
											}
											return ""
										}(),
										Manufacturer: func() string {
											if t, ok := DataloggerModels[d.Model]; ok {
												return t.Manufacturer
											}
											return ""
										}(),
										Vendor: func() string {
											if t, ok := DataloggerModels[d.Model]; ok {
												return t.Vendor
											}
											return ""
										}(),
										Model:        d.Model,
										SerialNumber: d.Serial,
										InstallationDate: func() *stationxml.DateTime {
											return &stationxml.DateTime{d.Start}
										}(),
										RemovalDate: func() *stationxml.DateTime {
											if time.Now().After(d.End) {
												return &stationxml.DateTime{d.End}
											}
											return nil
										}(),
									},
									Response: &stationxml.Response{
										Stages: stages,
										InstrumentSensitivity: &stationxml.Sensitivity{
											//TODO: check we may need to adjust gain for different frequency
											Gain: stationxml.Gain{
												Value: func() float64 {
													var gain float64 = 1.0
													for _, s := range stages {
														gain *= s.StageGain.Value
													}
													return gain
												}(),
												Frequency: freq,
											},
											InputUnits: func() stationxml.Units {
												var units stationxml.Units
												if len(stages) > 0 {
													s := stages[0]
													switch {
													case s.PolesZeros != nil:
														units = s.PolesZeros.BaseFilter.InputUnits
													case s.Coefficients != nil:
														units = s.Coefficients.BaseFilter.InputUnits
													case s.ResponseList != nil:
														units = s.ResponseList.BaseFilter.InputUnits
													case s.FIR != nil:
														units = s.FIR.BaseFilter.InputUnits
													case s.Polynomial != nil:
														units = s.Polynomial.BaseFilter.InputUnits
													}
												}
												return units
											}(),
											OutputUnits: func() stationxml.Units {
												var units stationxml.Units
												if len(stages) > 0 {
													s := stages[len(stages)-1]
													switch {
													case s.PolesZeros != nil:
														units = s.PolesZeros.BaseFilter.OutputUnits
													case s.Coefficients != nil:
														units = s.Coefficients.BaseFilter.OutputUnits
													case s.ResponseList != nil:
														units = s.ResponseList.BaseFilter.OutputUnits
													case s.FIR != nil:
														units = s.FIR.BaseFilter.OutputUnits
													case s.Polynomial != nil:
														units = s.Polynomial.BaseFilter.OutputUnits
													}
												}
												return units
											}(),
										},
									},
								})

							}
						}
					}
				}
			}
		}

		place := Place{
			Latitude:  station.Latitude,
			Longitude: station.Longitude,
		}

		sort.Sort(Channels(chans))

		stas[net.ExternalCode] = append(stas[net.ExternalCode], stationxml.Station{
			BaseNode: stationxml.BaseNode{
				Code:        station.Code,
				Description: net.Description,
				RestrictedStatus: func() stationxml.RestrictedStatus {
					switch net.Restricted {
					case true:
						return stationxml.StatusClosed
					default:
						return stationxml.StatusOpen
					}
				}(),
				StartDate: &(stationxml.DateTime{station.Start}),
				EndDate:   &(stationxml.DateTime{station.End}),
				Comments: []stationxml.Comment{
					stationxml.Comment{
						Id:    1,
						Value: "Location is given in " + station.Datum,
					},
				},
			},
			Latitude: stationxml.Latitude{LatitudeBase: stationxml.LatitudeBase{
				Float: stationxml.Float{
					Value: station.Latitude,
				}}, Datum: station.Datum},
			Longitude: stationxml.Longitude{LongitudeBase: stationxml.LongitudeBase{
				Float: stationxml.Float{
					Value: station.Longitude,
				}}, Datum: station.Datum},
			Elevation: stationxml.Distance{
				Float: stationxml.Float{Value: station.Elevation},
			},
			Site: stationxml.Site{
				Name: func() string {
					if station.Name != "" {
						return station.Name
					} else {
						return station.Code
					}
				}(),
				Description: func() string {
					if loc := Locations.Closest(place); loc != nil {
						if dist := loc.Distance(place); dist < 5.0 {
							return fmt.Sprintf("within 5 km of %s", loc.Name)
						} else {
							return fmt.Sprintf("%.0f km %s of %s", dist, loc.Compass(place), loc.Name)
						}
					}
					return ""
				}(),
			},
			CreationDate: stationxml.DateTime{station.Start},
			TerminationDate: func() *stationxml.DateTime {
				if time.Now().Before(station.End) {
					return nil
				}
				return &stationxml.DateTime{station.End}
			}(),
			Channels: chans,
		})
	}

	for networkCode, stationList := range stas {
		net, ok := networkMap[networkCode]
		if !ok {
			log.Fatalf("error: unknown network: %s", networkCode)
		}
		if !networkMatch.MatchString(net.ExternalCode) {
			if verbose {
				log.Printf("skipping network: %s [outside regexp match]", net.ExternalCode)
			}
			continue
		}

		var on, off *stationxml.DateTime
		for _, s := range stationList {
			if s.BaseNode.StartDate != nil {
				if on == nil || s.BaseNode.StartDate.Before(on.Time) {
					on = s.BaseNode.StartDate
				}
			}
			if s.BaseNode.EndDate != nil {
				if off == nil || s.BaseNode.EndDate.After(off.Time) {
					off = s.BaseNode.EndDate
				}
			}
		}
		networks = append(networks, stationxml.Network{
			BaseNode: stationxml.BaseNode{
				Code:        net.ExternalCode,
				Description: net.Description,
				RestrictedStatus: func() stationxml.RestrictedStatus {
					switch net.Restricted {
					case true:
						return stationxml.StatusClosed
					default:
						return stationxml.StatusOpen
					}
				}(),
				StartDate: on,
				EndDate:   off,
			},
			SelectedNumberStations: uint32(len(stationList)),
			Stations:               stationList,
		})
	}

	// render station xml
	root := stationxml.NewFDSNStationXML(source, sender, module, "", networks)
	if ok := root.IsValid(); ok != nil {
		log.Fatal("error: invalid stationxml file")
	}

	// marshal into xml
	x, err := root.Marshal()
	if err != nil {
		log.Fatalf("error: unable to marshal stationxml: %v", err)
	}

	// output as needed ...
	switch output {
	case "-":
		fmt.Fprintln(os.Stdout, string(x))
	default:
		if err := os.MkdirAll(filepath.Dir(output), 0755); err != nil {
			log.Fatalf("error: unable to create directory %s: %v", filepath.Dir(output), err)
		}
		if err := ioutil.WriteFile(output, x, 0644); err != nil {
			log.Fatalf("error: unable to write file %s: %v", output, err)
		}
	}
}
