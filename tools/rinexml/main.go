package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/GeoNet/delta/meta"
)

const (
	DateFormat     = "2006-01-02"
	DateTimeFormat = "2006-01-02T15:04:05"
)

var HeaderComments map[string]string = map[string]string{
	"linz": `
			Data supplied by the GeoNet project.
			GeoNet is core funded by EQC, with support from LINZ, and is operated by GNS on behalf of EQC and all New Zealanders.
			Contact: www.geonet.org.nz  email: info@geonet.org.nz.
		`,
	"geonet": `
			Data supplied by the GeoNet project.
			GeoNet is core funded by EQC and is operated by GNS on behalf of EQC and all New Zealanders.
			Contact: www.geonet.org.nz  email: info@geonet.org.nz.
		`,
	"gsi": `
			Data supplied by the GeoNet project and GSI (Tsukuba, Japan).
			GeoNet is core funded by EQC and is operated by GNS on behalf of EQC and all New Zealanders.
			Contact: www.geonet.org.nz  email: info@geonet.org.nz.
		`,
	"sagenz": `
			Data supplied by the GeoNet project as part of the SAGENZ project.
			GeoNet is core funded by EQC and is operated by GNS on behalf of EQC and all New Zealanders.
			Contact: www.geonet.org.nz  email: info@geonet.org.nz.
		`,
	"swpacific": `
			Data supplied by the GeoNet project as part of a joint project involving GNS, Ohio State Univ.,
			Pacific GPS Facility at Univ. Hawaii, and the governments of a number of SW Pacific states.
			GeoNet is core funded by EQC and is operated by GNS on behalf of EQC and all New Zealanders.
			Contact: www.geonet.org.nz  email: info@geonet.org.nz.
		`,
}

var DownloadNameFormats map[string]DownloadNameFormatXML = map[string]DownloadNameFormatXML{
	"x4": DownloadNameFormatXML{
		Type:   "long",
		Year:   "x4 A4 x*",
		Month:  "x8 A2 x*",
		Day:    "x10 A2 x*",
		Hour:   "x12 A2 x*",
		Minute: "x14 A2 x*",
		Second: "x16 A2 x*",
	},
	"x5": DownloadNameFormatXML{
		Type:   "long",
		Year:   "x5 A4 x*",
		Month:  "x9 A2 x*",
		Day:    "x11 A2 x*",
		Hour:   "x14 A2 x*",
		Minute: "x16 A2 x*",
		Second: "x18 A2 x*",
	},
	"x6": DownloadNameFormatXML{
		Type:    "short",
		Year:    "x6 A2",
		YearDay: "x9 A3",
	},
	"x15": DownloadNameFormatXML{
		Type:   "long",
		Year:   "x15 A4 x*",
		Month:  "x20 A2 x*",
		Day:    "x23 A2 x*",
		Hour:   "x26 A2 x*",
		Minute: "x28 A2 x*",
		Second: "x30 A2 x*",
	},
	"x17": DownloadNameFormatXML{
		Type:   "long",
		Year:   "x17 A4 x*",
		Month:  "x22 A2 x*",
		Day:    "x25 A2 x*",
		Hour:   "x28 A2 x*",
		Minute: "x30 A2 x*",
		Second: "x32 A2 x*",
	},
	"x18": DownloadNameFormatXML{
		Type:   "long",
		Year:   "x18 A4 x*",
		Month:  "x23 A2 x*",
		Day:    "x26 A2 x*",
		Hour:   "x29 A2 x*",
		Minute: "x31 A2 x*",
		Second: "x33 A2 x*",
	},
	"x19": DownloadNameFormatXML{
		Type:   "long",
		Year:   "x19 A4 x*",
		Month:  "x24 A2 x*",
		Day:    "x27 A2 x*",
		Hour:   "x30 A2 x*",
		Minute: "x32 A2 x*",
		Second: "x34 A2 x*",
	},
	"x21": DownloadNameFormatXML{
		Type:   "long",
		Year:   "x21 A4 x*",
		Month:  "x26 A2 x*",
		Day:    "x29 A2 x*",
		Hour:   "x32 A2 x*",
		Minute: "x34 A2 x*",
		Second: "x36 A2 x*",
	},
}

type SessionList []meta.Session

func (s SessionList) Len() int           { return len(s) }
func (s SessionList) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s SessionList) Less(i, j int) bool { return s[i].Start.After(s[j].Start) }

func main() {

	var verbose bool
	flag.BoolVar(&verbose, "verbose", false, "make noise")

	var output string
	flag.StringVar(&output, "output", "output", "output directory")

	var operational string
	flag.StringVar(&operational, "operational", "operational.xml", "operational status file name")

	var stopped string
	flag.StringVar(&stopped, "stopped", "stopped.xml", "stopped status file name")

	var network string
	flag.StringVar(&network, "network", "../../network", "base network directory")

	var install string
	flag.StringVar(&install, "install", "../../install", "base install directory")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Build a RINEX configuration site XML file from delta meta information\n")
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
		installedAntenna[i.Mark] = append(installedAntenna[i.Mark], i)
	}
	for i, _ := range installedAntenna {
		sort.Sort(sort.Reverse(meta.InstalledAntennaList(installedAntenna[i])))
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
		sort.Sort(sort.Reverse(meta.DeployedReceiverList(deployedReceivers[i])))
	}

	var installedRadomeList meta.InstalledRadomeList
	if err := meta.LoadList(filepath.Join(install, "radomes.csv"), &installedRadomeList); err != nil {
		fmt.Fprintf(os.Stderr, "error: unable to load radome installs: %v\n", err)
		os.Exit(-1)
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
		fmt.Fprintf(os.Stderr, "error: unable to load metsensors installs: %v\n", err)
		os.Exit(-1)
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
		fmt.Fprintf(os.Stderr, "error: unable to load mark list: %v\n", err)
		os.Exit(-1)
	}

	var monumentList meta.MonumentList
	if err := meta.LoadList(filepath.Join(network, "monuments.csv"), &monumentList); err != nil {
		fmt.Fprintf(os.Stderr, "error: unable to load monuments list: %v\n", err)
		os.Exit(-1)
	}
	monuments := make(map[string]meta.Monument)
	for _, m := range monumentList {
		monuments[m.Mark] = m
	}

	var sessionList meta.SessionList
	if err := meta.LoadList(filepath.Join(install, "sessions.csv"), &sessionList); err != nil {
		fmt.Fprintf(os.Stderr, "error: unable to load session list: %v\n", err)
		os.Exit(-1)
	}

	sort.Sort(SessionList(sessionList))

	var on, off []Mark

	sessions := make(map[string][]meta.Session)
	for _, s := range sessionList {
		sessions[s.Mark] = append(sessions[s.Mark], s)
	}

	for _, m := range markList {
		if _, ok := sessions[m.Code]; !ok {
			continue
		}
		if _, ok := installedAntenna[m.Code]; !ok {
			continue
		}
		if _, ok := deployedReceivers[m.Code]; !ok {
			continue
		}

		var list []CGPSSessionXML
		for _, s := range sessions[m.Code] {
			for _, a := range installedAntenna[m.Code] {
				if a.Start.After(s.End) || a.End.Before(s.Start) {
					continue
				}

				for _, r := range deployedReceivers[m.Code] {
					if (!r.Start.Before(s.End)) || (!r.End.After(s.Start)) {
						continue
					}
					if r.Start.After(a.End) || r.End.Before(a.Start) {
						continue
					}

					radome := "NONE"
					if _, ok := installedRadomes[m.Code]; ok {
						for _, v := range installedRadomes[m.Code] {
							if v.Start.After(a.End) || v.Start.After(r.End) {
								continue
							}
							if v.End.Before(a.Start) || v.End.Before(r.Start) {
								continue
							}
							radome = v.Model
						}
					}

					var metsensor *MetSensor
					if _, ok := installedMetSensors[m.Code]; ok {
						for _, v := range installedMetSensors[m.Code] {
							if (!v.Start.Before(s.End)) || (!v.End.After(s.Start)) {
								continue
							}
							metsensor = &MetSensor{
								Model:      v.Make,
								Type:       v.Model + " S/N " + v.Serial,
								HrAccuracy: strconv.FormatFloat(v.Accuracy.Humidity, 'g', -1, 64),
								PrAccuracy: strconv.FormatFloat(v.Accuracy.Pressure, 'g', -1, 64),
								TdAccuracy: strconv.FormatFloat(v.Accuracy.Temperature, 'g', -1, 64),
							}
						}
					}

					var firmware []FirmwareHistoryXML
					if _, ok := firmwareHistory[r.Model]; ok {
						if _, ok := firmwareHistory[r.Model][r.Serial]; ok {
							for i, _ := range firmwareHistory[r.Model][r.Serial] {
								v := firmwareHistory[r.Model][r.Serial][len(firmwareHistory[r.Model][r.Serial])-i-1]
								/*
									if v.End.Before(r.Start) || v.Start.After(r.End) {
										continue
									}
								*/
								firmware = append(firmware, FirmwareHistoryXML{
									StartTime: v.Start.Format(DateTimeFormat),
									StopTime: func() string {
										if time.Now().After(v.End) {
											return v.End.Format(DateTimeFormat)
										} else {
											return "open"
										}
									}(),
									Version: v.Version,
								})
							}
						}
					}

					list = append(list, CGPSSessionXML{
						StartTime: func() string {
							if r.Start.After(s.Start) && r.Start.After(a.Start) {
								return r.Start.Format(DateTimeFormat)
							} else if a.Start.After(s.Start) {
								return a.Start.Format(DateTimeFormat)
							} else {
								return s.Start.Format(DateTimeFormat)
							}
						}(),
						StopTime: func() string {
							if r.End.Before(s.End) && r.End.Before(a.End) {
								if time.Now().After(r.End) {
									return r.End.Format(DateTimeFormat)
								} else {
									return "open"
								}
							} else if a.End.Before(s.End) {
								if time.Now().After(a.End) {
									return a.End.Format(DateTimeFormat)
								} else {
									return "open"
								}
							} else {
								if time.Now().After(s.End) {
									return s.End.Format(DateTimeFormat)
								} else {
									return "open"
								}
							}
						}(),
						Receiver: ReceiverXML{
							SerialNumber:      r.Serial,
							IGSDesignation:    r.Model,
							FirmwareHistories: firmware,
						},
						InstalledCGPSAntenna: InstalledCGPSAntennaXML{
							Height:      Number{Units: "m", Value: fmt.Sprintf("%.4f", a.Vertical)},
							OffsetNorth: Number{Units: "m", Value: fmt.Sprintf("%.4f", a.North)},
							OffsetEast:  Number{Units: "m", Value: fmt.Sprintf("%.4f", a.East)},
							Radome:      radome,
							CGPSAntenna: CGPSAntennaXML{
								SerialNumber:   a.Serial,
								IGSDesignation: a.Model,
							},
						},
						MetSensor: metsensor,
						ObservationInterval: Number{
							Units: "s",
							Value: fmt.Sprintf("%.0f", s.Interval.Seconds()),
						},
						Operator: OperatorXML{
							Name:   s.Operator,
							Agency: s.Agency,
						},
						Rinex: RinexXML{
							HeaderCommentName: s.HeaderComment,
							HeaderCommentText: func() string {
								if t, ok := HeaderComments[s.HeaderComment]; ok {
									return strings.Replace(strings.Join(strings.Fields(t), " "), "email", " email", -1)
								}
								return ""
							}(),
						},
						DataFormat: func() string {
							parts := strings.Fields(s.Format)
							if len(parts) > 0 {
								return parts[0]
							}
							return "unknown"
						}(),
						DownloadNameFormat: func() DownloadNameFormatXML {
							parts := strings.Fields(s.Format)
							if len(parts) > 1 {
								if f, ok := DownloadNameFormats[parts[1]]; ok {
									return f
								}
							}
							return DownloadNameFormatXML{}
						}(),
					})
				}
			}
		}

		sort.Sort(CGPSSessionXMLs(list))

		if time.Now().Before(m.End) {
			on = append(on, Mark{
				Name:    m.Name,
				Code:    m.Code,
				Lat:     strconv.FormatFloat(m.Latitude, 'f', 14, 64),
				Lon:     strconv.FormatFloat(m.Longitude, 'f', 14, 64),
				Network: m.Network,
				Opened:  m.Start.Format(DateFormat),
			})
		} else {
			off = append(off, Mark{
				Name:    m.Name,
				Code:    m.Code,
				Lat:     strconv.FormatFloat(m.Latitude, 'f', 14, 64),
				Lon:     strconv.FormatFloat(m.Longitude, 'f', 14, 64),
				Network: m.Network,
				Opened:  m.Start.Format(DateFormat),
				Closed:  m.End.Format(DateFormat),
			})
		}

		x := NewSiteXML(
			MarkXML{
				GeodeticCode: m.Code,
				DomesNumber: func() string {
					if c, ok := monuments[m.Code]; ok {
						return c.DomesNumber
					}
					return ""
				}(),
			},
			LocationXML{
				Latitude:  m.Latitude,
				Longitude: m.Longitude,
				Height:    m.Elevation,
				Datum:     m.Datum,
			},
			list,
		)

		s, err := x.Marshal()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: unable to marsh xml: %v\n", err)
			os.Exit(-1)
		}

		s = []byte(strings.Replace(string(s), "<domes-number></domes-number>", "<domes-number />", -1))

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

	s, err := (Marks{Name: "CGPS Marks. Status: Operational. ", Marks: on}).Marshal()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: unable to marsh xml: %v\n", err)
		os.Exit(-1)
	}
	s = []byte(strings.Replace(string(s), "></mark>", "/>", -1))
	xmlfile := filepath.Join(output, operational)
	if err := os.MkdirAll(filepath.Dir(xmlfile), 0755); err != nil {
		fmt.Fprintf(os.Stderr, "error: unable to create dir: %v\n", err)
		os.Exit(-1)
	}
	if err := ioutil.WriteFile(xmlfile, s, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "error: unable to write file: %v\n", err)
		os.Exit(-1)
	}

	s, err = (Marks{Name: "CGPS Marks. Status: Stopped. ", Marks: off}).Marshal()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: unable to marsh xml: %v\n", err)
		os.Exit(-1)
	}
	s = []byte(strings.Replace(string(s), "></mark>", "/>", -1))
	xmlfile = filepath.Join(output, stopped)
	if err := os.MkdirAll(filepath.Dir(xmlfile), 0755); err != nil {
		fmt.Fprintf(os.Stderr, "error: unable to create dir: %v\n", err)
		os.Exit(-1)
	}
	if err := ioutil.WriteFile(xmlfile, s, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "error: unable to write file: %v\n", err)
		os.Exit(-1)
	}
}
