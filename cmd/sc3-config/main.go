package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/GeoNet/delta"
	"github.com/GeoNet/delta/meta"
)

type Settings struct {
	base       string // optional delta base directory
	output     string // output caps key directory
	configured string // configured networks
	exclude    string // list of stations to exclude
	include    string // list of stations to include
	purge      bool   // remove unknown station files
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
	flag.StringVar(&settings.output, "output", "key", "output caps configuration directory")
	flag.StringVar(&settings.configured, "networks", "AK,CB,CH,EC,FI,HB,KI,NM,NZ,OT,RA,RT,SC,SI,SM,SP,TP,TR,WL,TG", "comma separated list of networks to use")
	flag.StringVar(&settings.exclude, "exclude", "", "comma separated list of stations to skip (either NN_SSSS, or simply SSSS)")
	flag.StringVar(&settings.include, "include", "", "comma separated list of external stations to add (requires NN_SSSS)")
	flag.BoolVar(&settings.purge, "purge", false, "remove unknown single xml files")

	flag.Parse()

	set, err := delta.NewBase(settings.base)
	if err != nil {
		log.Fatalf("unable to load delta set: %v", err)
	}

	if err := os.MkdirAll(settings.output, 0700); err != nil {
		log.Fatalf("unable to create output directory %s: %v", settings.output, err)
	}

	files := make(map[string]string)

	if err := filepath.Walk(settings.output, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		switch d := filepath.Base(filepath.Dir(path)); d {
		case "global", "scautopick":
			files[filepath.Join(d, filepath.Base(path))] = path
		default:
			files[filepath.Base(path)] = path
		}
		return nil
	}); err != nil {
		log.Fatalf("unable to walk dir %s: %v", settings.output, err)
	}

	check := make(map[string]interface{})
	for _, n := range strings.Split(settings.configured, ",") {
		check[strings.TrimSpace(strings.ToUpper(n))] = true
	}

	extra := make(map[string]interface{})
	for _, s := range strings.Split(settings.include, ",") {
		extra[strings.TrimSpace(strings.ToUpper(s))] = true
	}

	var unmatch []string
	for _, p := range strings.Split(settings.exclude, ",") {
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
	stationMap := make(map[string]Station)

	lookup := make(map[string]meta.Network)
	for _, n := range set.Networks() {
		lookup[n.Code] = n
	}

	for _, station := range set.Stations() {
		if _, ok := check[station.Network]; !ok {
			if _, ok := extra[station.Code]; !ok {
				continue
			}
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
				if collection.Start.After(now) || collection.End.Before(now) {
					continue
				}

				key := ToKey(network.External, station.Code)

				if s, ok := stationMap[key]; ok {
					if !(collection.InstalledSensor.Location < s.Global.Location) {
						continue
					}
				}

				var code [2]byte
				copy(code[:], collection.Code())

				global := Global{
					Stream:   string(code[:]),
					Location: collection.Stream.Location,
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
					Location:     collection.Stream.Location,
					SamplingRate: collection.Stream.SamplingRate,
				}

				if _, ok := autoPicks[autoPick.Key()]; !ok {
					autoPicks[autoPick.Key()] = autoPick
				}

				if a, ok := autoPicks[autoPick.Key()]; ok {
					switch {
					case collection.Stream.Location < a.Location:
						autoPicks[autoPick.Key()] = autoPick
					case collection.Stream.SamplingRate > a.SamplingRate:
						autoPicks[autoPick.Key()] = autoPick
					}
				}

				station := Station{
					Global:   global,
					AutoPick: autoPick,
					Code:     station.Code,
					Network:  network.External,
				}

				stationMap[station.Key()] = station
			}
		}
	}

	var count int

	for _, g := range globals {
		if err := Store(g, settings.output); err != nil {
			log.Fatalf("unable to store global profile %s: %v", g.Key(), err)
		}

		delete(files, g.Path())

		count++
	}

	for _, a := range autoPicks {
		if err := Store(a, settings.output); err != nil {
			log.Fatalf("unable to store scautopick profile %s: %v", a.Key(), err)
		}

		delete(files, a.Path())

		count++
	}

	for k, s := range stationMap {
		for _, m := range unmatch {
			ok, err := filepath.Match(m, k)
			if err != nil {
				log.Fatalf("unable to match station %s (%s): %v", s.Key(), m, err)
			}
			if ok {
				continue
			}

			if err := Store(s, settings.output); err != nil {
				log.Fatalf("unable to store key output %s: %v", s.Key(), err)
			}

			delete(files, s.Path())

			count++
		}
	}

	for k, v := range files {
		log.Printf("removing extra file: %s", k)
		if err := os.Remove(v); err != nil {
			log.Fatalf("unable to remove file %s: %v", k, err)
		}
	}

	log.Printf("updated %d files, removed %d", count, len(files))
}
