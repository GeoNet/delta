package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/GeoNet/delta"
	"github.com/GeoNet/delta/meta"
)

type Settings struct {
	debug      bool   // run application in debug mode
	base       string // optional delta base directory
	output     string // output caps key directory
	configured string // configured networks
	exclude    string // list of stations to exclude
	include    string // list of stations to include
	skip       string // list of streams to skip
}

var (
	Version = "gaps-config"
	Release = "dev"
)

func main() {

	var settings Settings

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Build a set of gaps key files\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Version: %s (%s)\n", Version, Release)
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

	flag.BoolVar(&settings.debug, "debug", false, "run the application in debug mode")
	flag.StringVar(&settings.base, "base", "", "optional delta base directory")

	flag.StringVar(&settings.output, "output", "key", "output caps configuration directory")

	flag.StringVar(&settings.configured, "networks", "AK,CB,CH,EC,FI,HB,KI,NM,NZ,OT,RA,RT,SC,SI,SM,SP,TP,TR,WL,LG,TG", "comma separated list of networks to use")
	flag.StringVar(&settings.exclude, "exclude", "", "comma separated list of stations to skip (either NN_SSSS, or simply SSSS)")
	flag.StringVar(&settings.include, "include", "", "comma separated list of external stations to add (requires NN_SSSS)")
	flag.StringVar(&settings.skip, "skip", "", "comma separated list of streams to skip (requires NN_SSSS_LL)")

	flag.Parse()

	if settings.debug {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	slog.Info(Version, "release", Release)

	set, err := delta.NewBase(settings.base)
	if err != nil {
		log.Fatalf("unable to load delta set: %v", err)
	}

	if err := os.MkdirAll(settings.output, 0700); err != nil {
		log.Fatalf("unable to create output directory %s: %v", settings.output, err)
	}

	now := time.Now().UTC()

	files := make(map[string]string)

	if err := filepath.Walk(settings.output, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		switch d := filepath.Base(filepath.Dir(path)); d {
		case "global":
			files[filepath.Join(d, filepath.Base(path))] = path
		default:
			files[filepath.Base(path)] = path
		}
		return nil
	}); err != nil {
		log.Fatalf("unable to walk directory %q: %v", settings.output, err)
	}

	check := make(map[string]interface{})
	for _, n := range strings.Split(settings.configured, ",") {
		check[strings.TrimSpace(strings.ToUpper(n))] = true
	}

	extra := make(map[string]interface{})
	for _, v := range strings.Split(settings.include, ",") {
		extra[strings.TrimSpace(strings.ToUpper(v))] = true
	}

	skip := make(map[string]interface{})
	for _, v := range strings.Split(settings.skip, ",") {
		skip[strings.TrimSpace(strings.ToUpper(v))] = true
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

	globals := make(map[string]Global)
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

			stream := ToStream(network.External, station.Code, site.Location)
			if _, ok := skip[stream]; ok {
				continue
			}

			for _, collection := range set.Collections(site) {
				if collection.Start.After(now) || collection.End.Before(now) {
					continue
				}

				if v, ok := stationMap[ToKey(network.External, station.Code)]; ok {
					if collection.InstalledSensor.Location >= v.Global.Location {
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

				station := Station{
					Global:  global,
					Code:    station.Code,
					Network: network.External,
				}

				stationMap[station.Key()] = station
			}
		}
	}

	var count int
	for _, g := range globals {
		slog.Debug("gaps", slog.String("store", g.Path()))
		if err := Store(g, settings.output); err != nil {
			log.Fatalf("unable to store file %q: %v", g.Path(), err)
		}
		delete(files, g.Path())
		count++
	}

	for k, v := range stationMap {
		for _, m := range unmatch {
			ok, err := filepath.Match(m, k)
			if err != nil {
				log.Fatalf("unable to match file %q <=> %q", m, k)
			}
			if ok {
				continue
			}

			slog.Debug("gaps", slog.String("store", v.Path()))
			if err := Store(v, settings.output); err != nil {
				log.Fatalf("unable to store file %q: %v", v.Path(), err)
			}

			delete(files, v.Path())

			count++
		}
	}

	for k, v := range files {
		slog.Info("gaps", slog.String("remove", k))
		if err := os.Remove(v); err != nil {
			log.Fatalf("unable to remove file %q: %v", v, err)
		}
	}

	slog.Info("gaps", slog.Int("found", count), slog.Int("removed", len(files)))
}
