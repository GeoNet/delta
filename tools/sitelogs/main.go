package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/GeoNet/delta/meta"
)

//go:generate bash -c "rm -f config_auto.go; go run generate/*.go | gofmt > config_auto.go"

const DateTimeFormat = "2006-01-02T15:04Z"
const DateFormat = "2006-01-02"

func main() {

	var verbose bool
	flag.BoolVar(&verbose, "verbose", false, "make noise")

	var output string
	flag.StringVar(&output, "output", "output", "output directory")

	var logs string
	flag.StringVar(&logs, "logs", "logs", "logs output directory")

	var network string
	flag.StringVar(&network, "network", "../../network", "base network directory")

	var install string
	flag.StringVar(&install, "install", "../../install", "base install directory")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Build GNSS SiteLog XML files from delta meta information\n")
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

	var tplFuncMap template.FuncMap = template.FuncMap{
		"empty": func(d, s string) string {
			if s != "" {
				return s
			}
			return d
		},
		"tolower": func(s string) string {
			switch t := strings.ToLower(s); t {
			case "wyatt/agnew drilled-braced":
				return "Deep Wyatt/Agnew drilled-braced"
			default:
				return t
			}
		},
		"lines": func(p, s string) string {
			switch s {
			case "":
				return s
			default:
				return strings.Join(strings.Split(s, "\n"), "\n"+p)
			}
		},
		"plus": func(n int) string {
			return fmt.Sprintf("%-2s", strconv.Itoa(n+1))
		},
		"lat": func(s string) string {
			if f, err := strconv.ParseFloat(s, 64); err == nil {
				m := math.Abs(f-float64(int(f))) * 60.0
				return fmt.Sprintf("%+3d%02d%05.2f", int(f), int(m), (m-float64(int(m)))*60.0)
			}
			return ""
		},
		"lon": func(s string) string {
			if f, err := strconv.ParseFloat(s, 64); err == nil {
				m := math.Abs(f-float64(int(f))) * 60.0
				return fmt.Sprintf("%+3d%02d%05.2f", int(f), int(m), (m-float64(int(m)))*60.0)
			}
			return ""
		},
	}

	tmpl, err := template.New("").Funcs(tplFuncMap).Parse(sitelogTemplate)
	if err != nil {
		log.Fatalf("error: unable to compile template: %v", err)
	}

	var firmwareHistoryList meta.FirmwareHistoryList
	if err := meta.LoadList(filepath.Join(install, "firmware.csv"), &firmwareHistoryList); err != nil {
		log.Fatalf("error: unable to load firmware history: %v", err)
	}

	firmwareHistory := make(map[string]map[string][]meta.FirmwareHistory)
	for _, i := range firmwareHistoryList {
		if _, ok := firmwareHistory[i.Model]; !ok {
			firmwareHistory[i.Model] = make(map[string][]meta.FirmwareHistory)
		}
		firmwareHistory[i.Model][i.Serial] = append(firmwareHistory[i.Model][i.Serial], i)
	}

	for j, _ := range firmwareHistory {
		for k, _ := range firmwareHistory[j] {
			sort.Sort(meta.FirmwareHistoryList(firmwareHistory[j][k]))
		}
	}

	var installedAntennaList meta.InstalledAntennaList
	if err := meta.LoadList(filepath.Join(install, "antennas.csv"), &installedAntennaList); err != nil {
		log.Fatalf("error: unable to load antenna installs: %v", err)
	}

	installedAntenna := make(map[string][]meta.InstalledAntenna)
	for _, i := range installedAntennaList {
		installedAntenna[i.Mark] = append(installedAntenna[i.Mark], i)
	}
	for i, _ := range installedAntenna {
		sort.Sort(meta.InstalledAntennaList(installedAntenna[i]))
	}

	var deployedReceiverList meta.DeployedReceiverList
	if err := meta.LoadList(filepath.Join(install, "receivers.csv"), &deployedReceiverList); err != nil {
		log.Fatalf("error: unable to load receiver installs: %v", err)
	}

	deployedReceivers := make(map[string][]meta.DeployedReceiver)
	for _, i := range deployedReceiverList {
		deployedReceivers[i.Mark] = append(deployedReceivers[i.Mark], i)
	}
	for i, _ := range deployedReceivers {
		sort.Sort(meta.DeployedReceiverList(deployedReceivers[i]))
	}

	var installedRadomeList meta.InstalledRadomeList
	if err := meta.LoadList(filepath.Join(install, "radomes.csv"), &installedRadomeList); err != nil {
		log.Fatalf("error: unable to load radome installs: %v", err)
	}

	installedRadomes := make(map[string][]meta.InstalledRadome)
	for _, i := range installedRadomeList {
		installedRadomes[i.Mark] = append(installedRadomes[i.Mark], i)
	}
	for i, _ := range installedRadomes {
		sort.Sort(meta.InstalledRadomeList(installedRadomes[i]))
	}

	var installedMetSensorList meta.InstalledMetSensorList
	if err := meta.LoadList(filepath.Join(install, "metsensors.csv"), &installedMetSensorList); err != nil {
		log.Fatalf("error: unable to load metsensors list: %v", err)
	}

	installedMetSensors := make(map[string][]meta.InstalledMetSensor)
	for _, i := range installedMetSensorList {
		installedMetSensors[i.Mark] = append(installedMetSensors[i.Mark], i)
	}
	for i, _ := range installedMetSensors {
		sort.Sort(meta.InstalledMetSensorList(installedMetSensors[i]))
	}

	var markList meta.MarkList
	if err := meta.LoadList(filepath.Join(network, "marks.csv"), &markList); err != nil {
		log.Fatalf("error: unable to load mark list: %v", err)
	}

	var sessionList meta.SessionList
	if err := meta.LoadList(filepath.Join(install, "sessions.csv"), &sessionList); err != nil {
		log.Fatalf("error: unable to load session list: %v", err)
	}

	sessions := make(map[string][]meta.Session)
	for _, s := range sessionList {
		sessions[s.Mark] = append(sessions[s.Mark], s)
	}

	var monumentList meta.MonumentList
	if err := meta.LoadList(filepath.Join(network, "monuments.csv"), &monumentList); err != nil {
		log.Fatalf("error: unable to load monument list: %v", err)
	}
	monuments := make(map[string]meta.Monument)
	for _, m := range monumentList {
		monuments[m.Mark] = m
	}

	for _, m := range markList {

		if _, ok := monuments[m.Code]; !ok {
			continue
		}

		if _, ok := sessions[m.Code]; !ok {
			continue
		}
		if _, ok := installedAntenna[m.Code]; !ok {
			continue
		}
		if _, ok := deployedReceivers[m.Code]; !ok {
			continue
		}

		var receivers []GnssReceiver
		var antennas []GnssAntenna
		var metsensors []GnssMetSensor

		for _, m := range installedMetSensors[m.Reference.Code] {
			var session *meta.Session
			for i, s := range sessions[m.Mark] {
				if m.Start.After(s.End) || m.End.Before(s.Start) {
					continue
				}
				session = &sessions[m.Mark][i]
				break
			}
			if session == nil {
				continue
			}
			metsensors = append(metsensors, GnssMetSensor{
				Manufacturer:         m.Make,
				MetSensorModel:       m.Model,
				SerialNumber:         m.Serial,
				DataSamplingInterval: "360 sec",
				EffectiveDates:       "2000-02-05/CCYY-MM-DD",
				Notes:                "",
			})
		}

		for _, a := range installedAntenna[m.Code] {
			var session *meta.Session
			for i, s := range sessions[m.Code] {
				if a.Start.After(s.End) || a.End.Before(s.Start) {
					continue
				}
				session = &sessions[m.Code][i]
				break
			}
			if session == nil {
				continue
			}

			radome := "NONE"
			serial := ""
			if _, ok := installedRadomes[m.Code]; ok {
				for _, v := range installedRadomes[m.Code] {
					if v.Start.After(a.End) || v.End.Before(a.Start) {
						continue
					}
					radome = v.Model
					serial = v.Serial
				}
			}

			antennas = append(antennas, GnssAntenna{
				AntennaType:            a.Model,
				SerialNumber:           a.Serial,
				AntennaReferencePoint:  "BAM",
				MarkerArpUpEcc:         strconv.FormatFloat(a.Vertical, 'f', 4, 64),
				MarkerArpNorthEcc:      strconv.FormatFloat(a.North, 'f', 4, 64),
				MarkerArpEastEcc:       strconv.FormatFloat(a.East, 'f', 4, 64),
				AlignmentFromTrueNorth: "0",
				AntennaRadomeType:      radome,
				RadomeSerialNumber:     serial,
				AntennaCableType:       "",
				AntennaCableLength:     "",
				DateInstalled:          a.Start.Format(DateTimeFormat),
				DateRemoved: func() string {
					if time.Now().After(a.End) {
						return a.End.Format(DateTimeFormat)
					} else {
						return ""
					}
				}(),
				Notes: "",
			})
		}

		for _, r := range deployedReceivers[m.Code] {
			if _, ok := firmwareHistory[r.Model]; ok {
				if _, ok := firmwareHistory[r.Model][r.Serial]; ok {
					for i, _ := range firmwareHistory[r.Model][r.Serial] {

						v := firmwareHistory[r.Model][r.Serial][len(firmwareHistory[r.Model][r.Serial])-i-1]
						if v.End.Before(r.Start) || v.Start.After(r.End) {
							continue
						}

						var session *meta.Session
						for i, s := range sessions[m.Code] {
							if r.Start.After(s.End) || r.End.Before(s.Start) {
								continue
							}
							if v.Start.After(s.End) || v.End.Before(s.Start) {
								continue
							}
							session = &sessions[m.Code][i]
							break
						}
						if session == nil {
							continue
						}

						start := r.Start
						/*
							if start.Before(s.Start) {
								start = s.Start
							}
						*/
						if start.Before(v.Start) {
							start = v.Start
						}

						end := r.End
						/*
							if end.After(s.End) {
								end = s.End
							}
						*/
						if end.After(v.End) {
							end = v.End
						}

						receivers = append(receivers, GnssReceiver{
							ReceiverType:           r.Model,
							SatelliteSystem:        session.SatelliteSystem,
							SerialNumber:           r.Serial,
							FirmwareVersion:        v.Version,
							ElevationCutoffSetting: strconv.FormatFloat(session.ElevationMask, 'g', -1, 64),
							DateInstalled:          start.Format(DateTimeFormat),
							/*
								DateInstalled: func() string {
									if v.Start.Before(r.Start) {
										return r.Start.Format(DateTimeFormat)
									} else {
										return v.Start.Format(DateTimeFormat)
									}
								}(),
							*/
							DateRemoved: func() string {
								/*
									if v.End.After(r.End) {
										if time.Now().After(r.End) {
											return r.End.Format(DateTimeFormat)
										} else {
											return ""
										}
									} else {
										if time.Now().After(v.End) {
											return v.End.Format(DateTimeFormat)
										} else {
											return ""
										}
									}
								*/
								if time.Now().After(end) {
									return end.Format(DateTimeFormat)
								} else {
									return ""
								}
							}(),
							TemperatureStabilization: "",
							Notes: "",
						})
					}
				}
			}
		}

		sort.Sort(GnssReceivers(receivers))
		sort.Sort(GnssAntennas(antennas))

		monument := monuments[m.Code]

		X, Y, Z := WGS842ITRF(m.Latitude, m.Longitude, m.Elevation)

		x := SiteLog{

			EquipNameSpace:   equipNameSpace,
			ContactNameSpace: contactNameSpace,
			MiNameSpace:      miNameSpace,
			LiNameSpace:      liNameSpace,
			XmlNameSpace:     xmlNameSpace,
			XsiNameSpace:     xsiNameSpace,
			SchemaLocation:   schemaLocation,

			FormInformation: FormInformation{
				PreparedBy:   preparedBy,
				DatePrepared: time.Now().Format(DateFormat),
				ReportType:   "DYNAMIC",
			},

			SiteIdentification: SiteIdentification{
				SiteName:            m.Name,
				FourCharacterID:     m.Code,
				MonumentInscription: "",
				IersDOMESNumber:     monument.DomesNumber,
				CdpNumber:           "",
				MonumentDescription: monument.Type,
				HeightOfTheMonument: strconv.FormatFloat(-monument.GroundRelationship, 'g', -1, 64),
				MonumentFoundation:  monument.FoundationType,
				FoundationDepth:     strconv.FormatFloat(monument.FoundationDepth, 'f', 1, 64),
				MarkerDescription: func() string {
					switch monument.MarkType {
					case "Forced Centering":
						return "Forced Centering"
					default:
						return "unknown"
					}
				}(),
				DateInstalled:          m.Start.Format(DateTimeFormat),
				GeologicCharacteristic: "",
				BedrockType:            "",
				BedrockCondition:       "",
				FractureSpacing:        "",
				FaultZonesNearby:       "",
				DistanceActivity:       "",
				Notes:                  "",
			},
			SiteLocation: SiteLocation{
				/*
					City:          m.Place,
					State:         m.Region,
				*/
				Country: func(lat, lon float64) string {
					X, Y, _ := WGS842ITRF(lat, lon, 0.0)
					dist := float64(-1.0)
					country := "Unknown"
					for _, v := range countryList {
						x, y, _ := WGS842ITRF(v.lat, v.lon, 0.0)
						r := math.Sqrt((x-X)*(x-X) + (y-Y)*(y-Y))
						if dist < 0.0 || r < dist {
							country = v.name
							dist = r
						}
					}

					return country
				}(m.Latitude, m.Longitude),

				TectonicPlate: TectonicPlate(m.Latitude, m.Longitude),
				ApproximatePositionITRF: ApproximatePositionITRF{
					XCoordinateInMeters: strconv.FormatFloat(X, 'f', 1, 64),
					YCoordinateInMeters: strconv.FormatFloat(Y, 'f', 1, 64),
					ZCoordinateInMeters: strconv.FormatFloat(Z, 'f', 1, 64),
					LatitudeNorth:       strconv.FormatFloat(m.Latitude, 'g', -1, 64),
					LongitudeEast:       strconv.FormatFloat(m.Longitude, 'g', -1, 64),
					ElevationMEllips:    strconv.FormatFloat(m.Elevation, 'f', 1, 64),
				},
				Notes: "",
			},
			GnssReceivers:  receivers,
			GnssAntennas:   antennas,
			GnssMetSensors: metsensors,
			ContactAgency:  contactAgency,
			ResponsibleAgency: func() Agency {
				switch m.Network {
				case "LI":
					return responsibleAgency
				default:
					return Agency{
						MailingAddress: "\n",
					}
				}
			}(),
			MoreInformation: MoreInformation{
				PrimaryDataCenter:     primaryDatacentre,
				SecondaryDataCenter:   "",
				UrlForMoreInformation: urlForMoreInformation,
				HardCopyOnFile:        "",
				SiteMap:               "",
				SiteDiagram:           "",
				HorizonMask:           "",
				MonumentDescription:   "",
				SitePictures:          "",
				Notes:                 extraNotes + " " + m.Code,
				AntennaGraphicsWithDimensions: func() string {
					var graphs []string
					models := make(map[string]interface{})
					for _, a := range antennas {
						if _, ok := models[a.AntennaType]; ok {
							continue
						}
						if g, ok := antennaGraphs[a.AntennaType]; ok {
							b, err := hex.DecodeString(g)
							if err != nil {
								log.Printf("error: unable to decode antenna graph for: \"%s\"", a.AntennaType)
								continue
							}
							graphs = append(graphs, strings.Join([]string{a.AntennaType, "", string(b)}, "\n"))
						} else {
							log.Printf("warning: missing antenna graph for: \"%s\"", a.AntennaType)
						}
						models[a.AntennaType] = true
					}
					return strings.Join(graphs, "\n") // + "\n"
				}(),
				InsertTextGraphicFromAntenna: "",
			},
		}

		s, err := x.Marshal()
		if err != nil {
			log.Fatalf("error: unable to marshal xml: %v", err)
		}

		xmlfile := filepath.Join(output, strings.ToLower(m.Code)+".xml")
		if err := os.MkdirAll(filepath.Dir(xmlfile), 0755); err != nil {
			log.Fatalf("error: unable to create dir: %v", err)
		}
		if err := ioutil.WriteFile(xmlfile, s, 0644); err != nil {
			log.Fatalf("error: unable to write file: %v", err)
		}

		logfile := filepath.Join(logs, strings.ToLower(m.Code)+".log")
		if err := os.MkdirAll(filepath.Dir(logfile), 0755); err != nil {
			log.Fatalf("error: unable to create logs dir: %v", err)
		}
		f, err := os.Create(logfile)
		if err != nil {
			log.Fatalf("error: unable to create log file: %v", err)
		}
		defer f.Close()

		if err := tmpl.Execute(f, x); err != nil {
			log.Fatalf("error: unable to write log file: %v", err)
		}
	}
}
