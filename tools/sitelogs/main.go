package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/GeoNet/delta/meta"
)

const DateTimeFormat = "2006-01-02T15:04Z"
const DateFormat = "2006-01-02"

var primaryDatacentre = "ftp.geonet.org.nz"
var urlForMoreInformation = "www.geonet.org.nz"
var extraNotes = "additional information and pictures could be found at http://magma.geonet.org.nz/delta/app then search for CGPS mark"

var contactAgency = Agency{
	Agency:                "GNS Science",
	PreferredAbbreviation: "GNS",
	MailingAddress:        "1 Fairway Drive, Avalon 5010, PO Box 30-368, Lower Hutt, New Zealand",
	PrimaryContact: Contact{
		Name:               "GeoNet reception",
		TelephonePrimary:   "+64 4 570 1444",
		TelephoneSecondary: "",
		Fax:                "+64 4 570 4676",
		Email:              "info@geonet.org.nz",
	},
	SecondaryContact: Contact{
		Name:               "Elisabetta D'Anastasio",
		TelephonePrimary:   "+64 4 570 4744",
		TelephoneSecondary: "",
		Fax:                "",
		Email:              "e.danastasio@gns.cri.nz",
	},
	Notes: "",
}

var responsibleAgency = Agency{
	Agency:                "Land Information New Zealand",
	PreferredAbbreviation: "LINZ",
	MailingAddress:        "155 The Terrace, PO Box 5501, Wellington 6145 New Zealand",
	PrimaryContact: Contact{
		Name:               "LINZ Reception",
		TelephonePrimary:   "+64 4 460 0110",
		TelephoneSecondary: "",
		Fax:                "+64 4 472 2244",
		Email:              "positionz@linz.govt.nz",
	},
	SecondaryContact: Contact{
		Name:               "Paula Gentle",
		TelephonePrimary:   "+64 4 460 2757",
		TelephoneSecondary: "",
		Fax:                "",
		Email:              "pgentle@linz.govt.nz",
	},
	Notes: "CGPS site is part of the LINZ PositioNZ Network http://www.linz.govt.nz/positionz",
}

func country(lat, lon float64) string {
	var countries = []struct {
		name     string
		lat, lon float64
	}{
		{"New Zealand", -40.0, 174.0},
		{"Tonga", -21.2, -175.2},
		{"Samoa", -13.8, -172.1},
		{"Niue", -19.0, -169.9},
	}
	X, Y, _ := WGS842ITRF(lat, lon, 0.0)

	dist := float64(-1.0)
	country := "Unknown"
	for _, v := range countries {
		x, y, _ := WGS842ITRF(v.lat, v.lon, 0.0)
		r := math.Sqrt((x-X)*(x-X) + (y-Y)*(y-Y))
		if dist < 0.0 || r < dist {
			country = v.name
			dist = r
		}
	}

	return country
}

func main() {

	var verbose bool
	flag.BoolVar(&verbose, "verbose", false, "make noise")

	var output string
	flag.StringVar(&output, "output", "output", "output directory")

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

	var firmwareHistoryList meta.FirmwareHistoryList
	if err := meta.LoadList(filepath.Join(install, "firmware.csv"), &firmwareHistoryList); err != nil {
		fmt.Fprintf(os.Stderr, "error: unable to load firmware history: %v\n", err)
		os.Exit(-1)
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
		fmt.Fprintf(os.Stderr, "error: unable to load antenna installs: %v\n", err)
		os.Exit(-1)
	}

	installedAntenna := make(map[string][]meta.InstalledAntenna)
	for _, i := range installedAntennaList {
		installedAntenna[i.MarkCode] = append(installedAntenna[i.MarkCode], i)
	}
	for i, _ := range installedAntenna {
		sort.Sort(meta.InstalledAntennaList(installedAntenna[i]))
	}

	var deployedReceiverList meta.DeployedReceiverList
	if err := meta.LoadList(filepath.Join(install, "receivers.csv"), &deployedReceiverList); err != nil {
		fmt.Fprintf(os.Stderr, "error: unable to load receiver installs: %v\n", err)
		os.Exit(-1)
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
		fmt.Fprintf(os.Stderr, "error: unable to load radome installs: %v\n", err)
		os.Exit(-1)
	}

	installedRadomes := make(map[string][]meta.InstalledRadome)
	for _, i := range installedRadomeList {
		installedRadomes[i.MarkCode] = append(installedRadomes[i.MarkCode], i)
	}
	for i, _ := range installedRadomes {
		sort.Sort(meta.InstalledRadomeList(installedRadomes[i]))
	}

	var markList meta.MarkList
	if err := meta.LoadList(filepath.Join(network, "marks.csv"), &markList); err != nil {
		fmt.Fprintf(os.Stderr, "error: unable to load mark list: %v\n", err)
		os.Exit(-1)
	}

	var sessionList meta.SessionList
	if err := meta.LoadList(filepath.Join(install, "sessions.csv"), &sessionList); err != nil {
		fmt.Fprintf(os.Stderr, "error: unable to load session list: %v\n", err)
		os.Exit(-1)
	}

	sessions := make(map[string][]meta.Session)
	for _, s := range sessionList {
		sessions[s.MarkCode] = append(sessions[s.MarkCode], s)
	}

	var monumentList meta.MonumentList
	if err := meta.LoadList(filepath.Join(network, "monuments.csv"), &monumentList); err != nil {
		fmt.Fprintf(os.Stderr, "error: unable to load monument list: %v\n", err)
		os.Exit(-1)
	}
	monuments := make(map[string]meta.Monument)
	for _, m := range monumentList {
		monuments[m.MarkCode] = m
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

			/*
				start := a.Start
				if start.Before(s.Start) {
					start = s.Start
				}

				end := a.End
				if end.After(s.End) {
					end = s.End
				}
			*/

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
				MarkerArpUpEcc:         strconv.FormatFloat(a.Height, 'f', 4, 64),
				MarkerArpNorthEcc:      strconv.FormatFloat(a.North, 'f', 4, 64),
				MarkerArpEastEcc:       strconv.FormatFloat(a.East, 'f', 4, 64),
				AlignmentFromTrueNorth: "0",
				AntennaRadomeType:      radome,
				RadomeSerialNumber:     serial,
				AntennaCableType:       "",
				AntennaCableLength:     "",
				DateInstalled:          a.Start.Format(DateTimeFormat),
				//DateInstalled: start.Format(DateTimeFormat),
				DateRemoved: func() string {
					if time.Now().After(a.End) {
						return a.End.Format(DateTimeFormat)
						//if time.Now().After(end) {
						// return end.Format(DateTimeFormat)
					} else {
						return ""
					}
				}(),
				Notes: "",
			})
		}

		for _, r := range deployedReceivers[m.Code] {
			/*
				for _, s := range sessions[m.Code] {
					if r.Start.After(s.End) || r.End.Before(s.Start) {
						continue
					}
			*/
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
				//PreparedBy:   "",
				DatePrepared: time.Now().Format(DateFormat),
				ReportType:   "DYNAMIC",
			},

			SiteIdentification: SiteIdentification{
				SiteName:            m.Name,
				FourCharacterID:     m.Code,
				MonumentInscription: "",
				IersDOMESNumber:     monument.DomesNumber,
				CdpNumber:           "",
				MonumentDescription: func() string {
					switch monument.MonumentType {
					case "Pillar":
						return "pillar"
					default:
						return "unknown"
					}
				}(),
				HeightOfTheMonument: strconv.FormatFloat(-monument.GroundRelationship, 'g', -1, 64),
				MonumentFoundation: func() string {
					switch monument.MonumentType {
					case "Pillar":
						return "reinforced concrete"
					default:
						return "unknown"
					}
				}(),
				FoundationDepth: func() string {
					switch monument.MonumentType {
					case "Pillar":
						return "2.0"
					default:
						return "unknown"
					}
				}(),
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
				Country:       country(m.Latitude, m.Longitude),
				TectonicPlate: TectonicPlate(m.Latitude, m.Longitude),
				ApproximatePositionITRF: ApproximatePositionITRF{
					XCoordinateInMeters: strconv.FormatFloat(X, 'f', 5, 64),
					YCoordinateInMeters: strconv.FormatFloat(Y, 'f', 5, 64),
					ZCoordinateInMeters: strconv.FormatFloat(Z, 'f', 5, 64),
					LatitudeNorth:       strconv.FormatFloat(m.Latitude, 'g', -1, 64),
					LongitudeEast:       strconv.FormatFloat(m.Longitude, 'g', -1, 64),
					ElevationMEllips:    strconv.FormatFloat(m.Elevation, 'g', -1, 64),
				},
				Notes: "",
			},
			GnssReceivers:     receivers,
			GnssAntennas:      antennas,
			ContactAgency:     contactAgency,
			ResponsibleAgency: responsibleAgency,
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
				Notes:                 extraNotes,
				AntennaGraphicsWithDimensions: func() string {
					var graphs []string
					models := make(map[string]interface{})
					for _, a := range antennas {
						if _, ok := models[a.AntennaType]; ok {
							continue
						}
						if g, ok := antennaGraphs[a.AntennaType]; ok {
							graphs = append(graphs, g)
						}
						models[a.AntennaType] = true
					}
					return strings.Join(graphs, "\n") + "\n"
				}(),
				InsertTextGraphicFromAntenna: "",
			},
		}

		s, err := x.Marshal()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: unable to marshal xml: %v\n", err)
			os.Exit(-1)
		}

		xmlfile := filepath.Join(output, strings.ToLower(m.Code)+".xml")
		if err := os.MkdirAll(filepath.Dir(xmlfile), 0755); err != nil {
			fmt.Fprintf(os.Stderr, "error: unable to create dir: %v\n", err)
			os.Exit(-1)
		}
		if err := ioutil.WriteFile(xmlfile, s, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "error: unable to write file: %v\n", err)
			os.Exit(-1)
		}
	}
}
