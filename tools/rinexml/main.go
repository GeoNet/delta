package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/GeoNet/delta/meta"
)

var HeaderComments map[string]string = map[string]string{
	"linz": `
			Data supplied by the GeoNet project.
			GeoNet is core funded by EQC and is operated by GNS on behalf of EQC and all New Zealanders.
			Contact: www.geonet.org.nz  email: info@geonet.org.nz.
		`,
	"geonet": `
			Data supplied by the GeoNet project.
			GeoNet is core funded by EQC and is operated by GNS on behalf of EQC and all New Zealanders.
			Contact: www.geonet.org.nz  email: info@geonet.org.nz.
		`,
}

var DataFormats map[string]string = map[string]string{
	"5700":  "trimble_5700",
	"NetRS": "trimble_netrs",
	"NetR9": "trimble_netr9",
}

var DownloadNameFormats map[string]DownloadNameFormatXML = map[string]DownloadNameFormatXML{
	"5700": DownloadNameFormatXML{
		Type:   "long",
		Year:   "x4 A4 x*",
		Month:  "x8 A2 x*",
		Day:    "x10 A2 x*",
		Hour:   "x12 A2 x*",
		Minute: "x14 A2 x*",
		Second: "x16 A2 x*",
	},
	"NetRS": DownloadNameFormatXML{
		Type:   "long",
		Year:   "x4 A4 x*",
		Month:  "x8 A2 x*",
		Day:    "x10 A2 x*",
		Hour:   "x12 A2 x*",
		Minute: "x14 A2 x*",
		Second: "x16 A2 x*",
	},
	"NetR9": DownloadNameFormatXML{
		Type:   "long",
		Year:   "x4 A4 x*",
		Month:  "x8 A2 x*",
		Day:    "x10 A2 x*",
		Hour:   "x12 A2 x*",
		Minute: "x14 A2 x*",
		Second: "x16 A2 x*",
	},
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
				if _, ok := IGSModels[a.Model]; !ok {
					fmt.Fprintf(os.Stderr, "%s: no igs designation for: %s [%s], skipping\n", m.Code, a.Model, a.Serial)
					continue
				}

				for _, r := range deployedReceivers[m.Code] {
					if r.Start.After(a.End) || r.End.Before(a.Start) {
						continue
					}
					if _, ok := IGSModels[r.Model]; !ok {
						fmt.Fprintf(os.Stderr, "%s: no igs designation for: %s [%s]", m.Code, r.Model, r.Serial)
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
							if _, ok := IGSModels[v.Model]; !ok {
								fmt.Fprintf(os.Stderr, "%s: no igs designation for: %s [%s]", m.Code, v.Model, v.Serial)
								continue
							}
							radome = IGSModels[v.Model]
						}
					}

					var firmware []FirmwareHistoryXML
					if _, ok := firmwareHistory[r.Model]; ok {
						if _, ok := firmwareHistory[r.Model][r.Serial]; ok {
							for i, _ := range firmwareHistory[r.Model][r.Serial] {
								v := firmwareHistory[r.Model][r.Serial][len(firmwareHistory[r.Model][r.Serial])-i-1]
								if v.End.Before(r.Start) || v.Start.After(r.End) {
									continue
								}
								firmware = append(firmware, FirmwareHistoryXML{
									StartTime: v.Start.Format(meta.DateTimeFormat),
									StopTime: func() string {
										if time.Now().After(v.End) {
											return v.End.Format(meta.DateTimeFormat)
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
							if r.Start.After(a.Start) {
								return r.Start.Format(meta.DateTimeFormat)
							} else if a.Start.After(s.Start) {
								return a.Start.Format(meta.DateTimeFormat)
							} else {
								return s.Start.Format(meta.DateTimeFormat)
							}
						}(),
						StopTime: func() string {
							if r.End.Before(a.End) {
								if time.Now().After(r.End) {
									return r.End.Format(meta.DateTimeFormat)
								} else {
									return "open"
								}
							} else if a.End.Before(s.End) {
								if time.Now().After(a.End) {
									return a.End.Format(meta.DateTimeFormat)
								} else {
									return "open"
								}
							} else {
								if time.Now().After(s.End) {
									return s.End.Format(meta.DateTimeFormat)
								} else {
									return "open"
								}
							}
						}(),
						Receiver: ReceiverXML{
							SerialNumber:      r.Serial,
							IGSDesignation:    IGSModels[r.Model],
							FirmwareHistories: firmware,
						},
						InstalledCGPSAntenna: InstalledCGPSAntennaXML{
							Height:      Number{Units: "m", Value: a.Height},
							OffsetNorth: Number{Units: "m", Value: a.North},
							OffsetEast:  Number{Units: "m", Value: a.East},
							Radome:      radome,
							CGPSAntenna: CGPSAntennaXML{
								SerialNumber:   a.Serial,
								IGSDesignation: IGSModels[a.Model],
							},
						},
						ObservationInterval: Number{
							Units: "s",
							Value: s.Interval.Seconds(),
						},
						Operator: OperatorXML{
							Name:   s.Operator,
							Agency: s.Agency,
						},
						Rinex: RinexXML{
							HeaderCommentName: s.HeaderComment,
							HeaderCommentText: func() string {
								if t, ok := HeaderComments[s.HeaderComment]; ok {
									return strings.Join(strings.Fields(t), " ")
								}
								return ""
							}(),
						},
						DataFormat: func() string {
							if f, ok := DataFormats[r.Model]; ok {
								return f
							}
							return "unknown"
						}(),
						DownloadNameFormat: func() DownloadNameFormatXML {
							if f, ok := DownloadNameFormats[r.Model]; ok {
								return f
							}
							return DownloadNameFormatXML{}
						}(),
					})
				}
			}
		}

		x := NewSiteXML(
			MarkXML{
				GeodeticCode: m.Code,
				//DomesNumber:  m.DomesNumber,
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

		xmlfile := filepath.Join(output, m.Code+".xml")
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
