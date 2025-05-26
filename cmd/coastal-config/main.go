package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/GeoNet/delta"
	"github.com/GeoNet/delta/resp"
)

type Settings struct {
	base    string // optional delta base directory
	resp    string // optional delta resp directory
	network string // coastal gauge network code
	output  string // optional output file
}

func main() {

	var settings Settings

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Build a coastal config file\n")
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
	flag.StringVar(&settings.resp, "resp", "", "delta resp files")
	flag.StringVar(&settings.network, "network", "TG", "coastal gauge network code")
	flag.StringVar(&settings.output, "output", "", "output dart configuration file")

	flag.Parse()

	set, err := delta.NewBase(settings.base)
	if err != nil {
		log.Fatalf("unable to create delta set: %v", err)
	}

	// avoids the json null
	coastal := make([]Coastal, 0)

	externals := make(map[string]string)
	for _, n := range set.Networks() {
		externals[n.Code] = n.External
	}

	codes := make(map[string]string)
	for _, s := range set.Stations() {
		codes[s.Code] = s.Network
	}

	// check each site
	for _, s := range set.Sites() {
		// must have a network code
		n, ok := codes[s.Station]
		if !ok || n != settings.network {
			continue
		}

		// that code must have an external code
		e, ok := externals[n]
		if !ok {
			continue
		}

		for _, c := range set.Collections(s) {

			pair := resp.NewInstrumentResponse()

			// handle instruments with a normal response configution (e.g. sensor and datalogger)
			sensor, err := resp.LookupBase(settings.resp, c.Component.Response)
			if err != nil {
				log.Fatalf("unable to find response %q: %v", c.Component.Response, err)
			}
			if err := pair.SetSensor(sensor); err != nil {
				log.Fatalf("unable to set sensor response %q: %v", c.Component.Response, err)
			}

			datalogger, err := resp.LookupBase(settings.resp, c.Channel.Response)
			if err != nil {
				log.Fatalf("unable to find response %q: %v", c.Channel.Response, err)
			}
			if err := pair.SetDatalogger(datalogger); err != nil {
				log.Fatalf("unable to set datalogger response %q: %v", c.Component.Response, err)
			}

			r, err := pair.ResponseType()
			if err != nil {
				log.Fatalf("unable to find response type: %v", err)
			}

			if r.InstrumentPolynomial != nil {
				var factor, bias float64
				for _, c := range r.InstrumentPolynomial.Coefficients {
					switch c.Number {
					case 1:
						bias = c.Value
					case 2:
						factor = c.Value
					}
				}

				// will be sorted as per delta
				coastal = append(coastal, Coastal{
					Network:   e,
					Station:   s.Station,
					Location:  s.Location,
					Latitude:  s.Latitude,
					Longitude: s.Longitude,
					Detide:    NewDetide(set, s),
					Factor:    factor,
					Bias:      bias,
					Units:     r.InstrumentPolynomial.InputUnits.Name,
					Start:     c.Start,
					End:       c.End,
				})
			}
		}
	}

	switch {
	case settings.output != "":
		file, err := os.Create(settings.output)
		if err != nil {
			log.Fatalf("unable to create output file %q: %v", settings.output, err)
		}
		defer file.Close()

		if err := Encode(file, coastal); err != nil {
			log.Fatalf("unable to write output to %q: %v", settings.output, err)
		}
	default:
		if err := Encode(os.Stdout, coastal); err != nil {
			log.Fatalf("unable to write output: %v", err)
		}
	}
}
