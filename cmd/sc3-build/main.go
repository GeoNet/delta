package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/GeoNet/delta"
)

func main() {

	var settings Settings

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

	// application settings
	flag.StringVar(&settings.baseDir, "base", "", "delta base files")
	flag.StringVar(&settings.outputDir, "output", "key", "output sc3 configuration directory")
	flag.BoolVar(&settings.purgeFiles, "purge", false, "remove unknown single xml files")
	flag.IntVar(&settings.daysGrace, "grace", 30, "allow for a grace period in days after site changes")

	// network and station selection settings
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

	flag.Parse()

	// set recovers the delta tables
	set, err := delta.NewBase(settings.baseDir)
	if err != nil {
		log.Fatal(err)
	}

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
			files[p] = p
		}
	}

	// reference times
	now := time.Now().UTC()
	grace := now.AddDate(0, 0, -settings.daysGrace)

	globals := make(map[string]Global)
	autoPicks := make(map[string]AutoPick)
	stationMap := make(map[string]Station)

	for _, stn := range set.Stations() {

		// may need to check network codes
		if settings.ExcludeNetwork(stn.Network) {
			continue
		}

		network, ok := set.Network(stn.Network)
		if !ok {
			continue
		}

		for _, site := range set.Sites() {
			if site.Station != stn.Code {
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

				key := Station{
					Code:    stn.Code,
					Network: network.External,
				}.Key()

				if s, ok := stationMap[key]; ok {
					if !(site.Location < s.Global.Location) {
						continue
					}
				}

				global := Global{
					Stream:   collection.Code(),
					Location: site.Location,
				}
				if global.Style() == "" {
					continue
				}

				if _, ok := globals[global.Key()]; !ok {
					globals[global.Key()] = global
				}

				autoPick := AutoPick{
					Style:        global.Style(),
					Filter:       DefaultFilter,
					Correction:   DefaultCorrection,
					Location:     site.Location,
					SamplingRate: collection.Stream.SamplingRate,
				}

				if _, ok := autoPicks[autoPick.Key()]; !ok {
					autoPicks[autoPick.Key()] = autoPick
				}

				if a, ok := autoPicks[autoPick.Key()]; ok {
					switch {
					case site.Location < a.Location:
						autoPicks[autoPick.Key()] = autoPick
					case collection.Stream.SamplingRate > a.SamplingRate:
						autoPicks[autoPick.Key()] = autoPick
					}
				}

				station := Station{
					Global:   global,
					AutoPick: autoPick,
					Code:     stn.Code,
					Network:  network.External,
				}

				stationMap[station.Key()] = station
			}

		}
	}

	var count int

	// update global file settings
	for _, g := range globals {
		if err := g.Store(settings.outputDir); err != nil {
			log.Fatalf("unable to store global profile %s: %v", g.Key(), err)
		}

		delete(files, g.Path())

		count++
	}

	// update autopick file settings
	for _, a := range autoPicks {
		if err := a.Store(settings.outputDir); err != nil {
			log.Fatalf("unable to store scautopick profile %s: %v", a.Key(), err)
		}

		delete(files, a.Path())

		count++
	}

	// update station file settings
	for _, s := range stationMap {

		if settings.ExcludeStation(s) {
			continue
		}

		if err := s.Store(settings.outputDir); err != nil {
			log.Fatalf("unable to store key output %s: %v", s.Key(), err)
		}

		delete(files, s.Path())

		count++
	}

	// purge any remaining files
	for k, v := range files {
		log.Printf("removing extra file: %s", k)
		if err := os.Remove(v); err != nil {
			log.Fatalf("unable to remove file %s: %v", k, err)
		}
	}

	log.Printf("updated %d files, removed %d", count, len(files))
}
