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
		fmt.Fprintf(os.Stderr, "Provide NTRIP configuration file hiera settings\n")
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
	flag.StringVar(&base, "base", "", "delta base for config files")

	var input string
	flag.StringVar(&input, "input", "", "input base for ntrip config csv files")

	var output string
	flag.StringVar(&output, "output", "", "output config file")

	flag.Parse()

	set, err := delta.NewBase(base)
	if err != nil {
		log.Fatal(err)
	}

	caster, err := ntrip.NewCaster(input)
	if err != nil {
		log.Fatal(err)
	}

	config, err := Build(set, caster)
	if err != nil {
		log.Fatalf("unable to build config: %v", err)
	}

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
