package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/GeoNet/delta/internal/metadb"
	"github.com/GeoNet/delta/resp"
)

func main() {

	var base string
	flag.StringVar(&base, "base", "../..", "delta base files")

	var output string
	flag.StringVar(&output, "output", "keys", "output sc3 configuration directory")

	var networks string
	//TODO: this could be managed via configuration elsewhere.
	flag.StringVar(&networks, "networks", "AK,CB,CH,EC,FI,HB,KI,MN,NX,NZ,OT,RT,SC,SI,SM,SP,SX,TP,TR,WL", "comma separated list of networks to use")

	var exclude string
	flag.StringVar(&exclude, "exclude", "", "comma separated list of stations to skip (either NN_SSSS, or simply SSSS)")

	var include string
	flag.StringVar(&include, "include", "", "comma separated list of external stations to add (requires NN_SSSS)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Build a set of sc3 key files\n")
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

	check := make(map[string]interface{})
	for _, n := range strings.Split(networks, ",") {
		check[strings.TrimSpace(strings.ToUpper(n))] = true
	}

	var unmatch []string
	for _, p := range strings.Split(exclude, ",") {
		if s := func() string {
			switch r := strings.Split(strings.ToUpper(p), "_"); len(r) {
			case 2:
				return strings.TrimSpace(r[0]) + "_" + strings.TrimSpace(r[1])
			case 1:
				if strings.ContainsAny(r[0], "*?\\[") {
					return strings.TrimSpace(r[0])
				}
				return "*_" + strings.TrimSpace(r[0])
			default:
				return ""
			}
		}(); s != "" {
			unmatch = append(unmatch, s)
		}
	}

	now := time.Now().UTC()

	globals := make(map[string]Global)
	autoPicks := make(map[string]AutoPick)
	stations := make(map[string]Station)

	db := metadb.NewMetaDB(base)

	dbStations, err := db.Stations()
	if err != nil {
		log.Fatalf("unable to load metadata stations: %v", err)
	}

	for _, station := range dbStations {
		if _, ok := check[station.Network]; !ok {
			continue
		}

		installations, err := db.Installations(station.Code)
		if err != nil {
			log.Fatalf("unable to load metadata installations %s: %v", station.Code, err)
		}

		network, err := db.Network(station.Network)
		if err != nil {
			log.Fatalf("unable to load metadata network %s: %v", station.Network, err)
		}
		if network == nil {
			continue
		}

		for _, installation := range installations {
			if installation.Start.After(now) || installation.End.Before(now) {
				continue
			}

			for _, response := range resp.Streams(installation.Datalogger.Model, installation.Sensor.Model) {
				stream, err := db.StationLocationSamplingRateStartStream(
					station.Code,
					installation.Location,
					response.Datalogger.SampleRate,
					installation.Start)
				if stream == nil || err != nil {
					if err != nil {
						log.Fatalf("unable to load metadata streams %s: %v", station.Code, err)
					}
					continue
				}

				lookup := response.Channels(stream.Axial)
				for pin := range response.Components {
					if !(pin < len(lookup)) {
						continue
					}
					channel := lookup[pin]
					if len(channel) != 3 {
						continue
					}

					global := Global{
						Stream:   channel[0:2],
						Location: installation.Location,
					}
					if global.Style() == "" {
						continue
					}

					if _, ok := globals[global.Key()]; !ok {
						globals[global.Key()] = global
					}

					autoPick := AutoPick{
						Style:      global.Style(),
						Filter:     DefaultFilter,
						Correction: DefaultCorrection,
					}

					autoPicks[autoPick.Key()] = autoPick

					station := Station{
						Global:   global,
						AutoPick: autoPick,
						Code:     station.Code,
						Network:  network.External,
					}

					stations[station.Key()] = station
				}
			}
		}
	}

	var count int

	for _, g := range globals {
		if err := Store(g, output); err != nil {
			log.Fatalf("unable to store global profile %s: %v", g.Key(), err)
		}

		count++
	}

	for _, a := range autoPicks {
		if err := Store(a, output); err != nil {
			log.Fatalf("unable to store scautopick profile %s: %v", a.Key(), err)
		}

		count++
	}

	for k, s := range stations {
		for _, m := range unmatch {
			ok, err := filepath.Match(m, k)
			if err != nil {
				log.Fatalf("unable to match station %s (%s): %v", s.Key(), m, err)
			}
			if ok {
				continue
			}

			if err := Store(s, output); err != nil {
				log.Fatalf("unable to store key output %s: %v", s.Key(), err)
			}

			count++
		}
	}

	log.Printf("updated %d files", count)
}
