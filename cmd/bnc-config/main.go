package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/GeoNet/delta"
	"github.com/GeoNet/delta/internal/ntrip"
)

type Settings struct {
	base   string // delta base directory
	common string // ntrip common files directory
	input  string // ntrip input files directory
	extra  bool   // add aliases to mounts list
	output string // optional output file
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

	flag.Parse()

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
}
