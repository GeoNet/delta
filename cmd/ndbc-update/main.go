package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/GeoNet/delta/meta"
)

type Settings struct {
	base     string // base directory of delta files on disk
	network  string // dart network
	banks    string // what payloads are configured
	platform string // dart platform
	output   string // optional output file
}

func main() {

	var settings Settings

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Build a DART update meta-data message for NOAA/NDBC\n")
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

	flag.StringVar(&settings.base, "base", "", "base directory of delta files on disk")
	flag.StringVar(&settings.network, "network", "TD", "DART network code")
	flag.StringVar(&settings.banks, "banks", "P,S", "DART buoy payload banks")
	flag.StringVar(&settings.platform, "platform", "DART 4G", "platform to encode")
	flag.StringVar(&settings.output, "output", "", "optional output file")

	flag.Parse()

	set, err := meta.NewBase(settings.base)
	if err != nil {
		log.Fatal(err)
	}

	var banks []string
	for _, s := range strings.Split(settings.banks, ",") {
		if s = strings.TrimSpace(s); s != "" {
			banks = append(banks, s)
		}
	}

	var deployments []Deployment
	for _, dart := range set.Darts() {
		station, ok := set.Station(dart.Station)
		if !ok {
			continue
		}
		if station.Network != settings.network {
			continue
		}

		for _, site := range set.Sites() {
			if site.Station != station.Code {
				continue
			}

			for _, collection := range set.Collections(site) {
				deployments = append(deployments, Deployment{
					Buoy:       dart.Station,
					Deployment: site.Location,
					Name:       station.Name,
					Pid:        dart.Pid,
					Region:     dart.WmoIdentifier,
					Banks:      banks,
					Platform:   settings.platform,
					Latitude:   site.Latitude,
					Longitude:  site.Longitude,
					Serial:     collection.InstalledSensor.Serial,
					Depth:      site.Depth,
					Start:      collection.Start,
					End:        collection.End,
				})
			}
		}
	}

	sort.Slice(deployments, func(i, j int) bool {
		switch {
		case deployments[i].Buoy < deployments[j].Buoy:
			return true
		case deployments[i].Buoy > deployments[j].Buoy:
			return false
		case deployments[i].Start.Before(deployments[j].Start):
			return true
		default:
			return false
		}
	})

	res, err := Deployments(deployments).Process()
	if err != nil {
		log.Fatal(err)
	}

	switch {
	case settings.output != "":
		if err := os.WriteFile(settings.output, res, 0600); err != nil {
			log.Fatal(err)
		}
	default:
		fmt.Println(string(res))
	}
}
