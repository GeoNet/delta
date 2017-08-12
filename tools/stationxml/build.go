package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/GeoNet/delta/metadb"
	"github.com/GeoNet/delta/resp"

	"github.com/ozym/fdsn/stationxml"
)

type Channels []stationxml.Channel

func (c Channels) Len() int           { return len(c) }
func (c Channels) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c Channels) Less(i, j int) bool { return c[i].StartDate.Time.Before(c[j].StartDate.Time) }

func matcher(path, def string) (*regexp.Regexp, error) {

	// no list given
	if path == "" {
		return regexp.Compile(func() string {
			if def == "" {
				return "[A-Z0-9]+"
			}
			return def
		}())
	}

	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return regexp.Compile("^(" + strings.Join(strings.Fields(string(buf)), "|") + ")$")
}

type Builder struct {
	operational *time.Time
	networks    *regexp.Regexp
	stations    *regexp.Regexp
	channels    *regexp.Regexp
	sensors     *regexp.Regexp
	dataloggers *regexp.Regexp
	installed   bool
}

func SetInstalled(installed bool) func(*Builder) error {
	return func(b *Builder) error {
		b.installed = installed
		return nil
	}
}

func SetOperational(operational bool, offset time.Duration) func(*Builder) error {
	return func(b *Builder) error {
		t := time.Now().Add(-offset)
		if operational {
			b.operational = &t
		}
		return nil
	}
}

func SetNetworks(list, match string) func(*Builder) error {
	return func(b *Builder) error {
		re, err := matcher(list, match)
		if err != nil {
			return err
		}
		b.networks = re
		return nil
	}
}
func SetStations(list, match string) func(*Builder) error {
	return func(b *Builder) error {
		re, err := matcher(list, match)
		if err != nil {
			return err
		}
		b.stations = re
		return nil
	}
}
func SetChannels(list, match string) func(*Builder) error {
	return func(b *Builder) error {
		re, err := matcher(list, match)
		if err != nil {
			return err
		}
		b.channels = re
		return nil
	}
}
func SetSensors(list, match string) func(*Builder) error {
	return func(b *Builder) error {
		re, err := matcher(list, match)
		if err != nil {
			return err
		}
		b.sensors = re
		return nil
	}
}
func SetDataloggers(list, match string) func(*Builder) error {
	return func(b *Builder) error {
		re, err := matcher(list, match)
		if err != nil {
			return err
		}
		b.dataloggers = re
		return nil
	}
}

func NewBuilder(opts ...func(*Builder) error) (*Builder, error) {
	var b Builder

	for _, opt := range opts {
		if err := opt(&b); err != nil {
			return nil, err
		}
	}

	return &b, nil
}

func (b *Builder) MatchOperational(at time.Time) bool {
	if b.operational == nil {
		return true
	}
	return !(b.operational.After(at))
}
func (b *Builder) MatchNetwork(net string) bool {
	if b.networks == nil {
		return true
	}
	return b.networks.MatchString(net)
}
func (b *Builder) MatchStation(sta string) bool {
	if b.stations == nil {
		return true
	}
	return b.stations.MatchString(sta)
}
func (b *Builder) MatchChannel(cha string) bool {
	if b.channels == nil {
		return true
	}
	return b.channels.MatchString(cha)
}
func (b *Builder) MatchSensor(sensor string) bool {
	if b.sensors == nil {
		return true
	}
	return b.sensors.MatchString(sensor)
}
func (b *Builder) MatchDatalogger(logger string) bool {
	if b.dataloggers == nil {
		return true
	}
	return b.dataloggers.MatchString(logger)
}
func (b *Builder) Installed() bool {
	return b.installed
}

func (b *Builder) Construct(base string) ([]stationxml.Network, error) {

	mdb := metadb.NewMetaDB(base)

	var networks []stationxml.Network

	stations, err := mdb.Stations()
	if err != nil {
		return nil, err
	}

	stas := make(map[string][]stationxml.Station)
	for _, station := range stations {
		if !b.MatchStation(station.Code) {
			continue
		}
		network, err := mdb.Network(station.Network)
		if err != nil {
			return nil, err
		}
		if network == nil {
			continue
		}

		if !b.MatchNetwork(network.External) {
			continue
		}

		var channels []stationxml.Channel

		installations, err := mdb.Installations(station.Code)
		if err != nil {
			return nil, err
		}
		for _, installation := range installations {
			location, err := mdb.Site(station.Code, installation.Location)
			if err != nil {
				return nil, err
			}
			if location == nil {
				continue
			}
			if !b.MatchSensor(installation.Sensor.Model) {
				continue
			}
			if !b.MatchDatalogger(installation.Datalogger.Model) {
				continue
			}
			if !b.MatchOperational(installation.End) {
				continue
			}

			for _, response := range resp.Streams(installation.Datalogger.Model, installation.Sensor.Model) {
				stream, err := mdb.StationLocationSamplingRateStartStream(
					station.Code,
					installation.Location,
					response.Datalogger.SampleRate,
					installation.Start)
				if err != nil {
					return nil, err
				}
				if stream == nil {
					continue
				}

				var types []stationxml.Type
				for _, t := range response.Type {
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

				lookup := response.Channels(stream.Axial)
				for pin, comp := range response.Components {
					if !(pin < len(lookup)) {
						continue
					}
					if !b.MatchChannel(lookup[pin]) {
						continue
					}

					channel := lookup[pin]
					freq := response.Datalogger.Frequency
					dip := comp.Dip
					azimuth := installation.Sensor.Azimuth + comp.Azimuth

					// only rotate horizontal components
					if dip == 0.0 {
						if response.Sensor.Reversed {
							azimuth += 180.0
						}
						if response.Datalogger.Reversed {
							azimuth += 180.0
						}
						if stream.Reversed {
							azimuth += 180.0
						}
						// avoid negative zero
						dip = 0.0
						// bring into positive range
						for azimuth < 0.0 {
							azimuth += 360.0
						}
						for azimuth >= 360.0 {
							azimuth -= 360.0
						}
					} else {
						if response.Sensor.Reversed {
							dip *= -1.0
						}
						if response.Datalogger.Reversed {
							dip *= -1.0
						}
						if stream.Reversed {
							dip *= -1.0
						}
						// no azimuth on verticals
						azimuth = 0.0
					}

					tag := fmt.Sprintf(
						"%s.%s.%s",
						station.Code,
						location.Location,
						channel,
					)

					var stages []stationxml.ResponseStage
					for _, s := range append(response.Sensor.Stages, response.Datalogger.Stages...) {
						if s.StageSet == nil {
							continue
						}
						switch s.StageSet.GetType() {
						case "poly":
							stages = append(stages, polyResponseStage(s.StageSet.(resp.Polynomial), Stage{
								responseStage: s,
								count:         len(stages) + 1,
								id:            s.Filter,
								name: fmt.Sprintf(
									"%s.%04d.%03d.stage_%d",
									tag,
									installation.Start.Year(),
									installation.Start.YearDay(),
									len(stages)+1,
								),
								frequency: freq,
							}))
						case "paz":
							stages = append(stages, pazResponseStage(s.StageSet.(resp.PAZ), Stage{
								responseStage: s,
								count:         len(stages) + 1,
								id:            s.Filter,
								name: fmt.Sprintf(
									"%s.%04d.%03d.stage_%d",
									tag,
									installation.Start.Year(),
									installation.Start.YearDay(),
									len(stages)+1,
								),
								frequency: freq,
							}))
						case "a2d":
							stages = append(stages, a2dResponseStage(Stage{
								responseStage: s,
								count:         len(stages) + 1,
								id:            s.Filter,
								name: fmt.Sprintf(
									"%s.%04d.%03d.stage_%d",
									tag,
									installation.Start.Year(),
									installation.Start.YearDay(),
									len(stages)+1,
								),
								frequency: freq,
							}))
						case "fir":
							stages = append(stages, firResponseStage(s.StageSet.(resp.FIR), Stage{
								responseStage: s,
								count:         len(stages) + 1,
								id:            s.Filter,
								name: fmt.Sprintf(
									"%s.%04d.%03d.stage_%d",
									tag,
									installation.Start.Year(),
									installation.Start.YearDay(),
									len(stages)+1,
								),
								frequency: freq,
							}))
						}

					}

					channels = append(channels, stationxml.Channel{
						BaseNode: stationxml.BaseNode{
							Code:      channel, //response.Label + string(cha),
							StartDate: &stationxml.DateTime{installation.Start},
							EndDate:   &stationxml.DateTime{installation.End},
							RestrictedStatus: func() stationxml.RestrictedStatus {
								switch network.Restricted {
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
										switch location.Survey {
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
									Value: "Location is given in " + location.Datum,
								},
								stationxml.Comment{
									Id:    3,
									Value: "Sensor orientation not known",
								},
							},
						},
						LocationCode: location.Location,
						Latitude: stationxml.Latitude{
							LatitudeBase: stationxml.LatitudeBase{
								Float: stationxml.Float{
									Value: location.Latitude,
								},
							},
							Datum: location.Datum,
						},
						Longitude: stationxml.Longitude{
							LongitudeBase: stationxml.LongitudeBase{
								Float: stationxml.Float{
									Value: location.Longitude,
								},
							},
							Datum: location.Datum,
						},
						Elevation: stationxml.Distance{Float: stationxml.Float{Value: location.Elevation}},
						Depth:     stationxml.Distance{Float: stationxml.Float{Value: -installation.Sensor.Vertical}},
						Azimuth:   &stationxml.Azimuth{Float: stationxml.Float{Value: azimuth}},
						Dip:       &stationxml.Dip{Float: stationxml.Float{Value: dip}},
						Types:     types,
						SampleRateGroup: stationxml.SampleRateGroup{
							SampleRate: stationxml.SampleRate{Float: stationxml.Float{Value: response.SampleRate}},
							SampleRateRatio: func() *stationxml.SampleRateRatio {
								if response.SampleRate > 1.0 {
									return &stationxml.SampleRateRatio{
										NumberSamples: int32(response.SampleRate),
										NumberSeconds: 1,
									}
								} else {
									return &stationxml.SampleRateRatio{
										NumberSamples: 1,
										NumberSeconds: int32(1.0 / response.SampleRate),
									}
								}
							}(),
						},
						StorageFormat: response.StorageFormat,
						ClockDrift:    &stationxml.ClockDrift{Float: stationxml.Float{Value: response.ClockDrift}},
						Sensor: &stationxml.Equipment{
							ResourceId: "Sensor#" + installation.Sensor.Model + ":" + installation.Sensor.Serial,
							Type: func() string {
								if t, ok := resp.SensorModels[installation.Sensor.Model]; ok {
									return t.Type
								}
								return ""
							}(),
							Description: func() string {
								if t, ok := resp.SensorModels[installation.Sensor.Model]; ok {
									return t.Description
								}
								return ""
							}(),
							Manufacturer: func() string {
								if t, ok := resp.SensorModels[installation.Sensor.Model]; ok {
									return t.Manufacturer
								}
								return ""
							}(),
							Vendor: func() string {
								if t, ok := resp.SensorModels[installation.Sensor.Model]; ok {
									return t.Vendor
								}
								return ""
							}(),
							Model:        installation.Sensor.Model,
							SerialNumber: installation.Sensor.Serial,
							InstallationDate: func() *stationxml.DateTime {
								return &stationxml.DateTime{installation.Sensor.Start}
							}(),
							RemovalDate: func() *stationxml.DateTime {
								if time.Now().After(installation.Sensor.End) {
									return &stationxml.DateTime{installation.Sensor.End}
								}
								return nil
							}(),
						},

						DataLogger: &stationxml.Equipment{
							ResourceId: "Datalogger#" + installation.Datalogger.Model + ":" + installation.Datalogger.Serial,
							Type: func() string {
								if t, ok := resp.DataloggerModels[installation.Datalogger.Model]; ok {
									return t.Type
								}
								return ""
							}(),
							Description: func() string {
								if t, ok := resp.DataloggerModels[installation.Datalogger.Model]; ok {
									return t.Description
								}
								return ""
							}(),
							Manufacturer: func() string {
								if t, ok := resp.DataloggerModels[installation.Datalogger.Model]; ok {
									return t.Manufacturer
								}
								return ""
							}(),
							Vendor: func() string {
								if t, ok := resp.DataloggerModels[installation.Datalogger.Model]; ok {
									return t.Vendor
								}
								return ""
							}(),
							Model:        installation.Datalogger.Model,
							SerialNumber: installation.Datalogger.Serial,
							InstallationDate: func() *stationxml.DateTime {
								return &stationxml.DateTime{installation.Datalogger.Start}
							}(),
							RemovalDate: func() *stationxml.DateTime {
								if time.Now().After(installation.Datalogger.End) {
									return &stationxml.DateTime{installation.Datalogger.End}
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

		sort.Sort(Channels(channels))

		start, end := &(stationxml.DateTime{station.Start}), &(stationxml.DateTime{station.End})
		if b.Installed() && len(channels) > 0 {
			start, end = channels[0].StartDate, channels[len(channels)-1].EndDate
		}

		if !b.MatchOperational(end.Time) {
			continue
		}

		stas[network.External] = append(stas[network.External], stationxml.Station{
			BaseNode: stationxml.BaseNode{
				Code:        station.Code,
				Description: network.Description,
				RestrictedStatus: func() stationxml.RestrictedStatus {
					switch network.Restricted {
					case true:
						return stationxml.StatusClosed
					default:
						return stationxml.StatusOpen
					}
				}(),
				StartDate: start,
				EndDate:   end,
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
		network, err := mdb.Network(networkCode)
		if err != nil {
			return nil, err
		}
		if network == nil {
			continue
		}

		var start, end *stationxml.DateTime
		for _, s := range stationList {
			if s.BaseNode.StartDate != nil {
				if start == nil || s.BaseNode.StartDate.Before(start.Time) {
					start = s.BaseNode.StartDate
				}
			}
			if s.BaseNode.EndDate != nil {
				if end == nil || s.BaseNode.EndDate.After(end.Time) {
					end = s.BaseNode.EndDate
				}
			}
		}
		networks = append(networks, stationxml.Network{
			BaseNode: stationxml.BaseNode{
				Code:        network.External,
				Description: network.Description,
				RestrictedStatus: func() stationxml.RestrictedStatus {
					switch network.Restricted {
					case true:
						return stationxml.StatusClosed
					default:
						return stationxml.StatusOpen
					}
				}(),
				StartDate: start,
				EndDate:   end,
			},
			SelectedNumberStations: uint32(len(stationList)),
			Stations:               stationList,
		})
	}

	return networks, nil
}
