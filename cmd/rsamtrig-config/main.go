package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/GeoNet/delta"
	"github.com/GeoNet/delta/meta"
)

type Settings struct {
	base     string  // optional delta base directory
	level    float64 // default rsam level
	output   string  // output caps key directory
	stations string  // list of stations to exclude
	channels string  // list of stations to include
	purge    bool    // remove unknown station files
}

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

	flag.StringVar(&settings.base, "base", "", "optional delta base directory")
	flag.Float64Var(&settings.level, "level", 100, "provide a default trigger level")
	flag.StringVar(&settings.output, "output", "key", "output rsamtrig configuration directory")
	flag.StringVar(&settings.stations, "stations", "MAVZ:330,TRVZ:130,WHVZ:140,WSRZ:2780", "comma separated list of stations to configure")
	flag.StringVar(&settings.channels, "channels", "HHZ,EHZ", "comma separated list of channels to configure")
	flag.BoolVar(&settings.purge, "purge", false, "remove unknown single xml files")

	flag.Parse()

	set, err := delta.NewBase(settings.base)
	if err != nil {
		log.Fatalf("unable to load delta set: %v", err)
	}

	files := make(map[string]string)

	// create output directory if missing, needed for walking.
	if err := os.MkdirAll(settings.output, 0700); err != nil {
		log.Fatalf("unable to create output directory %s: %v", settings.output, err)
	}

	// find all files under the current output, needed for purging.
	if err := filepath.Walk(settings.output, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		switch d := filepath.Base(filepath.Dir(path)); d {
		case "global", "rsamtrig":
			files[filepath.Join(d, filepath.Base(path))] = path
		default:
			files[filepath.Base(path)] = path
		}
		return nil
	}); err != nil {
		log.Fatalf("unable to walk dir %s: %v", settings.output, err)
	}

	// process each station, expected label as "STN:LEVEL", if the level is
	// missing then a default value will be used.
	check := make(map[string]float64)
	for _, s := range strings.Split(settings.stations, ",") {
		if s = strings.TrimSpace(strings.ToUpper(s)); s == "" {
			continue
		}
		switch {
		case strings.Contains(s, ":"):
			parts := strings.Fields(strings.ReplaceAll(s, ":", " "))
			if len(parts) < 2 {
				continue
			}
			level, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				log.Fatalf("unable to decode station: %q", s)
			}
			check[parts[0]] = level

		default:
			check[s] = settings.level
		}
	}

	// only process the given set of channels
	channels := make(map[string]interface{})
	for _, s := range strings.Split(settings.channels, ",") {
		channels[strings.TrimSpace(strings.ToUpper(s))] = true
	}

	now := time.Now().UTC()

	globals := make(map[string]Global)
	rsamTrigs := make(map[string]RsamTrig)
	stationMap := make(map[string]Station)

	lookup := make(map[string]meta.Network)
	for _, n := range set.Networks() {
		lookup[n.Code] = n
	}

	// run through each station and its configuration
	for _, station := range set.Stations() {
		level, ok := check[station.Code]
		if !ok {
			continue
		}
		network, ok := lookup[station.Network]
		if !ok {
			continue
		}

		for _, site := range set.Sites() {
			if site.Station != station.Code {
				continue
			}

			for _, collection := range set.Collections(site) {
				// only interested in operational streams
				if collection.Start.After(now) || collection.End.Before(now) {
					continue
				}
				// only interested in configured channels
				if _, ok := channels[collection.Code()]; !ok {
					continue
				}

				key := ToKey(network.External, station.Code)

				// check that the lowest location code is being used
				if s, ok := stationMap[key]; ok {
					if !(collection.InstalledSensor.Location < s.Global.Location) {
						continue
					}
				}

				// set a global configuration (unique on channel code and location code)
				global := Global{
					Stream:   collection.Code(),
					Location: collection.Stream.Location,
				}

				if global.Style() == "" {
					continue
				}

				if _, ok := globals[global.Key()]; !ok {
					globals[global.Key()] = global
				}

				// build the rsam confgituration settings
				rsamTrig := RsamTrig{
					Station:      station.Code,
					Name:         DefaultName,
					Filter:       DefaultFilter,
					Base:         level,
					Location:     collection.Stream.Location,
					SamplingRate: collection.Stream.SamplingRate,
				}

				// first time through
				if _, ok := rsamTrigs[rsamTrig.Key()]; !ok {
					rsamTrigs[rsamTrig.Key()] = rsamTrig
				}

				// check that the lowest location code is being used
				// with the highest sampling rate.
				if r, ok := rsamTrigs[rsamTrig.Key()]; ok {
					switch {
					case collection.Stream.Location < r.Location:
						rsamTrigs[rsamTrig.Key()] = rsamTrig
					case collection.Stream.SamplingRate > r.SamplingRate:
						rsamTrigs[rsamTrig.Key()] = rsamTrig
					}
				}

				station := Station{
					Global:   global,
					RsamTrig: rsamTrig,
					Code:     station.Code,
					Network:  network.External,
				}

				stationMap[station.Key()] = station
			}
		}
	}

	var count int

	// store the generated global files
	for _, g := range globals {
		if err := Store(g, settings.output); err != nil {
			log.Fatalf("unable to store global profile %s: %v", g.Key(), err)
		}

		delete(files, g.Path())

		count++
	}

	// store the generated rsam trig files
	for _, r := range rsamTrigs {
		if err := Store(r, settings.output); err != nil {
			log.Fatalf("unable to store scautopick profile %s: %v", r.Key(), err)
		}

		delete(files, r.Path())

		count++
	}

	// store the generated station files
	for _, s := range stationMap {

		if err := Store(s, settings.output); err != nil {
			log.Fatalf("unable to store key output %s: %v", s.Key(), err)
		}

		delete(files, s.Path())

		count++
	}

	// clean unknown files if desired
	switch {
	case settings.purge:
		for k, v := range files {
			log.Printf("removing extra file: %s", k)
			if err := os.Remove(v); err != nil {
				log.Fatalf("unable to remove file %s: %v", k, err)
			}
		}
		log.Printf("updated %d files, removed %d", count, len(files))
	default:
		log.Printf("updated %d files", count)
	}

}
