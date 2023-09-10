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
		files[filepath.Base(path)] = path
		return nil
	}); err != nil {
		log.Fatalf("unable to walk dir %s: %v", settings.output, err)
	}

	check := make(map[string]interface{})
	for _, n := range strings.Split(settings.configured, ",") {
		if n = strings.TrimSpace(n); n != "" {
			check[strings.ToUpper(n)] = true
		}
	}

	var unmatch []string
	for _, p := range strings.Split(settings.exclude, ",") {
		if p = match(p); p != "" {
			unmatch = append(unmatch, p)
		}
	}

	lookup := make(map[string]meta.Network)
	for _, n := range set.Networks() {
		lookup[n.Code] = n
	}

	now := time.Now().UTC()

	list := make(map[string]Station)
	for _, station := range set.Stations() {
		if _, ok := check[station.Network]; !ok {
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
				if collection.Start.After(now) || collection.End.Before(now) {
					continue
				}

				list[ToKey(network.External, station.Code)] = Station{
					Network: network.External,
					Code:    station.Code,
				}
			}
		}
	}

	var count int

	for k, s := range list {
		for _, m := range unmatch {
			if ok, err := filepath.Match(m, k); ok || err != nil {
				if err != nil {
					log.Fatalf("unable to match station %s (%s): %v", k, m, err)
				}
				continue
			}

			if err := s.Output(settings.output); err != nil {
				log.Fatalf("unable to store key output %s: %v", k, err)
			}

			delete(files, s.Path())

			count++
		}
	}

	for _, p := range strings.Split(settings.include, ",") {
		r := strings.Split(strings.ToUpper(p), "_")
		if len(r) != 2 {
			continue
		}

		s := Station{
			Network: r[0],
			Code:    r[1],
		}

		if err := s.Output(settings.output); err != nil {
			log.Fatalf("unable to store key output %s: %v", p, err)
		}

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
