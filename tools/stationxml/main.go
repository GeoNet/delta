package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/cmplx"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/GeoNet/delta/meta"
	"github.com/ozym/fdsn/stationxml"
)

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

	var level string
	flag.StringVar(&level, "level", "", "stationxml level")

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
		fmt.Fprintf(os.Stderr, "unable to compile station regexp %s: %v\n", stationRegexp, err)
		os.Exit(-1)
	}

	// which networks to process
	networkMatch, err := regexp.Compile(networkRegexp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to compile network regexp %s: %v\n", networkRegexp, err)
		os.Exit(-1)
	}

	// load network details ...
	networkMap := make(map[string]meta.Network)
	{
		var n meta.NetworkList
		if err := meta.LoadList(filepath.Join(network, "networks.csv"), &n); err != nil {
			fmt.Fprintf(os.Stderr, "unable to load network list: %v\n", err)
			os.Exit(-1)
		}

		for _, v := range n {
			networkMap[v.NetworkCode] = v
		}
	}

	// load station details
	stationMap := make(map[string]meta.Station)
	{
		var s meta.StationList
		if err := meta.LoadList(filepath.Join(network, "stations.csv"), &s); err != nil {
			fmt.Fprintf(os.Stderr, "unable to load station list: %v\n", err)
			os.Exit(-1)
		}

		for _, v := range s {
			stationMap[v.Code] = v
		}
	}

	// sorted station ids
	var keys []string
	for k, _ := range stationMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// load connections, i.e. what's connected to where
	connectionMap := make(map[string]meta.ConnectionList)
	{
		{
			var cons meta.ConnectionList
			if err := meta.LoadList(filepath.Join(install, "connections.csv"), &cons); err != nil {
				fmt.Fprintf(os.Stderr, "unable to load connection list: %v\n", err)
				os.Exit(-1)
			}

			for _, c := range cons {
				if _, ok := connectionMap[c.StationCode]; ok {
					connectionMap[c.StationCode] = append(connectionMap[c.StationCode], c)
				} else {
					connectionMap[c.StationCode] = meta.ConnectionList{c}
				}
			}

		}

		var recs meta.InstalledRecorderList
		if err := meta.LoadList(filepath.Join(install, "recorders.csv"), &recs); err != nil {
			fmt.Fprintf(os.Stderr, "unable to load recorder list: %v\n", err)
			os.Exit(-1)
		}
		for _, r := range recs {
			c := meta.Connection{
				StationCode:  r.StationCode,
				LocationCode: r.LocationCode,
				Span: meta.Span{
					Start: r.Start,
					End:   r.End,
				},
				Place: r.StationCode,
				Role:  r.LocationCode,
			}
			if _, ok := connectionMap[c.StationCode]; ok {
				connectionMap[c.StationCode] = append(connectionMap[c.StationCode], c)
			} else {
				connectionMap[c.StationCode] = meta.ConnectionList{c}
			}
		}
	}

	// where the sensors are installed
	siteMap := make(map[string]map[string]meta.Site)
	{
		var locs meta.SiteList
		if err := meta.LoadList(filepath.Join(network, "sites.csv"), &locs); err != nil {
			fmt.Fprintf(os.Stderr, "unable to load site list: %v\n", err)
			os.Exit(-1)
		}

		for _, l := range locs {
			if _, ok := siteMap[l.StationCode]; !ok {
				siteMap[l.StationCode] = make(map[string]meta.Site)
			}
			siteMap[l.StationCode][l.LocationCode] = l
		}
	}

	// where the dataloggers were deployed
	dataloggerDeploys := make(map[string]meta.DeployedDataloggerList)
	{
		var loggers meta.DeployedDataloggerList
		if err := meta.LoadList(filepath.Join(install, "dataloggers.csv"), &loggers); err != nil {
			fmt.Fprintf(os.Stderr, "unable to load datalogger list: %v\n", err)
			os.Exit(-1)
		}
		for _, d := range loggers {
			if _, ok := dataloggerDeploys[d.Place]; ok {
				dataloggerDeploys[d.Place] = append(dataloggerDeploys[d.Place], d)
			} else {
				dataloggerDeploys[d.Place] = meta.DeployedDataloggerList{d}
			}
		}

		var recs meta.InstalledRecorderList
		if err := meta.LoadList(filepath.Join(install, "recorders.csv"), &recs); err != nil {
			fmt.Fprintf(os.Stderr, "unable to load recorder list: %v\n", err)
			os.Exit(-1)
		}
		for _, r := range recs {
			d := meta.DeployedDatalogger{
				Install: meta.Install{
					Equipment: meta.Equipment{
						Make:   r.Make,
						Model:  r.DataloggerModel,
						Serial: r.Serial,
					},
					Span: meta.Span{
						Start: r.Start,
						End:   r.End,
					},
				},
				Place: r.StationCode,
				Role:  r.LocationCode,
			}
			if _, ok := dataloggerDeploys[d.Place]; ok {
				dataloggerDeploys[d.Place] = append(dataloggerDeploys[d.Place], d)
			} else {
				dataloggerDeploys[d.Place] = meta.DeployedDataloggerList{d}
			}
		}
	}

	// sort each datalogger deployment
	for i, _ := range dataloggerDeploys {
		sort.Sort(dataloggerDeploys[i])
	}

	// build sensor installation details
	sensorInstalls := make(map[string]meta.InstalledSensorList)
	{
		{
			var sensors meta.InstalledSensorList
			if err := meta.LoadList(filepath.Join(install, "sensors.csv"), &sensors); err != nil {
				fmt.Fprintf(os.Stderr, "unable to load sensor list: %v\n", err)
				os.Exit(-1)
			}
			for _, s := range sensors {
				if _, ok := sensorInstalls[s.StationCode]; ok {
					sensorInstalls[s.StationCode] = append(sensorInstalls[s.StationCode], s)
				} else {
					sensorInstalls[s.StationCode] = meta.InstalledSensorList{s}
				}
			}
		}
		{
			var recorders meta.InstalledRecorderList
			if err := meta.LoadList(filepath.Join(install, "recorders.csv"), &recorders); err != nil {
				fmt.Fprintf(os.Stderr, "unable to load recorder list: %v\n", err)
				os.Exit(-1)
			}
			for _, r := range recorders {
				if _, ok := sensorInstalls[r.StationCode]; ok {
					sensorInstalls[r.StationCode] = append(sensorInstalls[r.StationCode], r.InstalledSensor)
				} else {
					sensorInstalls[r.StationCode] = meta.InstalledSensorList{r.InstalledSensor}
				}
			}
		}
		{
			var gauges meta.InstalledGaugeList
			if err := meta.LoadList(filepath.Join(install, "gauges.csv"), &gauges); err != nil {
				fmt.Fprintf(os.Stderr, "unable to load gauge list: %v\n", err)
				os.Exit(-1)
			}
			for _, g := range gauges {
				s := meta.InstalledSensor{
					Install: meta.Install{
						Equipment: meta.Equipment{
							Make:   g.Make,
							Model:  g.Model,
							Serial: g.Serial,
						},
						Span: meta.Span{
							Start: g.Start,
							End:   g.End,
						},
					},
					Offset: meta.Offset{
						Height: g.Height,
						North:  g.North,
						East:   g.East,
					},
					Orientation: meta.Orientation{
						Dip:     g.Dip,
						Azimuth: g.Azimuth,
					},
					StationCode:  g.StationCode,
					LocationCode: g.LocationCode,
				}

				if _, ok := sensorInstalls[s.StationCode]; ok {
					sensorInstalls[s.StationCode] = append(sensorInstalls[s.StationCode], s)
				} else {
					sensorInstalls[s.StationCode] = meta.InstalledSensorList{s}
				}
			}
		}
	}

	// sort each sensor install
	for i, _ := range sensorInstalls {
		sort.Sort(sensorInstalls[i])
	}

	var networks []stationxml.Network

	// load response details
	resmap := make(map[string]map[string][]Stream)
	{
		for _, r := range Responses {
			for _, l := range r.Dataloggers {
				for _, a := range l.Dataloggers {
					if _, ok := resmap[a]; !ok {
						resmap[a] = make(map[string][]Stream)
					}
					for _, x := range r.Sensors {
						for _, b := range x.Sensors {
							resmap[a][b] = append(resmap[a][b], Stream{
								Datalogger: l,
								Sensor:     x,
							})
						}
					}
				}
			}
		}
	}

	components := make(map[string]map[int]SensorComponent)

	for k, v := range SensorModels {
		if _, ok := components[k]; !ok {
			components[k] = make(map[int]SensorComponent)
		}
		for n, p := range v.Components {
			components[k][n] = p
		}
	}

	stas := make(map[string][]stationxml.Station)
	for _, k := range keys {

		log.Printf("checking station: %s\n", k)
		if !stationMatch.MatchString(k) {
			log.Printf("skipping station: %s [outside regexp match]\n", k)
			continue
		}

		v := stationMap[k]
		if _, ok := siteMap[k]; !ok {
			log.Printf("skipping station: %s [no site map entry]\n", k)
			continue
		}
		n, ok := networkMap[v.Network]
		if !ok {
			fmt.Fprintf(os.Stderr, "error: unknown network %s\n", v.Network)
			os.Exit(-1)
		}

		if _, ok := connectionMap[k]; !ok {
			log.Printf("skipping station: %s [no connection map entry]\n", k)
			continue
		}

		if _, ok := sensorInstalls[k]; !ok {
			log.Printf("skipping station: %s [no sensor installs entry]\n", k)
			continue
		}

		var chans []stationxml.Channel
		if strings.ToLower(level) != "station" {

			for _, c := range connectionMap[k] {
				log.Printf("checking station channel: %s %s\n", k, c.LocationCode)

				if _, ok := dataloggerDeploys[c.Place]; !ok {
					log.Printf("skipping station channel: %s %s [no deployed datalogger]\n", k, c.Place)
					continue
				}
				l, ok := siteMap[k][c.LocationCode]
				if !ok {
					log.Printf("skipping station channel: %s %s [no site map]\n", k, c.LocationCode)
					continue
				}

				for _, s := range sensorInstalls[k] {
					log.Printf("checking station sensor install: %s %s\n", k, s.LocationCode)
					switch {
					case s.LocationCode != c.LocationCode:
						continue
					case s.Start.After(c.End):
						continue
					case s.End.Before(c.Start):
						continue
					}
					for _, d := range dataloggerDeploys[c.Place] {
						log.Printf("checking station datalogger deploys: %s %s %s\n", k, c.Place, d.Role)
						switch {
						case d.Role != c.Role:
							continue
						case d.Start.After(c.End):
							continue
						case d.End.Before(c.Start):
							continue
						}
						on := s.Start
						if d.Start.After(on) {
							on = d.Start
						}
						off := s.End
						if d.End.Before(off) {
							off = d.End
						}

						log.Printf("checking station datalogger response: %s %s\n", k, d.Model)
						if _, ok := resmap[d.Model]; !ok {
							log.Printf("skipping station datalogger response: %s \"%s\" no resmap\n", k, d.Model)
							continue
						}
						if _, ok := resmap[d.Model][s.Model]; !ok {
							log.Printf("skipping station datalogger response: %s \"%s\" \"%s\" no resmap\n", k, d.Model, s.Model)
							continue
						}
						log.Printf("checking station response: %s %s %s\n", k, d.Model, s.Model)
						for _, r := range resmap[d.Model][s.Model] {
							if r.Datalogger.Match != "" && !regexp.MustCompile(r.Datalogger.Match).MatchString(v.Code) {
								log.Println("skipping!", r.Datalogger.Match, v.Code)
								continue
							}
							if r.Datalogger.Skip != "" && regexp.MustCompile(r.Datalogger.Skip).MatchString(v.Code) {
								log.Println("skipping!", r.Datalogger.Skip, v.Code)
								continue
							}

							if r.Sensor.Match != "" && !regexp.MustCompile(r.Sensor.Match).MatchString(v.Code) {
								log.Println("skipping!", r.Sensor.Match, v.Code)
								continue
							}
							if r.Sensor.Skip != "" && regexp.MustCompile(r.Sensor.Skip).MatchString(v.Code) {
								log.Println("skipping!", r.Sensor.Skip, v.Code)
								continue
							}

							var types []stationxml.Type
							for _, j := range r.Type {
								switch j {
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
								continue
							}

							var f []string
							f = append(f, r.Sensor.Filters...)
							f = append(f, r.Datalogger.Filters...)

							for p, j := range r.Channels {
								if p != 0 {
									continue
								}
								c, ok := components[s.Model][p]
								if !ok {
									continue
								}

								azimuth := s.Azimuth + c.Azimuth
								for azimuth < 0.0 {
									azimuth += 360.0
								}
								for azimuth > 360.0 {
									azimuth -= 360.0
								}

								count := 1
								var stages []stationxml.ResponseStage
								for _, v := range f {
									if _, ok := Filters[v]; !ok {
										fmt.Fprintf(os.Stderr, "error: unknown filter %s\n", v)
										os.Exit(-1)
									}
									for _, s := range Filters[v] {
										switch s.Type {
										case "poly":
											p, ok := Polynomials[s.Lookup]
											if !ok {
												fmt.Fprintf(os.Stderr, "error: unknown polynomial filter %s\n", s.Lookup)
												os.Exit(-1)
											}
											var coeffs []stationxml.Coefficient
											for n, c := range p.Coefficients {
												coeffs = append(coeffs, stationxml.Coefficient{
													Number: uint32(n) + 1,
													Value:  c.Value,
												})
											}

											poly := stationxml.Polynomial{
												BaseFilter: stationxml.BaseFilter{
													ResourceId:  "Polynomial#" + v,
													Name:        fmt.Sprintf("%s.%s.%s%c.%04d.%03d.stage_%d", k, l.LocationCode, r.Label, j, on.Year(), on.YearDay(), count),
													InputUnits:  stationxml.Units{Name: s.InputUnits},
													OutputUnits: stationxml.Units{Name: s.OutputUnits},
												},
												ApproximationType: func() stationxml.ApproximationType {
													switch p.ApproximationType {
													case "MACLAURIN":
														return stationxml.ApproximationTypeMaclaurin
													default:
														return stationxml.ApproximationTypeUnknown
													}
												}(),
												FrequencyLowerBound:     stationxml.Frequency{stationxml.Float{Value: p.FrequencyLowerBound}},
												FrequencyUpperBound:     stationxml.Frequency{stationxml.Float{Value: p.FrequencyUpperBound}},
												ApproximationLowerBound: strconv.FormatFloat(p.ApproximationLowerBound, 'g', -1, 64),
												ApproximationUpperBound: strconv.FormatFloat(p.ApproximationUpperBound, 'g', -1, 64),
												MaximumError:            p.MaximumError,
												Coefficients:            coeffs,
											}
											stages = append(stages, stationxml.ResponseStage{
												Number:     stationxml.Counter(uint32(count)),
												Polynomial: &poly,
												StageGain: stationxml.Gain{
													Value: func() float64 {
														if p.Gain != 0.0 {
															return p.Gain
														}
														return 1.0
													}(),
													Frequency: s.Frequency,
												},
											})
										case "paz":
											if _, ok := PAZs[s.Lookup]; !ok {
												fmt.Fprintf(os.Stderr, "error: unknown paz filter %s\n", s.Lookup)
												os.Exit(-1)
											}
											var poles []stationxml.PoleZero
											for j, p := range PAZs[s.Lookup].Poles {
												poles = append(poles, stationxml.PoleZero{
													Number:    uint32(j),
													Real:      stationxml.FloatNoUnit{Value: real(p)},
													Imaginary: stationxml.FloatNoUnit{Value: imag(p)},
												})
											}
											var zeros []stationxml.PoleZero
											for j, z := range PAZs[s.Lookup].Zeros {
												zeros = append(zeros, stationxml.PoleZero{
													Number:    uint32(len(PAZs[s.Lookup].Poles) + j),
													Real:      stationxml.FloatNoUnit{Value: real(z)},
													Imaginary: stationxml.FloatNoUnit{Value: imag(z)},
												})
											}
											paz := stationxml.PolesZeros{
												BaseFilter: stationxml.BaseFilter{
													ResourceId:  "PolesZeros#" + v,
													Name:        fmt.Sprintf("%s.%s.%s%c.%04d.%03d.stage_%d", k, l.LocationCode, r.Label, j, on.Year(), on.YearDay(), count),
													InputUnits:  stationxml.Units{Name: s.InputUnits},
													OutputUnits: stationxml.Units{Name: s.OutputUnits},
												},
												PzTransferFunction: func() stationxml.PzTransferFunction {
													switch PAZs[s.Lookup].Code {
													case "A":
														return stationxml.PZFunctionLaplaceRadiansPerSecond
													case "B":
														return stationxml.PZFunctionLaplaceHertz
													case "D":
														return stationxml.PZFunctionLaplaceZTransform
													default:
														return stationxml.PZFunctionUnknown
													}
												}(),
												NormalizationFactor: func() float64 {
													w := complex(0.0, -2.0*math.Pi*s.Frequency)
													h := complex(float64(1.0), float64(0.0))
													for _, p := range PAZs[s.Lookup].Poles {
														switch PAZs[s.Lookup].Code {
														case "B":
															h /= (w - 2.0*math.Pi*p)
														default:
															h /= (w - p)
														}
													}
													for _, z := range PAZs[s.Lookup].Zeros {
														switch PAZs[s.Lookup].Code {
														case "B":
															h *= (w - 2.0*math.Pi*z)
														default:
															h *= (w - z)
														}
													}
													return 1.0 / cmplx.Abs(h)
												}(),

												NormalizationFrequency: stationxml.Frequency{
													stationxml.Float{Value: s.Frequency},
												},
												Zeros: zeros,
												Poles: poles,
											}
											stages = append(stages, stationxml.ResponseStage{
												Number:     stationxml.Counter(uint32(count)),
												PolesZeros: &paz,
												StageGain: stationxml.Gain{
													Value: func() float64 {
														if s.Gain != 0.0 {
															return s.Gain
														}
														return 1.0
													}(),
													Frequency: s.Frequency,
												},
											})

										case "a2d":
											coefs := stationxml.Coefficients{
												BaseFilter: stationxml.BaseFilter{
													ResourceId:  "Coefficients#" + v,
													Name:        fmt.Sprintf("%s.%s.%s%c.%04d.%03d.stage_%d", k, l.LocationCode, r.Label, j, on.Year(), on.YearDay(), count),
													InputUnits:  stationxml.Units{Name: s.InputUnits},
													OutputUnits: stationxml.Units{Name: s.OutputUnits},
												},
												CfTransferFunctionType: stationxml.CfFunctionDigital,
											}
											stages = append(stages, stationxml.ResponseStage{
												Number:       stationxml.Counter(uint32(count)),
												Coefficients: &coefs,
												StageGain: stationxml.Gain{
													Value: func() float64 {
														if s.Gain != 0.0 {
															return s.Gain
														}
														return 1.0
													}(),
													Frequency: s.Frequency,
												},
												Decimation: &stationxml.Decimation{
													InputSampleRate: stationxml.Frequency{stationxml.Float{Value: s.SampleRate}},
													Factor: func() int32 {
														if s.Decimate != 0 {
															return s.Decimate
														}
														return 1
													}(),
													Delay:      stationxml.Float{Value: s.Delay},
													Correction: stationxml.Float{Value: s.Correction},
												},
											})

										case "fir":
											if _, ok := FIRs[s.Lookup]; !ok {
												fmt.Fprintf(os.Stderr, "error: unknown fir filter %s\n", s.Lookup)
												os.Exit(-1)
											}
											var coeffs []stationxml.NumeratorCoefficient
											for j, c := range FIRs[s.Lookup].Factors {
												coeffs = append(coeffs, stationxml.NumeratorCoefficient{
													Coefficient: int32(j + 1),
													Value:       c,
												})
											}

											fir := stationxml.FIR{
												BaseFilter: stationxml.BaseFilter{
													ResourceId:  "FIR#" + v,
													Name:        fmt.Sprintf("%s.%s.%s%c.%04d.%03d.stage_%d", k, l.LocationCode, r.Label, j, on.Year(), on.YearDay(), count),
													InputUnits:  stationxml.Units{Name: s.InputUnits},
													OutputUnits: stationxml.Units{Name: s.OutputUnits},
												},
												Symmetry: func() stationxml.Symmetry {
													switch strings.ToUpper(FIRs[s.Lookup].Symmetry) {
													case "EVEN":
														return stationxml.SymmetryEven
													case "ODD":
														return stationxml.SymmetryOdd
													default:
														return stationxml.SymmetryNone
													}
												}(),
												NumeratorCoefficients: coeffs,
											}
											stages = append(stages, stationxml.ResponseStage{
												Number: stationxml.Counter(uint32(count)),
												FIR:    &fir,
												StageGain: stationxml.Gain{
													Value: func() float64 {
														if s.Gain != 0.0 {
															return s.Gain
														}
														return 1.0
													}(),
													Frequency: s.Frequency,
												},
												Decimation: &stationxml.Decimation{
													InputSampleRate: stationxml.Frequency{stationxml.Float{Value: FIRs[s.Lookup].Decimation * s.SampleRate}},
													Factor: func() int32 {
														if s.Decimate != 0 {
															return s.Decimate
														}
														return 1
													}(),
													Delay:      stationxml.Float{Value: s.Delay},
													Correction: stationxml.Float{Value: s.Correction},
												},
											})
										default:
											fmt.Fprintf(os.Stderr, "error: unknown filter type %s: %s\n", v, s.Type)
											os.Exit(-1)
										}
										count++
									}
								}

								chans = append(chans, stationxml.Channel{
									BaseNode: stationxml.BaseNode{
										Code:      r.Label + string(j),
										StartDate: &stationxml.DateTime{on},
										EndDate:   &stationxml.DateTime{off},
										RestrictedStatus: func() stationxml.RestrictedStatus {
											switch n.Restricted {
											case true:
												return stationxml.StatusClosed
											default:
												return stationxml.StatusOpen
											}
										}(),
									},
									LocationCode: l.LocationCode,
									Latitude: stationxml.Latitude{
										LatitudeBase: stationxml.LatitudeBase{
											Float: stationxml.Float{
												Value: l.Latitude,
											},
										},
									},
									Longitude: stationxml.Longitude{
										LongitudeBase: stationxml.LongitudeBase{
											Float: stationxml.Float{
												Value: l.Longitude,
											},
										},
									},
									Elevation: stationxml.Distance{Float: stationxml.Float{Value: l.Elevation}},
									Depth:     stationxml.Distance{Float: stationxml.Float{Value: -s.Height}},
									Azimuth:   &stationxml.Azimuth{Float: stationxml.Float{Value: azimuth}},
									Dip:       &stationxml.Dip{Float: stationxml.Float{Value: c.Dip}},
									Types:     types,
									SampleRateGroup: stationxml.SampleRateGroup{
										SampleRate: stationxml.SampleRate{Float: stationxml.Float{Value: r.Rate}},
										SampleRateRatio: func() *stationxml.SampleRateRatio {
											if r.Rate > 1.0 {
												return &stationxml.SampleRateRatio{
													NumberSamples: int32(r.Rate),
													NumberSeconds: 1,
												}
											} else {
												return &stationxml.SampleRateRatio{
													NumberSamples: 1,
													NumberSeconds: int32(1.0 / r.Rate),
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
											Gain: stationxml.Gain{
												Value: func() float64 {
													var gain float64 = 1.0
													for _, s := range stages {
														gain *= s.StageGain.Value
													}
													return gain
												}(),
												Frequency: func() float64 {
													var freq float64
													if len(stages) > 0 {
														freq = stages[0].StageGain.Frequency
													}
													return freq
												}(),
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
			Latitude:  v.Latitude,
			Longitude: v.Longitude,
		}

		stas[n.ExternalCode] = append(stas[n.ExternalCode], stationxml.Station{
			BaseNode: stationxml.BaseNode{
				Code:        v.Code,
				Description: n.Description,
				RestrictedStatus: func() stationxml.RestrictedStatus {
					switch n.Restricted {
					case true:
						return stationxml.StatusClosed
					default:
						return stationxml.StatusOpen
					}
				}(),
				StartDate: &(stationxml.DateTime{v.Start}),
				EndDate:   &(stationxml.DateTime{v.End}),
			},
			Latitude: stationxml.Latitude{LatitudeBase: stationxml.LatitudeBase{
				Float: stationxml.Float{
					Value: v.Latitude,
				}}, Datum: v.Datum},
			Longitude: stationxml.Longitude{LongitudeBase: stationxml.LongitudeBase{
				Float: stationxml.Float{
					Value: v.Longitude,
				}}, Datum: v.Datum},
			Elevation: stationxml.Distance{
				Float: stationxml.Float{Value: v.Elevation},
			},
			Site: stationxml.Site{
				Name: func() string {
					if v.Name != "" {
						return v.Name
					} else {
						return v.Code
					}
				}(),
				Description: func() string {
					if l := Locations.Closest(place); l != nil {
						if d := l.Distance(place); d < 5.0 {
							return fmt.Sprintf("within 5 km of %s", l.Name)
						} else {
							return fmt.Sprintf("%.0f km %s of %s", d, l.Compass(place), l.Name)
						}
					}
					return ""
				}(),
			},
			CreationDate: stationxml.DateTime{v.Start},
			TerminationDate: func() *stationxml.DateTime {
				if time.Now().Before(v.End) {
					return nil
				}
				return &stationxml.DateTime{v.End}
			}(),
			Channels: chans,
		})
	}

	for k, _ := range stas {
		n, ok := networkMap[k]
		if !ok {
			fmt.Fprintf(os.Stderr, "error: unknown network: %s\n", k)
			os.Exit(-1)
		}
		var on, off *stationxml.DateTime

		if !networkMatch.MatchString(n.ExternalCode) {
			log.Printf("skipping network: %s [outside regexp match]\n", n.ExternalCode)
			continue
		}

		for _, v := range stas[k] {
			if v.BaseNode.StartDate != nil {
				if on == nil || v.BaseNode.StartDate.Before(on.Time) {
					on = v.BaseNode.StartDate
				}
			}
			if v.BaseNode.EndDate != nil {
				if off == nil || v.BaseNode.EndDate.After(off.Time) {
					off = v.BaseNode.EndDate
				}
			}
		}
		networks = append(networks, stationxml.Network{
			BaseNode: stationxml.BaseNode{
				Code:        n.ExternalCode,
				Description: n.Description,
				RestrictedStatus: func() stationxml.RestrictedStatus {
					switch n.Restricted {
					case true:
						return stationxml.StatusClosed
					default:
						return stationxml.StatusOpen
					}
				}(),
				StartDate: on,
				EndDate:   off,
			},
			SelectedNumberStations: uint32(len(stas[k])),
			Stations:               stas[k],
		})
	}

	root := stationxml.NewFDSNStationXML(source, sender, module, "", networks)
	if ok := root.IsValid(); ok != nil {
		fmt.Fprintf(os.Stderr, "error: invalid stationxml file\n")
		os.Exit(-1)
	}

	// marshal into xml
	x, err := root.Marshal()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: unable to marshal stationxml\n")
		os.Exit(-1)
	}

	// output as needed ...
	if output != "-" {
		if err := os.MkdirAll(filepath.Dir(output), 0755); err != nil {
			fmt.Fprintf(os.Stderr, "error: unable to create directory %s: %v\n", filepath.Dir(output), err)
			os.Exit(-1)
		}
		if err := ioutil.WriteFile(output, x, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "error: unable to write file %s: %v\n", output, err)
			os.Exit(-1)
		}
	} else {
		fmt.Fprintln(os.Stdout, string(x))
	}
}
