package main

import (
	"flag"
	"fmt"
	"regexp"
	"sort"
	"time"
	//	"sort"
	//"io/ioutil"
	"log"
	"os"
	"strings"
	//"path/filepath"
	//"text/template"
	//"time"

	"github.com/GeoNet/delta/internal/build/v1.2"
	"github.com/GeoNet/delta/internal/stationxml/v1.2"

	"github.com/GeoNet/delta"
	"github.com/GeoNet/delta/meta"
	"github.com/GeoNet/delta/resp"
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

	/*
		var stationxml string
		flag.StringVar(&stationxml, "stationxml", "", "set stationxml version")
	*/

	/*
		var sensorRegexp string
		flag.StringVar(&sensorRegexp, "sensors", ".*", "regexp selection of sensors")

		var dataloggerRegexp string
		flag.StringVar(&dataloggerRegexp, "dataloggers", ".*", "regexp selection of dataloggers")

		var installed bool
		flag.BoolVar(&installed, "installed", false, "set station times based on installation dates")

		var operational bool
		flag.BoolVar(&operational, "operational", false, "only output operational channels")

		var active bool
		flag.BoolVar(&active, "active", false, "only output stations with active channels")

		var offset time.Duration
		flag.DurationVar(&offset, "operational-offset", 0, "provide a recently closed window for operational only requests")

		var single bool
		flag.BoolVar(&single, "single", false, "produce single station xml files")

		var directory string
		flag.StringVar(&directory, "directory", "xml", "where to store station xml files")

		var purge bool
		flag.BoolVar(&purge, "purge", false, "remove unknown single xml files")

		var plate string
		flag.StringVar(&plate, "template", "station_{{.Code}}_{{with $s := index .Stations 0}}{{$s.Code}}{{end}}.xml", "how to name the single station xml files")
	*/

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

	resps := make(map[string][]byte)
	for _, c := range set.Components() {
		if c.Response == "" {
			continue
		}
		data, err := resp.LookupBase(lookup, c.Response)
		if err != nil {
			log.Fatalf("unable to decode response %s: %v", c.Response, err)
		}
		if data == nil {
			continue
		}
		resps[c.Response] = data
	}
	for _, c := range set.Channels() {
		if c.Response == "" {
			continue
		}
		data, err := resp.LookupBase(lookup, c.Response)
		if err != nil {
			log.Fatalf("unable to decode response %s: %v", c.Response, err)
		}
		if data == nil {
			continue
		}
		resps[c.Response] = data
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

							at := c.InstalledSensor.Start
							if c.DeployedDatalogger.Start.After(at) {
								at = c.DeployedDatalogger.Start
							}
							if c.Connection != nil && c.Connection.Start.After(at) {
								at = c.Connection.Start
							}

							prefix := fmt.Sprintf("%s.%s.%s.%s.", site.Station, site.Location, c.Code(), at.Format("2006.002"))

							resp := build.NewResponse(prefix, c.InstalledSensor.Serial, LegacyFrequency(c.Code()))
							if site.Station == "KAVZ" {
								resp = build.NewResponse(prefix, c.InstalledSensor.Serial, 1.0)
							}

							var gain, bias float64
							if v.Gain != nil {
								gain, bias = v.Gain.Factor, v.Gain.Bias
							}

							var preamp float64
							if v.Preamp != nil {
								preamp = v.Preamp.Gain
							}

							res, err := func() (*stationxml.ResponseType, error) {
								if c.Component.SamplingRate != 0 {
									derived, ok := resps[c.Component.Response]
									if !ok || !(len(derived) > 0) {
										return nil, nil
									}
									return resp.Derived(derived)
								}
								sensor, ok := resps[c.Component.Response]
								if !ok || !(len(sensor) > 0) {
									return nil, nil
								}
								if err := resp.Sensor(gain, bias, sensor); err != nil {
									return nil, err
								}
								datalogger, ok := resps[c.Channel.Response]
								if !ok || !(len(datalogger) > 0) {
									return nil, nil
								}
								if err := resp.Datalogger(preamp, datalogger); err != nil {
									return nil, err
								}
								if err := resp.Normalise(); err != nil {
									return nil, err
								}

								//return Resp(resp, prefix, c.InstalledSensor.Serial, LegacyFrequency(c.Code()), gain, bias, preamp, resps[c.Component.Response], resps[c.Channel.Response])
								return resp.ResponseType(), nil

							}()
							if err != nil {
								log.Fatal(err)
								//return err
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

	/*
		nets := make(map[string]meta.Network)
		for _, n := range set.Networks() {
			if !network.MatchString(n.Code) {
				continue
			}
			if !external.MatchString(n.External) {
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
	*/

	/*
		networks := make(map[meta.Network][]meta.Station)
		for _, n := range set.Networks() {
			if !network.MatchString(n.Code) {
				continue
			}
			if !external.MatchString(n.External) {
				continue
			}
			ext, ok := set.Network(n.External)
			if !ok {
				continue
			}
			for _, s := range set.Stations() {
				if s.Network != n.Code {
					continue
				}
				if !station.MatchString(s.Code) {
					continue
				}
				networks[ext] = append(networks[ext], s)
			}
		}
	*/

	//resps := make(map[string][]byte)

	/*
		cols := make(map[string][]meta.Collection)
		for _, c := range set.Collections() {
			if _, ok := stns[c.InstalledSensor.Station]; !ok {
				continue
			}

			if !location.MatchString(c.InstalledSensor.Location) {
				continue
			}
			if !channel.MatchString(c.Code()) {
				continue
			}

			cols[c.InstalledSensor.Station] = append(cols[c.InstalledSensor.Station], c)
		}

		log.Println(len(cols), time.Since(start))
	*/

	return

	//sites := make(map[meta.Site][]meta.Collection)
	/*
		for _, list := range networks {
			for _, sta := range list {
				for _, x := range set.Sites() {
					if x.Station != sta.Code {
						continue
					}
					if !location.MatchString(x.Location) {
						continue
					}
					for _, c := range set.Collections() {

						cols[sta.Code] = append(cols[sta.Code], c)
						if !channel.MatchString(c.Code()) {
							continue
						}

						sites[x] = append(sites[x], c)
						if r := c.Component.Response; r != "" {
							d, err := resp.LookupBase(lookup, r)
							if err != nil {
								log.Fatalf("unable to recover response file %q: %v", r, err)
							}
							resps[r] = d
						}
						if r := c.Channel.Response; r != "" {
							d, err := resp.LookupBase(lookup, r)
							if err != nil {
								log.Fatalf("unable to recover response file %q: %v", r, err)
							}
							resps[r] = d
						}
					}
				}
			}
		}
	*/

	/*
		somethings := make(map[meta.Collection][]meta.Something)
		for _, v := range cols {
			for _, c := range v {
	*/
	/*
		if c.Stream.Station != "AWAZ" {
			continue
		}
		log.Println(c)
		for _, x := range set.Somethings(c) {
			log.Println("\t", x.Gain.IsNominal(), x)
		}
		return
	*/
	/*
						somethings[c] = append(somethings[c], set.Somethings(c)...)
					}
				}


		nets := make(map[string]meta.Network)
		for _, n := range set.Networks() {
			nets[n.Code] = n
		}

		stns := make(map[string]meta.Station)
		for _, n := range set.Stations() {
			stns[n.Code] = n
		}
	*/

	/*
		switch stationxml {
		default:
			res, err := build10(source, sender, module, places, nets, stns, somethings, cols, networks, sites, resps)
			if err != nil {
				log.Fatal(err)
			}
			switch output {
			case "", "-":
				fmt.Fprintln(os.Stdout, string(res))
			default:
				file, err := os.Create(output)
				if err != nil {
					log.Fatalf("error: unable to create file %s: %v", output, err)
				}
				if _, err := file.Write(res); err != nil {
					log.Fatalf("error: unable to write file %s: %v", output, err)
				}
				if err := file.Close(); err != nil {
					log.Fatalf("error: unable to close file %s: %v", output, err)
				}
			}
		}
	*/

	/**
	builder, err := NewBuilder(
		SetBase(base),
		SetInstalled(installed),
		SetActive(active),
		SetOperational(operational, offset),
		SetNetworks(networkRegexp),
		SetExternal(externalRegexp),
		SetStations(stationRegexp),
		SetChannels(channelRegexp),
		SetSensors(sensorRegexp),
		SetDataloggers(dataloggerRegexp),
	)
	if err != nil {
		log.Fatalf("unable to make builder: %v", err)
	}

	// build a representation of the network
	networks, err := builder.Construct()
	if err != nil {
		log.Fatalf("error: unable to build networks list: %v", err)
	}

	switch {
	case single:
		tmpl, err := template.New("single").Parse(plate)
		if err != nil {
			log.Fatalf("unable to parse single xml file template: %v", err)
		}

		if err := os.MkdirAll(directory, 0755); err != nil {
			log.Fatalf("error: unable to create directory %s: %v", directory, err)
		}

		files := make(map[string]string)

		if purge {
			if err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if info.IsDir() {
					return nil
				}
				files[filepath.Base(path)] = path
				return nil
			}); err != nil {
				log.Fatalf("unable to walk dir %s: %v", directory, err)
			}
		}

		var count, updated int
		for _, n := range networks {
			for _, s := range n.Stations {
				node := stationxml.Network{
					BaseNode:               n.BaseNode,
					TotalNumberStations:    1,
					SelectedNumberStations: 1,
					Stations:               []stationxml.Station{s},
				}

				root := stationxml.NewFDSNStationXML(source, sender, module, "", []stationxml.Network{node})
				if ok := root.IsValid(); ok != nil {
					log.Fatalf("error: invalid stationxml file")
				}

				// marshal into xml
				res, err := root.MarshalIndent()
				if err != nil {
					log.Fatalf("error: unable to marshal stationxml: %v", err)
				}

				var name bytes.Buffer
				if err := tmpl.Execute(&name, &node); err != nil {
					log.Fatalf("unable to encode single xml filename: %s", err)
				}

				path := filepath.Join(directory, name.String())
				if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
					log.Fatalf("error: unable to create directory %s: %v", filepath.Dir(path), err)
				}

				delete(files, name.String())

				// has anything changed?
				if !compare(path, res) {
					if err := ioutil.WriteFile(path, res, 0600); err != nil {
						log.Fatalf("error: unable to write file %s: %v", path, err)
					}
					updated++
				}

				count++
			}
		}

		for k, v := range files {
			log.Printf("removing extra file: %s", k)
			if err := os.Remove(v); err != nil {
				log.Fatalf("unable to remove file %s: %v", k, err)
			}
		}

		log.Printf("built %d files, updated %d, removed %d", count, updated, len(files))

	default:
		// render full station xml
		root := stationxml.NewFDSNStationXML(source, sender, module, "", networks)
		if ok := root.IsValid(); ok != nil {
			log.Fatalf("error: invalid stationxml file")
		}

		// marshal into xml
		res, err := root.MarshalIndent()
		if err != nil {
			log.Fatalf("error: unable to marshal stationxml: %v", err)
		}

		// output as needed ...
		switch output {
		case "", "-":
			fmt.Fprintln(os.Stdout, string(res))
		default:
			if err := os.MkdirAll(filepath.Dir(output), 0755); err != nil {
				log.Fatalf("error: unable to create directory %s: %v", filepath.Dir(output), err)
			}
			if err := ioutil.WriteFile(output, res, 0600); err != nil {
				log.Fatalf("error: unable to write file %s: %v", output, err)
			}
		}
	}
	**/
}
