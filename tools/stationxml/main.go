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
	"github.com/GeoNet/delta/resp"
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

	var stationList string
	flag.StringVar(&stationList, "station-list", "", "regex selection of stations from file")

	var channelRegexp string
	flag.StringVar(&channelRegexp, "channels", "[A-Z0-9]+", "regex selection of channels")

	var channelList string
	flag.StringVar(&channelList, "channel-list", "", "regex selection of channels from file")

	var networkRegexp string
	flag.StringVar(&networkRegexp, "networks", "[A-Z0-9]+", "regex selection of networks")

	var networkList string
	flag.StringVar(&networkList, "network-list", "", "regex selection of networks from file")

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

	if stationList != "" {
		s, err := loadRegexpList(stationList)
		if err != nil || s == nil {
			log.Fatalf("unable to load station list regexp %s: %v", stationList, err)
		}
		stationRegexp = string(s)
	}

	// which stations to process
	stationMatch, err := regexp.Compile(stationRegexp)
	if err != nil {
		log.Fatalf("unable to compile station regexp %s: %v", stationRegexp, err)
	}

	if channelList != "" {
		s, err := loadRegexpList(channelList)
		if err != nil || s == nil {
			log.Fatalf("unable to load channel list regexp %s: %v", channelList, err)
		}
		channelRegexp = string(s)
	}

	// which stations to process
	channelMatch, err := regexp.Compile(channelRegexp)
	if err != nil {
		log.Fatalf("unable to compile channel regexp %s: %v", channelRegexp, err)
	}

	if networkList != "" {
		s, err := loadRegexpList(networkList)
		if err != nil || s == nil {
			log.Fatalf("unable to load network list regexp %s: %v", networkList, err)
		}
		networkRegexp = string(s)
	}

	// which networks to process
	networkMatch, err := regexp.Compile(networkRegexp)
	if err != nil {
		log.Fatalf("unable to compile network regexp %s: %v", networkRegexp, err)
	}

	// load meta information
	metaData, err := NewMeta(network, install)
	if err != nil {
		log.Fatalf("unable to load meta data: %v", err)
	}

	var networks []stationxml.Network

	stas := make(map[string][]stationxml.Station)
	for _, sta := range metaData.GetStationKeys() {
		station := metaData.GetStation(sta)
		if station == nil {
			continue
		}

		if !stationMatch.MatchString(sta) {
			continue
		}

		if verbose {
			log.Printf("checking station: %s", sta)
		}

		net := metaData.GetNetwork(station.Network)
		if net == nil {
			continue
		}

		if !networkMatch.MatchString(net.External) {
			continue
		}

		var channels []stationxml.Channel

		for _, conn := range metaData.GetConnections(sta) {
			if verbose {
				log.Printf("checking station channel: %s %s", sta, conn.Location)
			}

			if metaData.GetStreams(sta, conn.Location) == nil {
				continue
			}

			deploys := metaData.GetDeploys(conn.Place)
			if deploys == nil {
				log.Printf("skipping station channel: %s %s [no deployed datalogger]", sta, conn.Place)
				continue
			}

			l := metaData.GetSite(sta, conn.Location)
			if l == nil {
				log.Printf("skipping station channel: %s %s [no site map]", sta, conn.Location)
				continue
			}

			for _, sensorInstall := range metaData.GetInstalls(sta) {
				if verbose {
					log.Printf("checking station sensor install: %s %s", sta, sensorInstall.Location)
				}
				switch {
				case sensorInstall.Location != conn.Location:
					continue
				case sensorInstall.Start.After(conn.End):
					continue
				case sensorInstall.End.Before(conn.Start):
					continue
				case sensorInstall.Start == conn.End:
					continue
				}
				for _, dataloggerDeploy := range deploys {
					if verbose {
						log.Printf("checking station datalogger deploys: %s %s %s", sta, conn.Place, dataloggerDeploy.Role)
					}
					switch {
					case dataloggerDeploy.Role != conn.Role:
						continue
					case dataloggerDeploy.Start.After(conn.End):
						continue
					case dataloggerDeploy.End.Before(conn.Start):
						continue
					case dataloggerDeploy.Start == conn.End:
						continue
					case dataloggerDeploy.Start.After(sensorInstall.End):
						continue
					case dataloggerDeploy.End.Before(sensorInstall.Start):
						continue
					case dataloggerDeploy.Start == sensorInstall.End:
						continue
					case sensorInstall.End == dataloggerDeploy.Start:
						continue
					}

					on := conn.Start
					if sensorInstall.Start.After(on) {
						on = sensorInstall.Start
					}
					if dataloggerDeploy.Start.After(on) {
						on = dataloggerDeploy.Start
					}
					off := conn.End
					if sensorInstall.End.Before(off) {
						off = sensorInstall.End
					}
					if dataloggerDeploy.End.Before(off) {
						off = dataloggerDeploy.End
					}

					if verbose {
						log.Printf("checking station response: %s %s %s", sta, dataloggerDeploy.Model, sensorInstall.Model)
					}
					for _, r := range GetResponseStreams(dataloggerDeploy.Model, sensorInstall.Model) {
						stream := metaData.GetStream(sta, conn.Location, r.Datalogger.SampleRate, on)
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

						labels := r.Channels
						if stream.Axial {
							labels = strings.Replace(labels, "N", "1", -1)
							labels = strings.Replace(labels, "E", "2", -1)
						} else {
							labels = strings.Replace(labels, "1", "N", -1)
							labels = strings.Replace(labels, "2", "E", -1)
						}

						model, ok := resp.SensorModels[sensorInstall.Model]
						if !ok {
							continue
						}

						freq := r.Datalogger.Frequency
						for n := 0; n < len(labels) && n < len(model.Components); n++ {
							cha, comp := labels[n], model.Components[n]

							if !channelMatch.MatchString(r.Label + string(cha)) {
								continue
							}

							dip := comp.Dip
							azimuth := sensorInstall.Azimuth + comp.Azimuth

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

							tag := fmt.Sprintf("%s.%s.%s%c", sta, l.Location, r.Label, cha)

							var stages []stationxml.ResponseStage
							for _, s := range append(r.Sensor.Stages, r.Datalogger.Stages...) {
								if s.StageSet == nil {
									continue
								}
								switch s.StageSet.GetType() {
								case "poly":
									stages = append(stages, polyResponseStage(s.StageSet.(resp.Polynomial), Stage{
										responseStage: s,
										count:         len(stages) + 1,
										id:            s.Filter,
										name:          fmt.Sprintf("%s.%04d.%03d.stage_%d", tag, on.Year(), on.YearDay(), len(stages)+1),
										frequency:     freq,
									}))
								case "paz":
									stages = append(stages, pazResponseStage(s.StageSet.(resp.PAZ), Stage{
										responseStage: s,
										count:         len(stages) + 1,
										id:            s.Filter,
										name:          fmt.Sprintf("%s.%04d.%03d.stage_%d", tag, on.Year(), on.YearDay(), len(stages)+1),
										frequency:     freq,
									}))
								case "a2d":
									stages = append(stages, a2dResponseStage(Stage{
										responseStage: s,
										count:         len(stages) + 1,
										id:            s.Filter,
										name:          fmt.Sprintf("%s.%04d.%03d.stage_%d", tag, on.Year(), on.YearDay(), len(stages)+1),
										frequency:     freq,
									}))
								case "fir":
									stages = append(stages, firResponseStage(s.StageSet.(resp.FIR), Stage{
										responseStage: s,
										count:         len(stages) + 1,
										id:            s.Filter,
										name:          fmt.Sprintf("%s.%04d.%03d.stage_%d", tag, on.Year(), on.YearDay(), len(stages)+1),
										frequency:     freq,
									}))
								}
							}

							channels = append(channels, stationxml.Channel{
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
								LocationCode: l.Location,
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
								Depth:     stationxml.Distance{Float: stationxml.Float{Value: -sensorInstall.Vertical}},
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
									ResourceId: "Sensor#" + sensorInstall.Model + ":" + sensorInstall.Serial,
									Type: func() string {
										if t, ok := resp.SensorModels[sensorInstall.Model]; ok {
											return t.Type
										}
										return ""
									}(),
									Description: func() string {
										if t, ok := resp.SensorModels[sensorInstall.Model]; ok {
											return t.Description
										}
										return ""
									}(),
									Manufacturer: func() string {
										if t, ok := resp.SensorModels[sensorInstall.Model]; ok {
											return t.Manufacturer
										}
										return ""
									}(),
									Vendor: func() string {
										if t, ok := resp.SensorModels[sensorInstall.Model]; ok {
											return t.Vendor
										}
										return ""
									}(),
									Model:        sensorInstall.Model,
									SerialNumber: sensorInstall.Serial,
									InstallationDate: func() *stationxml.DateTime {
										return &stationxml.DateTime{sensorInstall.Start}
									}(),
									RemovalDate: func() *stationxml.DateTime {
										if time.Now().After(sensorInstall.End) {
											return &stationxml.DateTime{sensorInstall.End}
										}
										return nil
									}(),
								},

								DataLogger: &stationxml.Equipment{
									ResourceId: "Datalogger#" + dataloggerDeploy.Model + ":" + dataloggerDeploy.Serial,
									Type: func() string {
										if t, ok := resp.DataloggerModels[dataloggerDeploy.Model]; ok {
											return t.Type
										}
										return ""
									}(),
									Description: func() string {
										if t, ok := resp.DataloggerModels[dataloggerDeploy.Model]; ok {
											return t.Description
										}
										return ""
									}(),
									Manufacturer: func() string {
										if t, ok := resp.DataloggerModels[dataloggerDeploy.Model]; ok {
											return t.Manufacturer
										}
										return ""
									}(),
									Vendor: func() string {
										if t, ok := resp.DataloggerModels[dataloggerDeploy.Model]; ok {
											return t.Vendor
										}
										return ""
									}(),
									Model:        dataloggerDeploy.Model,
									SerialNumber: dataloggerDeploy.Serial,
									InstallationDate: func() *stationxml.DateTime {
										return &stationxml.DateTime{dataloggerDeploy.Start}
									}(),
									RemovalDate: func() *stationxml.DateTime {
										if time.Now().After(dataloggerDeploy.End) {
											return &stationxml.DateTime{dataloggerDeploy.End}
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

		sort.Sort(Channels(channels))

		stas[net.External] = append(stas[net.External], stationxml.Station{
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
					place := Place{
						Latitude:  station.Latitude,
						Longitude: station.Longitude,
					}
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
			Channels: channels,
		})
	}

	for networkCode, stationList := range stas {
		var net *meta.Network
		if net = metaData.GetNetwork(networkCode); net == nil {
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
				Code:        net.External,
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
