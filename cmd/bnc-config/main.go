package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/GeoNet/delta"
	"github.com/GeoNet/delta/internal/ntrip"
)

type Settings struct {
	base        string   // delta base directory
	common      string   // ntrip common files directory
	input       string   // ntrip input files directory
	extra       bool     // add aliases to mounts list
	output      string   // optional output file
	sklPath     string   // optional path to write skeleton files
	skippingSkl []string // optional list of marks to skip skeleton file generation
}

func main() {

	var settings Settings

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Provide BNC configuration file hiera settings\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  %s [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "General Options:\n")
		fmt.Fprintf(os.Stderr, "\n")
		flag.PrintDefaults()
	}

	flag.StringVar(&settings.base, "base", "", "delta base directory for config files")
	flag.StringVar(&settings.common, "common", "", "ntrip common csv file directory")
	flag.StringVar(&settings.input, "input", "", "ntrip input csv config file directory")
	flag.BoolVar(&settings.extra, "extra", false, "add aliases to mounts list")
	flag.StringVar(&settings.output, "output", "", "optional output file")
	flag.StringVar(&settings.sklPath, "skl", "", "optional path to write skeleton files")
	var skippingSkl string
	flag.StringVar(&skippingSkl, "skip-skl", "", "optional comma separated list of marks to skip skeleton file generation")

	flag.Parse()

	if skippingSkl != "" {
		settings.skippingSkl = strings.Split(skippingSkl, ",")
	}

	// set recovers the delta tables
	set, err := delta.NewBase(settings.base)
	if err != nil {
		log.Fatal(err)
	}

	caster, err := ntrip.NewCaster(settings.common, settings.input)
	if err != nil {
		log.Fatal(err)
	}

	// generate the configuration structures
	config, err := NewConfig(set, caster, settings.extra)
	if err != nil {
		log.Fatalf("unable to build config: %v", err)
	}

	// sort to help with merging
	config.Sort()

	// update the configuration yaml file
	switch {
	case settings.output != "":
		if err := config.WriteFile(settings.output); err != nil {
			log.Fatalf("unable to write config file %s: %v", settings.output, err)
		}
	default:
		if err := config.Write(os.Stdout); err != nil {
			log.Fatalf("unable to write config: %v", err)
		}
	}

	// generate skeleton file for each mount
	if settings.sklPath != "" {
		t := time.Now().UTC().Unix() // skeleton needs a reference time for the installations
		for _, m := range config.Mounts {
			if !slices.Contains(settings.skippingSkl, m.Mark) {
				if s, err := skeleton(m.Mark, set, t); err == nil {
					if len(s) > 0 { // empty means no valid mark during the reference time, simply skip without error
						err = os.WriteFile(filepath.Join(settings.sklPath, fmt.Sprintf("%s00NZL.SKL", m.Mark)), []byte(s), 0600)
						if err != nil {
							log.Fatalf("couldn't write skeleton file: %s", err)
						}
					}
				} else {
					log.Fatalf("couldn't create skeleton: %s", err)
				}
			}
		}
	}
}
