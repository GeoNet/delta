package main

import (
	"flag"
	"fmt"
	"github.com/GeoNet/delta/meta"
	"github.com/GeoNet/kit/gloria_pb"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func main() {

	//Meta Loading Code grabbed from rinexml
	var verbose bool
	flag.BoolVar(&verbose, "verbose", false, "make noise")

	var output string
	flag.StringVar(&output, "output", "output", "output directory")

	var network string
	flag.StringVar(&network, "network", "../../network", "base network directory")

	var install string
	flag.StringVar(&install, "install", "../../install", "base install directory")

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, "\n")
		fmt.Fprint(os.Stderr, "Build a set of ProtoBuf files for Gloria from delta meta information\n")
		fmt.Fprint(os.Stderr, "\n")
		fmt.Fprint(os.Stderr, "Usage:\n")
		fmt.Fprint(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  %s [options]\n", os.Args[0])
		fmt.Fprint(os.Stderr, "\n")
		fmt.Fprint(os.Stderr, "Options:\n")
		fmt.Fprint(os.Stderr, "\n")
		flag.PrintDefaults()
		fmt.Fprint(os.Stderr, "\n")
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

	for j := range firmwareHistory {
		for k := range firmwareHistory[j] {
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
	for i := range installedAntenna {
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
	for i := range deployedReceivers {
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
	for i := range installedRadomes {
		sort.Sort(meta.InstalledRadomeList(installedRadomes[i]))
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

	var marks = gloria_pb.Marks{Marks: make(map[string]*gloria_pb.Mark)}

	for _, m := range markList {

		mark_pb := gloria_pb.Mark{
			Code: m.Code,
			Point: &gloria_pb.Point{
				Latitude:  m.Latitude,
				Longitude: m.Longitude,
				Elevation: m.Elevation,
			},
			DomesNumber:      monuments[m.Code].DomesNumber,
			DeployedReceiver: make([]*gloria_pb.DeployedReceiver, 0),
			InstalledAntenna: make([]*gloria_pb.InstalledAntenna, 0),
			InstalledRadome:  make([]*gloria_pb.InstalledRadome, 0),
		}

		// Higher download priority marks are scheduled for download first.
		// Download rate can be restricted here as well.
		switch {
		case m.Igs:
			mark_pb.Download = &gloria_pb.Download{Priority: 1000}
			mark_pb.Distribution = &gloria_pb.Distribution{Igs: true}
		case m.Network == "LI":
			mark_pb.Download = &gloria_pb.Download{Priority: 100}
			mark_pb.Distribution = &gloria_pb.Distribution{Linz: true}
		default:
			mark_pb.Download = &gloria_pb.Download{Priority: 0}
		}

		if m.Network == "LI" {
			mark_pb.Comment = `Data supplied by the GeoNet project.  GeoNet is core
funded by EQC, with support from LINZ, and is
operated by GNS on behalf of EQC and all New Zealanders.
Contact: www.geonet.org.nz  email: info@geonet.org.nz`
		} else {
			mark_pb.Comment = `Data supplied by the GeoNet project.  GeoNet is core
funded by EQC and is operated by GNS on behalf of
EQC and all New Zealanders.
Contact: www.geonet.org.nz  email: info@geonet.org.nz`
		}

		recList := deployedReceivers[m.Code]
		for _, rec := range recList {
			rec_pb := gloria_pb.DeployedReceiver{
				Receiver: &gloria_pb.Receiver{
					Model:        rec.Model,
					SerialNumber: rec.Serial,
					Firmware:     make([]*gloria_pb.Firmware, 0),
				},
				Span: &gloria_pb.Span{
					Start: rec.Start.Unix(),
					End:   rec.End.Unix(),
				},
			}

			firmList := firmwareHistory[rec.Model][rec.Serial]
			for _, firm := range firmList {
				firm_pb := gloria_pb.Firmware{
					Version: firm.Version,
					Span: &gloria_pb.Span{
						Start: firm.Start.Unix(),
						End:   firm.End.Unix(),
					},
				}
				rec_pb.Receiver.Firmware = append(rec_pb.Receiver.Firmware, &firm_pb)
			}
			mark_pb.DeployedReceiver = append(mark_pb.DeployedReceiver, &rec_pb)
		}

		antList := installedAntenna[m.Code]
		for _, ant := range antList {
			ant_pb := gloria_pb.InstalledAntenna{
				Antenna: &gloria_pb.Antenna{
					Model:        ant.Model,
					SerialNumber: ant.Serial,
				},
				Offset: &gloria_pb.Offset{
					Vertical: ant.Vertical,
					North:    ant.North,
					East:     ant.East,
				},
				Span: &gloria_pb.Span{
					Start: ant.Start.Unix(),
					End:   ant.End.Unix(),
				},
			}
			mark_pb.InstalledAntenna = append(mark_pb.InstalledAntenna, &ant_pb)
		}

		radList := installedRadomes[m.Code]
		for _, rad := range radList {
			rad_pb := gloria_pb.InstalledRadome{
				Radome: &gloria_pb.Radome{
					Model: rad.Model,
				},
				Span: &gloria_pb.Span{
					Start: rad.Start.Unix(),
					End:   rad.Start.Unix(),
				},
			}
			mark_pb.InstalledRadome = append(mark_pb.InstalledRadome, &rad_pb)
		}

		b, err := proto.Marshal(&mark_pb)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: unable to marsh xml: %v\n", err)
			os.Exit(-1)
		}

		pbfile := filepath.Join(output, strings.ToLower(m.Code)+".pb")
		if err := os.MkdirAll(filepath.Dir(pbfile), 0755); err != nil {
			fmt.Fprintf(os.Stderr, "error: unable to create dir: %v\n", err)
			os.Exit(-1)
		}
		if err := ioutil.WriteFile(pbfile, b, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "error: unable to write file: %v\n", err)
			os.Exit(-1)
		}

		marks.Marks[m.Code] = &mark_pb
	}

	b, err := proto.Marshal(&marks)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: unable to marshal marks index protobuf: %v\n", err)
		os.Exit(-1)
	}
	if err := ioutil.WriteFile(filepath.Join(output, "mark-index.pb"), b, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "error: unable to write file: %v\n", err)
		os.Exit(-1)
	}
}
