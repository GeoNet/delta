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

	"github.com/GeoNet/delta/fdsn"
	"github.com/GeoNet/delta/meta"
	"github.com/GeoNet/delta/resp"
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
	flag.StringVar(&network, "network", "../network", "base network directory")

	var installs string
	flag.StringVar(&installs, "installs", "../installs", "base installs directory")

	var config string
	flag.StringVar(&config, "config", "../config", "base config directory")

	var responses string
	flag.StringVar(&responses, "responses", "../responses", "base response directory")

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

	stationMatch, err := regexp.Compile(stationRegexp)
	if err != nil {
		panic(err)
	}

	networkMatch, err := regexp.Compile(networkRegexp)
	if err != nil {
		panic(err)
	}

	// load response paz mappings ...
	pazMap := make(map[string]resp.PAZ)
	{
		paz, err := resp.LoadPAZFile(filepath.Join(responses, "paz.toml"))
		if err != nil {
			panic(err)

		}
		for _, p := range paz {
			pazMap[p.Name] = p
		}
	}

	// load response fir mappings ...
	firMap := make(map[string]resp.FIR)
	{
		fir, err := resp.LoadFIRFile(filepath.Join(responses, "fir.toml"))
		if err != nil {
			panic(err)

		}
		for _, f := range fir {
			firMap[f.Name] = f
		}
	}

	// load response paz mappings ...
	polyMap := make(map[string]resp.Polynomial)
	{
		poly, err := resp.LoadPolynomialFile(filepath.Join(responses, "polynomial.toml"))
		if err != nil {
			panic(err)

		}
		for _, p := range poly {
			polyMap[p.Name] = p
		}
	}

	// load network details ...
	networkMap := make(map[string]meta.Network)
	{
		var n meta.NetworkList
		if err := meta.LoadList(filepath.Join(network, "networks.csv"), &n); err != nil {
			panic(err)
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
			panic(err)
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
			if err := meta.LoadList(filepath.Join(config, "connections.csv"), &cons); err != nil {
				panic(err)
			}

			for _, c := range cons {
				if _, ok := connectionMap[c.StationCode]; ok {
					connectionMap[c.StationCode] = append(connectionMap[c.StationCode], c)
				} else {
					connectionMap[c.StationCode] = meta.ConnectionList{c}
				}
			}

		}
		{
			var cons meta.ConnectionList
			if err := meta.LoadList(filepath.Join(config, "tsunami.csv"), &cons); err != nil {
				panic(err)
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
		if err := meta.LoadList(filepath.Join(installs, "recorders.csv"), &recs); err != nil {
			panic(err)
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
			panic(err)
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
		if err := meta.LoadList(filepath.Join(installs, "dataloggers.csv"), &loggers); err != nil {
			panic(err)
		}
		for _, d := range loggers {
			if _, ok := dataloggerDeploys[d.Place]; ok {
				dataloggerDeploys[d.Place] = append(dataloggerDeploys[d.Place], d)
			} else {
				dataloggerDeploys[d.Place] = meta.DeployedDataloggerList{d}
			}
		}

		var recs meta.InstalledRecorderList
		if err := meta.LoadList(filepath.Join(installs, "recorders.csv"), &recs); err != nil {
			panic(err)
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
			if err := meta.LoadList(filepath.Join(installs, "sensors.csv"), &sensors); err != nil {
				panic(err)
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
			if err := meta.LoadList(filepath.Join(installs, "recorders.csv"), &recorders); err != nil {
				panic(err)
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
			if err := meta.LoadList(filepath.Join(installs, "gauges.csv"), &gauges); err != nil {
				panic(err)
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

	var networks []fdsn.Network

	// load response details
	resmap := make(map[string]map[string][]resp.Stream)
	{
		resps, err := resp.LoadResponseFile(filepath.Join(responses, "response.toml"))
		if err != nil {
			panic(err)
		}

		for _, r := range resps {
			for _, l := range r.Dataloggers {
				for _, a := range l.Dataloggers {
					if _, ok := resmap[a]; !ok {
						resmap[a] = make(map[string][]resp.Stream)
					}
					for _, x := range r.Sensors {
						for _, b := range x.Sensors {
							resmap[a][b] = append(resmap[a][b], resp.Stream{
								Datalogger: l,
								Sensor:     x,
							})
						}
					}
				}
			}
		}
	}
	for a, v := range resmap {
		for b, _ := range v {
			fmt.Println(a, b)
		}
	}

	// load datalogger details
	dataloggers := make(map[string]resp.DataloggerModel)
	{
		d, err := resp.LoadDataloggerModelFile(filepath.Join(responses, "datalogger.toml"))
		if err != nil {
			panic(err)
		}

		for _, v := range d {
			dataloggers[v.Model] = v
		}
	}

	// load sensor & component details
	sensors := make(map[string]resp.SensorModel)
	components := make(map[string]map[int]resp.SensorComponent)
	{
		s, err := resp.LoadSensorModelFile(filepath.Join(responses, "sensor.toml"))
		if err != nil {
			panic(err)
		}

		for _, v := range s {
			if _, ok := components[v.Model]; !ok {
				components[v.Model] = make(map[int]resp.SensorComponent)
			}
			for n, p := range v.Pins {
				components[v.Model][n] = p
			}
			sensors[v.Model] = v
		}
	}

	// load component details
	filters := make(map[string]resp.Filter)
	{
		f, err := resp.LoadFilterFile(filepath.Join(responses, "filter.toml"))
		if err != nil {
			panic(err)
		}

		for _, v := range f {
			filters[v.Name] = v
		}
	}

	stas := make(map[string][]fdsn.Station)
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
			panic("unknown network " + v.Network)
		}

		if _, ok := connectionMap[k]; !ok {
			log.Printf("skipping station: %s [no connection map entry]\n", k)
			continue
		}

		if _, ok := sensorInstalls[k]; !ok {
			log.Printf("skipping station: %s [no sensor installs entry]\n", k)
			continue
		}

		var chans []fdsn.Channel
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
							if r.Match != "" && !regexp.MustCompile(r.Match).MatchString(v.Code) {
								log.Println("skipping!", r.Match, v.Code)
								continue
							}
							if r.Skip != "" && regexp.MustCompile(r.Skip).MatchString(v.Code) {
								log.Println("skipping!", r.Skip, v.Code)
								continue
							}

							var types []fdsn.Type
							for _, j := range r.Type {
								switch j {
								case 'c', 'C':
									types = append(types, fdsn.TypeContinuous)
								case 't', 'T':
									types = append(types, fdsn.TypeTriggered)
								case 'g', 'G':
									types = append(types, fdsn.TypeGeophysical)
								case 'w', 'W':
									types = append(types, fdsn.TypeWeather)
								}
							}

							if _, ok := components[s.Model]; !ok {
								continue
							}

							var f []string
							f = append(f, r.SNFilters...)
							f = append(f, r.DLFilters...)

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
								var stages []fdsn.ResponseStage
								for _, v := range f {
									if _, ok := filters[v]; !ok {
										panic("missing filter: " + v)
									}
									for _, s := range filters[v].Stages {
										switch s.Type {
										case "poly":
											p, ok := polyMap[s.Lookup]
											if !ok {
												panic("unknown poly filter: " + s.Lookup)
											}
											var coeffs []fdsn.Coefficient
											for n, c := range p.Coefficients {
												coeffs = append(coeffs, fdsn.Coefficient{
													Number: uint32(n) + 1,
													Value:  c.Value,
												})
											}

											poly := fdsn.Polynomial{
												BaseFilter: fdsn.BaseFilter{
													ResourceId:  "Polynomial#" + v,
													Name:        fmt.Sprintf("%s.%s.%s%c.%04d.%03d.stage_%d", k, l.LocationCode, r.Label, j, on.Year(), on.YearDay(), count),
													InputUnits:  fdsn.Units{Name: s.InputUnits},
													OutputUnits: fdsn.Units{Name: s.OutputUnits},
												},
												ApproximationType: func() fdsn.ApproximationType {
													switch p.ApproximationType {
													case "MACLAURIN":
														return fdsn.ApproximationTypeMaclaurin
													default:
														return fdsn.ApproximationTypeUnknown
													}
												}(),
												FrequencyLowerBound:     fdsn.Frequency{fdsn.Float{Value: p.FrequencyLowerBound}},
												FrequencyUpperBound:     fdsn.Frequency{fdsn.Float{Value: p.FrequencyUpperBound}},
												ApproximationLowerBound: strconv.FormatFloat(p.ApproximationLowerBound, 'g', -1, 64),
												ApproximationUpperBound: strconv.FormatFloat(p.ApproximationUpperBound, 'g', -1, 64),
												MaximumError:            p.MaximumError,
												Coefficients:            coeffs,
											}
											stages = append(stages, fdsn.ResponseStage{
												Number:     fdsn.Counter(uint32(count)),
												Polynomial: &poly,
												StageGain: fdsn.Gain{
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
											if _, ok := pazMap[s.Lookup]; !ok {
												panic("unknown paz filter: " + s.Lookup)
											}
											var poles []fdsn.PoleZero
											for j, p := range pazMap[s.Lookup].Poles {
												poles = append(poles, fdsn.PoleZero{
													Number:    uint32(j),
													Real:      fdsn.FloatNoUnit{Value: real(p.Value)},
													Imaginary: fdsn.FloatNoUnit{Value: imag(p.Value)},
												})
											}
											var zeros []fdsn.PoleZero
											for j, z := range pazMap[s.Lookup].Zeros {
												zeros = append(zeros, fdsn.PoleZero{
													Number:    uint32(len(pazMap[s.Lookup].Poles) + j),
													Real:      fdsn.FloatNoUnit{Value: real(z.Value)},
													Imaginary: fdsn.FloatNoUnit{Value: imag(z.Value)},
												})
											}
											paz := fdsn.PolesZeros{
												BaseFilter: fdsn.BaseFilter{
													ResourceId:  "PolesZeros#" + v,
													Name:        fmt.Sprintf("%s.%s.%s%c.%04d.%03d.stage_%d", k, l.LocationCode, r.Label, j, on.Year(), on.YearDay(), count),
													InputUnits:  fdsn.Units{Name: s.InputUnits},
													OutputUnits: fdsn.Units{Name: s.OutputUnits},
												},
												PzTransferFunctionType: func() fdsn.PzTransferFunction {
													switch pazMap[s.Lookup].Code {
													case "A":
														return fdsn.PZFunctionLaplaceRadiansPerSecond
													case "B":
														return fdsn.PZFunctionLaplaceHertz
													case "D":
														return fdsn.PZFunctionLaplaceZTransform
													default:
														return fdsn.PZFunctionUnknown
													}
												}(),
												NormalizationFactor: func() float64 {
													w := complex(0.0, -2.0*math.Pi*s.Frequency)
													h := complex(float64(1.0), float64(0.0))
													for _, p := range pazMap[s.Lookup].Poles {
														switch pazMap[s.Lookup].Code {
														case "B":
															h /= (w - 2.0*math.Pi*p.Value)
														default:
															h /= (w - p.Value)
														}
													}
													for _, z := range pazMap[s.Lookup].Zeros {
														switch pazMap[s.Lookup].Code {
														case "B":
															h *= (w - 2.0*math.Pi*z.Value)
														default:
															h *= (w - z.Value)
														}
													}
													return 1.0 / cmplx.Abs(h)
												}(),

												NormalizationFrequency: fdsn.Frequency{
													fdsn.Float{Value: s.Frequency},
												},
												Zeros: zeros,
												Poles: poles,
											}
											stages = append(stages, fdsn.ResponseStage{
												Number:     fdsn.Counter(uint32(count)),
												PolesZeros: &paz,
												StageGain: fdsn.Gain{
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
											coefs := fdsn.Coefficients{
												BaseFilter: fdsn.BaseFilter{
													ResourceId:  "Coefficients#" + v,
													Name:        fmt.Sprintf("%s.%s.%s%c.%04d.%03d.stage_%d", k, l.LocationCode, r.Label, j, on.Year(), on.YearDay(), count),
													InputUnits:  fdsn.Units{Name: s.InputUnits},
													OutputUnits: fdsn.Units{Name: s.OutputUnits},
												},
												CfTransferFunctionType: fdsn.CfFunctionDigital,
											}
											stages = append(stages, fdsn.ResponseStage{
												Number:       fdsn.Counter(uint32(count)),
												Coefficients: &coefs,
												StageGain: fdsn.Gain{
													Value: func() float64 {
														if s.Gain != 0.0 {
															return s.Gain
														}
														return 1.0
													}(),
													Frequency: s.Frequency,
												},
												Decimation: &fdsn.Decimation{
													InputSampleRate: fdsn.Frequency{fdsn.Float{Value: s.SampleRate}},
													Factor: func() int32 {
														if s.Factor != 0 {
															return s.Factor
														}
														return 1
													}(),
													Delay:      fdsn.Float{Value: s.Delay},
													Correction: fdsn.Float{Value: s.Correction},
												},
											})

										case "fir":
											if _, ok := firMap[s.Lookup]; !ok {
												panic("unknown fir filter: " + s.Lookup)
											}
											var coeffs []fdsn.NumeratorCoefficient
											for j, c := range firMap[s.Lookup].Factors {
												coeffs = append(coeffs, fdsn.NumeratorCoefficient{
													Coefficient: int32(j + 1),
													Value:       c,
												})
											}

											fir := fdsn.FIR{
												BaseFilter: fdsn.BaseFilter{
													ResourceId:  "FIR#" + v,
													Name:        fmt.Sprintf("%s.%s.%s%c.%04d.%03d.stage_%d", k, l.LocationCode, r.Label, j, on.Year(), on.YearDay(), count),
													InputUnits:  fdsn.Units{Name: s.InputUnits},
													OutputUnits: fdsn.Units{Name: s.OutputUnits},
												},
												Symmetry: func() fdsn.Symmetry {
													switch strings.ToUpper(firMap[s.Lookup].Symmetry) {
													case "NONE":
														return fdsn.SymmetryNone
													case "EVEN":
														return fdsn.SymmetryEven
													case "ODD":
														return fdsn.SymmetryOdd

													default:
														panic("unknown fir symmetry: " + firMap[s.Lookup].Symmetry)
													}
												}(),
												NumeratorCoefficients: coeffs,
											}
											stages = append(stages, fdsn.ResponseStage{
												Number: fdsn.Counter(uint32(count)),
												FIR:    &fir,
												StageGain: fdsn.Gain{
													Value: func() float64 {
														if s.Gain != 0.0 {
															return s.Gain
														}
														return 1.0
													}(),
													Frequency: s.Frequency,
												},
												Decimation: &fdsn.Decimation{
													InputSampleRate: fdsn.Frequency{fdsn.Float{Value: firMap[s.Lookup].Decimation * s.SampleRate}},
													Factor: func() int32 {
														if s.Factor != 0 {
															return s.Factor
														}
														return 1
													}(),
													Delay:      fdsn.Float{Value: s.Delay},
													Correction: fdsn.Float{Value: s.Correction},
												},
											})
										default:
											panic("unknown filter type: " + v + ":" + s.Type)
										}
										count++
									}
								}

								chans = append(chans, fdsn.Channel{
									BaseNode: fdsn.BaseNode{
										Code:      r.Label + string(j),
										StartDate: &fdsn.DateTime{on},
										EndDate:   &fdsn.DateTime{off},
										RestrictedStatus: func() fdsn.RestrictedStatus {
											switch n.Restricted {
											case true:
												return fdsn.StatusClosed
											default:
												return fdsn.StatusOpen
											}
										}(),
									},
									LocationCode: l.LocationCode,
									Latitude: fdsn.Latitude{
										LatitudeBase: fdsn.LatitudeBase{
											Float: fdsn.Float{
												Value: l.Latitude,
											},
										},
									},
									Longitude: fdsn.Longitude{
										LongitudeBase: fdsn.LongitudeBase{
											Float: fdsn.Float{
												Value: l.Longitude,
											},
										},
									},
									Elevation: fdsn.Distance{Float: fdsn.Float{Value: l.Elevation}},
									Depth:     fdsn.Distance{Float: fdsn.Float{Value: -s.Height}},
									Azimuth:   &fdsn.Azimuth{Float: fdsn.Float{Value: azimuth}},
									Dip:       &fdsn.Dip{Float: fdsn.Float{Value: c.Dip}},
									Types:     types,
									SampleRateGroup: fdsn.SampleRateGroup{
										SampleRate: fdsn.SampleRate{Float: fdsn.Float{Value: r.Rate}},
										SampleRateRatio: func() *fdsn.SampleRateRatio {
											if r.Rate > 1.0 {
												return &fdsn.SampleRateRatio{
													NumberSamples: int32(r.Rate),
													NumberSeconds: 1,
												}
											} else {
												return &fdsn.SampleRateRatio{
													NumberSamples: 1,
													NumberSeconds: int32(1.0 / r.Rate),
												}
											}
										}(),
									},
									StorageFormat: r.StorageFormat,
									ClockDrift:    &fdsn.ClockDrift{Float: fdsn.Float{Value: r.ClockDrift}},
									Sensor: &fdsn.Equipment{
										ResourceId: "Sensor#" + s.Model + ":" + s.Serial,
										Type: func() string {
											if t, ok := sensors[s.Model]; ok {
												return t.Type
											}
											return ""
										}(),
										Description: func() string {
											if t, ok := sensors[s.Model]; ok {
												return t.Description
											}
											return ""
										}(),
										Manufacturer: func() string {
											if t, ok := sensors[s.Model]; ok {
												return t.Manufacturer
											}
											return ""
										}(),
										Vendor: func() string {
											if t, ok := sensors[s.Model]; ok {
												return t.Vendor
											}
											return ""
										}(),
										Model:        s.Model,
										SerialNumber: s.Serial,
										InstallationDate: func() *fdsn.DateTime {
											return &fdsn.DateTime{s.Start}
										}(),
										RemovalDate: func() *fdsn.DateTime {
											if time.Now().After(s.End) {
												return &fdsn.DateTime{s.End}
											}
											return nil
										}(),
									},

									DataLogger: &fdsn.Equipment{
										ResourceId: "Datalogger#" + d.Model + ":" + d.Serial,
										Type: func() string {
											if t, ok := dataloggers[d.Model]; ok {
												return t.Type
											}
											return ""
										}(),
										Description: func() string {
											if t, ok := dataloggers[d.Model]; ok {
												return t.Description
											}
											return ""
										}(),
										Manufacturer: func() string {
											if t, ok := dataloggers[d.Model]; ok {
												return t.Manufacturer
											}
											return ""
										}(),
										Vendor: func() string {
											if t, ok := dataloggers[d.Model]; ok {
												return t.Vendor
											}
											return ""
										}(),
										Model:        d.Model,
										SerialNumber: d.Serial,
										InstallationDate: func() *fdsn.DateTime {
											return &fdsn.DateTime{d.Start}
										}(),
										RemovalDate: func() *fdsn.DateTime {
											if time.Now().After(d.End) {
												return &fdsn.DateTime{d.End}
											}
											return nil
										}(),
									},
									Response: &fdsn.Response{
										Stages: stages,
										InstrumentSensitivity: &fdsn.Sensitivity{
											Gain: fdsn.Gain{
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
											InputUnits: func() fdsn.Units {
												var units fdsn.Units
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
											OutputUnits: func() fdsn.Units {
												var units fdsn.Units
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

		p := Place{Latitude: v.Latitude, Longitude: v.Longitude}

		stas[n.ExternalCode] = append(stas[n.ExternalCode], fdsn.Station{
			BaseNode: fdsn.BaseNode{
				Code:        v.Code,
				Description: n.Description,
				RestrictedStatus: func() fdsn.RestrictedStatus {
					switch n.Restricted {
					case true:
						return fdsn.StatusClosed
					default:
						return fdsn.StatusOpen
					}
				}(),
				StartDate: &(fdsn.DateTime{v.Start}),
				EndDate:   &(fdsn.DateTime{v.End}),
			},
			Latitude: fdsn.Latitude{LatitudeBase: fdsn.LatitudeBase{
				Float: fdsn.Float{
					Value: v.Latitude,
				}}, Datum: v.Datum},
			Longitude: fdsn.Longitude{LongitudeBase: fdsn.LongitudeBase{
				Float: fdsn.Float{
					Value: v.Longitude,
				}}, Datum: v.Datum},
			Elevation: fdsn.Distance{
				Float: fdsn.Float{Value: v.Elevation},
			},
			Site: fdsn.Site{
				Name: func() string {
					if v.Name != "" {
						return v.Name
					} else {
						return v.Code
					}
				}(),
				Description: func() string {
					if l := Locations.Closest(p); l != nil {
						if d := l.Distance(p); d < 5.0 {
							return fmt.Sprintf("within 5 km of %s", l.Name)
						} else {
							return fmt.Sprintf("%.0f km %s of %s", d, l.Compass(p), l.Name)
						}
					}
					return ""
				}(),
			},
			CreationDate: fdsn.DateTime{v.Start},
			TerminationDate: func() *fdsn.DateTime {
				if time.Now().Before(v.End) {
					return nil
				}
				return &fdsn.DateTime{v.End}
			}(),
			Channels: chans,
		})
	}

	for k, _ := range stas {
		n, ok := networkMap[k]
		if !ok {
			panic("unknown network " + k)
		}
		var on, off *fdsn.DateTime

		if !networkMatch.MatchString(n.ExternalCode) {
			log.Printf("skipping network: %s [outside regexp match]\n", n.ExternalCode)
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
		networks = append(networks, fdsn.Network{
			BaseNode: fdsn.BaseNode{
				Code:        n.ExternalCode,
				Description: n.Description,
				RestrictedStatus: func() fdsn.RestrictedStatus {
					switch n.Restricted {
					case true:
						return fdsn.StatusClosed
					default:
						return fdsn.StatusOpen
					}
				}(),
				StartDate: on,
				EndDate:   off,
			},
			SelectedNumberStations: uint32(len(stas[k])),
			Stations:               stas[k],
		})
	}

	root := fdsn.NewFDSNStationXML(source, sender, module, "", networks)
	if ok := root.IsValid(); ok != nil {
		panic("invalid stationxml file")
	}

	// marshal into xml
	x, err := root.Marshal()
	if err != nil {
		panic(err)
	}

	// output as needed ...
	if output != "-" {
		if err := os.MkdirAll(filepath.Dir(output), 0755); err != nil {
			panic(err)
		}
		if err := ioutil.WriteFile(output, x, 0644); err != nil {
			panic(err)
		}
	} else {
		fmt.Fprintln(os.Stdout, string(x))
	}
}
