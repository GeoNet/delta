package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/GeoNet/delta"
	"github.com/GeoNet/delta/internal/ntrip"
)

func main() {

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

	var base string
	flag.StringVar(&base, "base", "", "delta config base")

	var input string
	flag.StringVar(&input, "input", "", "input ntrip csv config files")

	var extra bool
	flag.BoolVar(&extra, "extra", false, "add aliases to mounts list")

	var output string
	flag.StringVar(&output, "output", "", "output config file")

	flag.Parse()

	// set recovers the delta tables
	set, err := delta.NewBase(base)
	if err != nil {
		log.Fatal(err)
	}

	caster, err := ntrip.NewCaster(input)
	if err != nil {
		log.Fatal(err)
	}

	// generate the configuration structures
	config, err := Build(set, caster, extra)
	if err != nil {
		log.Fatalf("unable to build config: %v", err)
	}

	// sort to help with merging
	config.Sort()

	// update the configuration yaml file
	switch output {
	case "":
		if err := config.Write(os.Stdout); err != nil {
			log.Fatalf("unable to write config: %v", err)
		}
	default:
		if err := config.WriteFile(output); err != nil {
			log.Fatalf("unable to write config file %s: %v", output, err)
		}
	}
}
