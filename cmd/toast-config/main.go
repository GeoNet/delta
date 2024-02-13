package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/GeoNet/delta"
)

type Settings struct {
	base    string // options delta base file directory
	agency  string // seiscomp agency id
	dart    string // DART network code
	coastal string // coastal network code
	output  string // optional output file
}

func main() {

	var settings Settings

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Build a toast config XML file\n")
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
	flag.StringVar(&settings.agency, "agency", "WEL(GNS_Test)", "seiscomp agency id")
	flag.StringVar(&settings.dart, "dart", "TD", "dart buoy network code")
	flag.StringVar(&settings.coastal, "coastal", "TG,LG", "coast tsunami gauge network code")
	flag.StringVar(&settings.output, "output", "", "output dart configuration file")

	flag.Parse()

	set, err := delta.NewBase(settings.base)
	if err != nil {
		log.Fatalf("unable to load delta base files: %v", err)
	}

	toast, err := NewToast(set, settings.agency, settings.coastal, settings.dart)
	if err != nil {
		log.Fatalf("unable to build toast config: %v", err)
	}

	switch {
	case settings.output != "":
		// output file has been given
		file, err := os.Create(settings.output)
		if err != nil {
			log.Fatalf("unable to create output file %q: %v", settings.output, err)
		}
		defer file.Close()
		if err := toast.Encode(file); err != nil {
			log.Fatalf("unable to marshal output file %q: %v", settings.output, err)
		}
	default:
		if err := toast.Encode(os.Stdout); err != nil {
			log.Fatalf("unable to marshal output: %v", err)
		}
	}
}
