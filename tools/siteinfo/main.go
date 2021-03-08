package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/GeoNet/delta/meta"
)

func main() {

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Build a GAMIT station info file from delta meta information\n")
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

	var output string
	flag.StringVar(&output, "output", "station.info.geonet", "output info file")

	var base string
	flag.StringVar(&base, "base", "../..", "base delta directory")

	flag.Parse()

	var firmwareHistoryList meta.FirmwareHistoryList
	if err := meta.LoadList(filepath.Join(base, "install", "firmware.csv"), &firmwareHistoryList); err != nil {
		log.Fatalf("error: unable to load firmware history: %v\n", err)
	}

	firmwareHistory := make(map[string]map[string][]meta.FirmwareHistory)
	for _, i := range firmwareHistoryList {
		if _, ok := firmwareHistory[i.Model]; !ok {
			firmwareHistory[i.Model] = make(map[string][]meta.FirmwareHistory)
		}
		firmwareHistory[i.Model][i.Serial] = append(firmwareHistory[i.Model][i.Serial], i)
	}

	for j := range firmwareHistory {
		for k := range firmwareHistory[j] {
			sort.Sort(meta.FirmwareHistoryList(firmwareHistory[j][k]))
		}
	}

	var installedAntennaList meta.InstalledAntennaList
	if err := meta.LoadList(filepath.Join(base, "install", "antennas.csv"), &installedAntennaList); err != nil {
		log.Fatalf("error: unable to load antenna installs: %v\n", err)
	}

	installedAntenna := make(map[string][]meta.InstalledAntenna)
	for _, i := range installedAntennaList {
		installedAntenna[i.Mark] = append(installedAntenna[i.Mark], i)
	}
	for i := range installedAntenna {
		sort.Sort(sort.Reverse(meta.InstalledAntennaList(installedAntenna[i])))
	}

	var deployedReceiverList meta.DeployedReceiverList
	if err := meta.LoadList(filepath.Join(base, "install", "receivers.csv"), &deployedReceiverList); err != nil {
		log.Fatalf("error: unable to load receiver installs: %v\n", err)
	}

	deployedReceivers := make(map[string][]meta.DeployedReceiver)
	for _, i := range deployedReceiverList {
		deployedReceivers[i.Mark] = append(deployedReceivers[i.Mark], i)
	}
	for i := range deployedReceivers {
		sort.Sort(sort.Reverse(meta.DeployedReceiverList(deployedReceivers[i])))
	}

	var markList meta.MarkList
	if err := meta.LoadList(filepath.Join(base, "network", "marks.csv"), &markList); err != nil {
		log.Fatalf("error: unable to load mark list: %v\n", err)
	}

	var monumentList meta.MonumentList
	if err := meta.LoadList(filepath.Join(base, "network", "monuments.csv"), &monumentList); err != nil {
		log.Fatalf("error: unable to load monuments list: %v\n", err)
	}
	monuments := make(map[string]meta.Monument)
	for _, m := range monumentList {
		monuments[m.Mark] = m
	}

	var installedRadomeList meta.InstalledRadomeList
	if err := meta.LoadList(filepath.Join(base, "install", "radomes.csv"), &installedRadomeList); err != nil {
		fmt.Fprintf(os.Stderr, "error: unable to load radome installs: %v\n", err)
		os.Exit(-1)
	}

	installedRadomes := make(map[string][]meta.InstalledRadome)
	for _, i := range installedRadomeList {
		installedRadomes[i.Mark] = append(installedRadomes[i.Mark], i)
	}
	for i := range installedRadomes {
		sort.Sort(meta.InstalledRadomeList(installedRadomes[i]))
	}

	var sessionList meta.SessionList
	if err := meta.LoadList(filepath.Join(base, "install", "sessions.csv"), &sessionList); err != nil {
		log.Fatalf("error: unable to load session list: %v\n", err)
	}

	sort.Slice(sessionList, func(i, j int) bool {
		return sessionList[i].Start.After(sessionList[j].Start)
	})

	sessions := make(map[string][]meta.Session)
	for _, s := range sessionList {
		sessions[s.Mark] = append(sessions[s.Mark], s)
	}

	var info []Info
	var update time.Time

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

		for _, s := range sessions[m.Code] {
			for _, r := range deployedReceivers[m.Code] {
				if r.Start.After(s.End) || r.End.Before(s.Start) {
					continue
				}

				var firmware []FirmwareHistory
				if _, ok := firmwareHistory[r.Model]; ok {
					if _, ok := firmwareHistory[r.Model][r.Serial]; ok {
						for i := range firmwareHistory[r.Model][r.Serial] {
							v := firmwareHistory[r.Model][r.Serial][len(firmwareHistory[r.Model][r.Serial])-i-1]
							firmware = append(firmware, FirmwareHistory{
								Start:   v.Start,
								End:     v.End,
								Version: v.Version,
							})
						}
					}
				}

				sort.Slice(firmware, func(i, j int) bool {
					return firmware[i].Less(firmware[j])
				})

				for _, a := range installedAntenna[m.Code] {
					if a.Start.After(s.End) || a.End.Before(s.Start) {
						continue
					}
					if a.Start.After(r.End) || a.End.Before(r.Start) {
						continue
					}

					for _, x := range firmware {
						if x.Start.After(s.End) || x.End.Before(s.Start) {
							continue
						}
						if x.Start.After(r.End) || x.End.Before(r.Start) {
							continue
						}
						if x.Start.After(r.End) || x.End.Before(r.Start) {
							continue
						}

						start, end := s.Start, s.End
						if r.Start.After(start) {
							start = r.Start
						}

						if r.End.Before(end) {
							end = r.End
						}

						if a.Start.After(start) {
							start = a.Start
						}

						if a.End.Before(end) {
							end = a.End
						}

						if x.Start.After(start) {
							start = x.Start
						}

						if x.End.Before(end) {
							end = x.End
						}

						if end.Before(start) {
							continue
						}

						radome := "NONE"
						if _, ok := installedRadomes[m.Code]; ok {
							for _, v := range installedRadomes[m.Code] {
								if v.Start.After(end) {
									continue
								}
								if v.End.Before(start) {
									continue
								}
								radome = v.Model
							}
						}

						if update.Before(start) {
							update = start
						}

						if update.Before(end) && (time.Since(end) > 0) {
							update = end
						}

						info = append(info, Info{
							Site:           m.Code,
							Name:           m.Name,
							Start:          start,
							End:            end,
							ReceiverType:   r.Model,
							ReceiverSerial: r.Serial,
							AntennaType:    a.Model,
							AntennaSerial:  a.Serial,
							Version:        x.Version,
							Software: func() string {
								if f, ok := firmwareMap[x.Version]; ok {
									return f
								}
								if strings.Contains(x.Version, "-") {
									return strings.ReplaceAll(x.Version, "-", "")
								}
								return x.Version
							}(),
							HeightCod: "DHARP",
							Height:    a.Vertical,
							North:     a.North,
							East:      a.East,
							Radome:    radome,
						})
					}
				}
			}
		}
	}

	sort.Slice(info, func(i, j int) bool {
		switch {
		case info[i].Site < info[j].Site:
			return true
		case info[i].Site > info[j].Site:
			return false
		case info[i].Start.Before(info[j].Start):
			return true
		default:
			return false
		}
	})

	list := squashInfo(info...)

	if err := os.MkdirAll(filepath.Dir(output), 0755); err != nil {
		log.Fatal(err)
	}

	file, err := os.Create(output)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if err := encodeInfo(file, update, list...); err != nil {
		log.Fatal(err)
	}
}
