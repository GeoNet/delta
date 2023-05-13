package main

/*
Build template files for configuring a caps system based on site details.
*/

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/GeoNet/delta"
)

func main() {

	var settings Settings

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Build CAPS configuration files from delta meta information\n")
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

	flag.StringVar(&settings.baseDir, "base", "", "delta base files")
	flag.IntVar(&settings.daysGrace, "grace", 30, "number of days grace for stations with recent changes")
	flag.StringVar(&settings.outputDir, "output", "key", "output caps configuration directory")
	flag.BoolVar(&settings.purgeFiles, "purge", false, "remove unknown single xml files")
	flag.Func("network", "add specific network(s), will skip all others (use ! prefix to exclude specific network)", func(s string) error {
		for _, s := range strings.Split(s, ",") {
			if s := strings.TrimSpace(s); (len(s) > 1) && !strings.HasPrefix(s, "!") {
				settings.includeNetworks = append(settings.includeNetworks, strings.ToUpper(s))
			}
			if s := strings.TrimSpace(s); (len(s) > 2) && strings.HasPrefix(s, "!") {
				settings.excludeNetworks = append(settings.excludeNetworks, strings.TrimLeft(strings.ToUpper(s), "!"))
			}
		}
		return nil
	})
	flag.Func("station", "add specific station(s) (requires SSSS, NN_SSSS, *_SSSS, or NN_*) (use ! prefix to exclude specific station) ", func(s string) error {
		for _, s := range strings.Split(s, ",") {
			if s := strings.TrimSpace(s); !strings.HasPrefix(s, "!") {
				settings.includeStations = append(settings.includeStations, NewStation(s))
			}
			if s := strings.TrimSpace(s); strings.HasPrefix(s, "!") {
				settings.excludeStations = append(settings.excludeStations, NewStation(strings.TrimLeft(s, "!")))
			}
		}
		return nil
	})
	flag.Func("extra", "add extra station(s) which may be outside of delta (requires exact NN_SSSS)", func(s string) error {
		for _, s := range strings.Split(s, ",") {
			settings.extraStations = append(settings.extraStations, NewStation(strings.TrimSpace(s)))
		}
		return nil
	})

	flag.Parse()

	// reference time with a possible grace period
	now := time.Now().UTC()
	grace := now.AddDate(0, 0, -settings.daysGrace)

	// set recovers the delta tables
	set, err := delta.NewBase(settings.baseDir)
	if err != nil {
		log.Fatal(err)
	}

	// where the keys will be stored
	if err := os.MkdirAll(settings.outputDir, 0755); err != nil {
		log.Fatalf("unable to create output directory %s: %v", settings.outputDir, err)
	}

	// if purging, then gather a list of existing files that can be ticked off
	files := make(map[string]string)
	if settings.purgeFiles {
		paths, err := settings.Walk()
		if err != nil {
			log.Fatalf("unable to walk dir %s: %v", settings.outputDir, err)
		}
		for _, p := range paths {
			files[filepath.Base(p)] = p
		}
	}

	list := make(map[string]Station)

	for _, station := range set.Stations() {

		// may need to check network codes
		if settings.ExcludeNetwork(station.Network) {
			continue
		}

		// the network needs to exist
		network, ok := set.Network(station.Network)
		if !ok {
			continue
		}

		// check station sites
		for _, site := range set.Sites() {
			if site.Station != station.Code {
				continue
			}

			for _, collection := range set.Collections(site) {

				// must be currently operational
				if collection.Span.Start.After(now) {
					continue
				}
				if collection.Span.End.Before(grace) {
					continue
				}

				s := Station{
					Code:    station.Code,
					Network: network.External,
				}

				list[s.Key()] = s
			}
		}
	}

	var count int

	// remove any stations that would be excluded
	for k, s := range list {
		if settings.ExcludeStation(s) {
			continue
		}

		if err := s.Output(settings.outputDir); err != nil {
			log.Fatalf("unable to store key output %s: %v", k, err)
		}

		// remove from purge list
		delete(files, s.Path())

		count++
	}

	// add any extra stations outside of delta
	for _, s := range settings.extraStations {

		if err := s.Output(settings.outputDir); err != nil {
			log.Fatalf("unable to store key output %s: %v", s, err)
		}

		// remove from purge list
		delete(files, s.Path())

		count++
	}

	for k, v := range files {
		log.Printf("removing extra file: %s", k)
		if err := os.Remove(v); err != nil {
			log.Fatalf("unable to remove file %s: %v", k, err)
		}
	}

	log.Printf("updated %d files, removed %d", count, len(files))
}
