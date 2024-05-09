package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/GeoNet/delta"
)

type Settings struct {
	base    string // options delta base file directory
	dart    string // DART network code
	coastal string // coastal network code
	geomag  string // geomag network code
	enviro  string // envirosensor network code
	manual  string // manualcollect network code
	output  string // optional output file
	extra   string // add extra stations
}

func main() {

	var settings Settings

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Build a tilde domain config file\n")
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
	flag.StringVar(&settings.dart, "dart", "TD", "dart buoy network code")
	flag.StringVar(&settings.coastal, "coastal", "TG,LG", "coast tsunami gauge network code")
	flag.StringVar(&settings.enviro, "enviro", "EN", "envirosensor network code")
	flag.StringVar(&settings.manual, "manual", "MC", "manualcollect network code")
	flag.StringVar(&settings.geomag, "geomag", "GM", "geomagnetic network code")
	flag.StringVar(&settings.extra, "extra", "GM=SM_SMHS_50", "attach extra stations and locations to a network, e.g. GM=SM_SMHS_50")
	flag.StringVar(&settings.output, "output", "", "output dart configuration file")

	flag.Parse()

	extra := make(map[string][]string)
	for _, s := range strings.Split(settings.extra, ",") {
		if s = strings.TrimSpace(s); s != "" {
			keys := strings.Split(s, "=")
			if len(keys) != 2 {
				log.Fatalf("invalid extra entry %q: should be of the form \"GM=SM_SMHS_50\" with \"GM\" being network to extend", s)
			}
			extra[keys[0]] = append(extra[keys[0]], keys[1])
		}
	}

	set, err := delta.NewBase(settings.base)
	if err != nil {
		log.Fatalf("unable to load delta base files: %v", err)
	}

	var tilde Tilde

	// update dart domain
	if err := tilde.Dart(set, settings.dart); err != nil {
		log.Fatalf("unable to build dart configuration: %v", err)
	}

	// update environment sensor domain
	if err := tilde.EnviroSensor(set, settings.enviro); err != nil {
		log.Fatalf("unable to build envirosensor configuration: %v", err)
	}

	// update manual collection domain
	if err := tilde.ManualCollection(set, settings.manual); err != nil {
		log.Fatalf("unable to build manualcollect configuration: %v", err)
	}

	// update coastal domain
	if err := tilde.Coastal(set, settings.coastal); err != nil {
		log.Fatalf("unable to build coastal configuration: %v", err)
	}

	// update geomag domain
	if err := tilde.Geomag(set, settings.geomag, extra[settings.geomag]...); err != nil {
		log.Fatalf("unable to build geomag configuration: %v", err)
	}

	switch {
	case settings.output != "":
		// output file has been given
		file, err := os.Create(settings.output)
		if err != nil {
			log.Fatalf("unable to create output file %q: %v", settings.output, err)
		}
		defer file.Close()

		if err := tilde.MarshalIndent(file, "", "  "); err != nil {
			log.Fatalf("unable to marshal output file %q: %v", settings.output, err)
		}
	default:
		if err := tilde.MarshalIndent(os.Stdout, "", "  "); err != nil {
			log.Fatalf("unable to marshal output: %v", err)
		}
	}
}
