package main

/**
build template files for configuring a caps system based on site details
*/

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/GeoNet/delta/internal/metadb"
	"github.com/GeoNet/delta/resp"
)

const contents = `###
### Delivered by puppet
###
capslink:sl4caps
`

func match(s string) string {
	switch r := strings.Split(strings.ToUpper(s), "_"); len(r) {
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
}

type Station struct {
	Network string
	Code    string
}

func (s Station) Key() string {
	return fmt.Sprintf("%s_%s", strings.ToUpper(s.Network), strings.ToUpper(s.Code))
}

func (s Station) Path() string {
	return fmt.Sprintf("station_%s", s.Key())
}

func (s Station) Store(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	if err := ioutil.WriteFile(path, []byte(contents), 0644); err != nil {
		return err
	}
	return nil
}

func (s Station) Output(base string) error {
	if err := s.Store(filepath.Join(base, s.Path())); err != nil {
		return err
	}
	return nil
}

func main() {

	var base string
	flag.StringVar(&base, "base", "../..", "delta base files")

	var output string
	flag.StringVar(&output, "output", "key", "output caps configuration directory")

	var networks string
	//TODO: this could be managed via configuration elsewhere.
	flag.StringVar(&networks, "networks", "AK,CB,CH,EC,FI,HB,KI,MN,NX,NZ,OT,RT,SC,SI,SM,SP,SX,TP,TR,WL", "comma separated list of networks to use")

	var exclude string
	flag.StringVar(&exclude, "exclude", "", "comma separated list of stations to skip (either NN_SSSS, or simply SSSS)")

	var include string
	flag.StringVar(&include, "include", "", "comma separated list of external stations to add (requires NN_SSSS)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Build a set of caps key files\n")
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
		if s := match(p); s != "" {
			unmatch = append(unmatch, s)
		}
	}

	now := time.Now().UTC()

	list := make(map[string]Station)

	db := metadb.NewMetaDB(base)

	stations, err := db.Stations()
	if err != nil {
		log.Fatalf("unable to load metadata stations: %v", err)
	}

	for _, station := range stations {
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
				if err != nil {
					log.Fatalf("unable to load metadata streams %s: %v", station.Code, err)
				}
				if stream == nil {
					continue
				}

				for range response.Channels(stream.Axial) {
					//TODO: filter here, maybe regexp on channel, or look at equipment.

					s := Station{
						Code:    station.Code,
						Network: network.External,
					}

					list[s.Key()] = s
				}
			}
		}
	}

	var count int

	if err := os.MkdirAll(filepath.Dir(output), 0755); err != nil {
		log.Fatalf("unable to create output directory %s: %v", output, err)
	}

	for k, s := range list {
		for _, m := range unmatch {
			if ok, err := filepath.Match(m, k); ok || err != nil {
				if err != nil {
					log.Fatalf("unable to match station %s (%s): %v", k, m, err)
				}
				continue
			}

			if err := s.Output(output); err != nil {
				log.Fatalf("unable to store key output %s: %v", k, err)
			}

			count++
		}
	}

	for _, p := range strings.Split(include, ",") {
		r := strings.Split(strings.ToUpper(p), "_")
		if len(r) != 2 {
			continue
		}

		s := Station{
			Network: r[0],
			Code:    r[1],
		}

		if err := s.Output(output); err != nil {
			log.Fatalf("unable to store key output %s: %v", p, err)
		}

		count++
	}

	log.Printf("updated %d key files", count)
}
