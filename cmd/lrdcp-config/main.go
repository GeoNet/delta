package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/GeoNet/delta"
)

type Settings struct {
	base   string // optional delta base directory
	output string // optional output file
}

func main() {

	var settings Settings

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Build a LRDCP domain config file\n")
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

	flag.StringVar(&settings.base, "base", "", "delta base files")
	flag.StringVar(&settings.output, "output", "", "output lrdcp configuration file")

	flag.Parse()

	set, err := delta.NewBase(settings.base)
	if err != nil {
		log.Fatalf("unable to read delta set: %v", err)
	}

	var lrdcp Lrdcp

	if err := lrdcp.Load(set); err != nil {
		log.Fatalf("unable to load lrdcp config: %v", err)
	}

	switch {
	case settings.output != "":
		// output file has been given
		file, err := os.Create(settings.output)
		if err != nil {
			log.Fatalf("unable to create output file %q: %v", settings.output, err)
		}
		defer file.Close()
		if err := lrdcp.MarshalIndent(file, "", "  "); err != nil {
			log.Fatalf("unable to marshal output file %q: %v", settings.output, err)
		}
	default:
		if err := lrdcp.MarshalIndent(os.Stdout, "", "  "); err != nil {
			log.Fatalf("unable to marshal output: %v", err)
		}
	}
}
