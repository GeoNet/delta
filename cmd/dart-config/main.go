package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/GeoNet/delta"
)

type Settings struct {
	base    string // optional delta base directory
	network string // dart buoy network code
	band    string // dart buoy band code
	output  string // optional output file
}

func main() {

	var settings Settings

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Build a dart deployment file\n")
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

	flag.StringVar(&settings.base, "base", "", "delta base files")
	flag.StringVar(&settings.network, "network", "TD", "dart buoy network code")
	flag.StringVar(&settings.band, "band", "W", "dart buoy stream band code")
	flag.StringVar(&settings.output, "output", "", "output dart configuration file")

	flag.Parse()

	set, err := delta.NewBase(settings.base)
	if err != nil {
		log.Fatalf("unable to create delta set: %v", err)
	}

	// avoids the json null
	deployments := make([]Deployment, 0)

	externals := make(map[string]string)
	for _, n := range set.Networks() {
		externals[n.Code] = n.External
	}

	codes := make(map[string]string)
	for _, s := range set.Stations() {
		codes[s.Code] = s.Network
	}

	// check each site
	for _, s := range set.Sites() {
		// must have a network code
		n, ok := codes[s.Station]
		if !ok || n != settings.network {
			continue
		}

		// that code must have an external code
		e, ok := externals[n]
		if !ok {
			continue
		}

		for _, c := range set.Collections(s) {
			if settings.band != c.Stream.Band {
				continue
			}
			for _, x := range set.Corrections(c) {
				var correction time.Duration
				if x.Timing != nil {
					correction = x.Timing.Correction
				}
				// will be sorted as per delta
				deployments = append(deployments, Deployment{
					Network:          e,
					Buoy:             s.Station,
					Location:         s.Location,
					Latitude:         s.Latitude,
					Longitude:        s.Longitude,
					Detide:           NewDetide(set, s),
					Depth:            s.Depth,
					TimingCorrection: correction,
					Start:            c.Start,
					End:              c.End,
				})
			}
		}
	}

	switch {
	case settings.output != "":
		// need to have a base directory
		if err := os.MkdirAll(filepath.Dir(settings.output), 0700); err != nil {
			log.Fatalf("unable to make output directory to %q: %v", filepath.Dir(settings.output), err)
		}
		// output file has been given
		file, err := os.Create(settings.output)
		if err != nil {
			log.Fatalf("unable to create output file %q: %v", settings.output, err)
		}
		defer file.Close()

		if err := Encode(file, deployments); err != nil {
			log.Fatalf("unable to write output to %q: %v", settings.output, err)
		}
	default:
		if err := Encode(os.Stdout, deployments); err != nil {
			log.Fatalf("unable to write output: %v", err)
		}
	}
}
