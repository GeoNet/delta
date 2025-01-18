package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/GeoNet/delta"
)

type Settings struct {
	path    string // delta database file
	network string // dart buoy network code
	output  string // optional output file
}

func main() {

	var settings Settings

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Build a dart deployment config file\n")
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

	flag.StringVar(&settings.path, "path", "delta.db", "delta database file")
	flag.StringVar(&settings.network, "network", "TD", "dart buoy network code")
	flag.StringVar(&settings.output, "output", "", "output dart configuration file")

	flag.Parse()

	db, err := delta.NewDB(settings.path)
	if err != nil {
		log.Fatalf("unable to read delta db: %v", err)
	}

	sites, err := QuerySites(db, settings.network)
	if err != nil {
		log.Fatalf("unable to query delta db: %v", err)
	}

	// avoids the json null
	deployments := make([]Deployment, 0)

	for _, site := range sites {

		// get the site detiding details
		detide, err := QueryDetide(db, site.Buoy, site.Start, site.End)
		if err != nil {
			log.Fatalf("unable to query delta db: %v", err)
		}

		// will be sorted as per delta
		deployments = append(deployments, Deployment{
			Network:          site.Network,
			Buoy:             site.Buoy,
			Location:         site.Location,
			Latitude:         site.Latitude,
			Longitude:        site.Longitude,
			Detide:           detide,
			Depth:            site.Depth,
			TimingCorrection: site.Correction,
			Start:            site.Start,
			End:              site.End,
		})

		//TODO: remove for production deployment
		// for completeness do it again!
		deployments = append(deployments, Deployment{
			Network:          site.Network,
			Buoy:             site.Buoy,
			Location:         site.Location,
			Latitude:         site.Latitude,
			Longitude:        site.Longitude,
			Detide:           detide,
			Depth:            site.Depth,
			TimingCorrection: site.Correction,
			Start:            site.Start,
			End:              site.End,
		})
	}

	switch {
	case settings.output != "":
		file, err := os.Create(settings.output)
		if err != nil {
			log.Fatalf("unable to create output file %q: %v", settings.output, err)
		}
		defer file.Close()

		if err := Encode(file, deployments); err != nil {
			log.Fatalf("unable to write output to %q: %v", settings.output, err)
		}
	default:
		if err := Encode(os.Stdout, deployments); err != nil {
			log.Fatalf("unable to write output: %v", err)
		}
	}
}
