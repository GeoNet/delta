package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/GeoNet/delta"
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

	set, err := delta.NewBase(settings.base)
	if err != nil {
		log.Fatal(err)
	}

	var banks []string
	for _, s := range strings.Split(settings.banks, ",") {
		if s = strings.TrimSpace(s); s != "" {
			banks = append(banks, s)
		}
	}

	network, ok := set.Network(settings.network)
	if !ok {
		log.Fatalf("unable to find DART network %s", settings.network)
	}

	darts := make(map[string]meta.Dart)
	for _, d := range set.Darts() {
		darts[d.Station] = d
	}

	stations := make(map[string]meta.Station)
	for _, s := range set.Stations() {
		if s.Network != network.Code {
			continue
		}
		stations[s.Code] = s
	}

	sites := make(map[string][]meta.Site)
	for _, s := range set.Sites() {
		if _, ok := stations[s.Station]; !ok {
			continue
		}
		sites[s.Station] = append(sites[s.Station], s)
	}

	sensors := make(map[string][]meta.InstalledSensor)
	for _, s := range set.InstalledSensors() {
		if _, ok := sites[s.Station]; !ok {
			continue
		}
		sensors[s.Station] = append(sensors[s.Station], s)
	}

	var deployments []Deployment
	for k, v := range sites {
		dart, ok := darts[k]
		if !ok {
			continue
		}
		for _, s := range v {
			deployments = append(deployments, Deployment{
				Buoy:       k,
				Deployment: s.Location,
				Name:       stations[k].Name,
				Pid:        dart.Pid,
				Region:     dart.WmoIdentifier,
				Banks:      banks,
				Platform:   settings.platform,
				Latitude:   s.Latitude,
				Longitude:  s.Longitude,
				Serial:     "",
				Depth:      s.Depth,
				Start:      s.Start,
				End:        s.End,
			})
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
