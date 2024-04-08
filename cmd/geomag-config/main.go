package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/GeoNet/delta"
	"github.com/GeoNet/delta/internal/stationxml"
	"github.com/GeoNet/delta/resp"
)

// Settings holds the application configuration.
type Settings struct {
	base    string // optional delta base directory
	resp    string // optional delta response directory
	network string // geomag network code
	output  string // optional output file
}

func main() {

	var settings Settings

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Build a geomag processing config file\n")
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
	flag.StringVar(&settings.resp, "resp", "", "delta base files")
	flag.StringVar(&settings.network, "network", "GM", "geomag network code")
	flag.StringVar(&settings.output, "output", "", "output geomag configuration file")

	flag.Parse()

	set, err := delta.NewBase(settings.base)
	if err != nil {
		log.Fatalf("unable to create delta set: %v", err)
	}

	// avoids the json null
	configs := make([]Config, 0)

	// external network lookup
	externals := make(map[string]string)
	for _, n := range set.Networks() {
		externals[n.Code] = n.External
	}

	// network codes
	codes := make(map[string]string)
	for _, s := range set.Stations() {
		codes[s.Code] = s.Network
	}

	// check each site, skip any that don't match the network
	for _, site := range set.Sites() {
		// must have a network code
		n, ok := codes[site.Station]
		if !ok || n != settings.network {
			continue
		}

		// that code must have an external code
		external, ok := externals[n]
		if !ok {
			continue
		}

		// examine the collection of information for each site
		for _, collection := range set.Collections(site) {

			// find any corrections that might be needed, e.g. gain or calibration
			for _, correction := range set.Corrections(collection) {

				pair := stationxml.NewResponse()

				// adjust for corrections
				if cal := correction.SensorCalibration; cal != nil {
					pair.SetCalibration(cal.ScaleFactor, cal.ScaleBias, cal.ScaleAbsolute)
				}
				if gain := correction.Gain; gain != nil {
					pair.SetGain(gain.Scale.Factor, gain.Scale.Bias, gain.Absolute)
				}
				if correction.Telemetry != nil {
					pair.SetTelemetry(correction.Telemetry.ScaleFactor)
				}
				if correction.Preamp != nil {
					pair.SetPreamp(correction.Preamp.ScaleFactor)
				}

				switch {
				// handle instruments with a single response configution (usually no datalogger)
				case collection.Component.SamplingRate != 0:

					derived, err := resp.LookupBase(settings.resp, collection.Component.Response)
					if err != nil {
						log.Fatal(err)
					}

					// generate the derived response
					r, err := pair.Derived(derived)
					if err != nil {
						log.Fatal(err)
					}
					if r.InstrumentSensitivity != nil {
						configs = append(configs, Config{
							Srcname:     strings.Join([]string{external, site.Station, site.Location, collection.Code()}, "_"),
							Network:     external,
							Station:     site.Station,
							Location:    site.Location,
							Channel:     collection.Code(),
							ScaleFactor: 1.0 / r.InstrumentSensitivity.Value,
							ScaleBias:   0.0,
							InputUnits:  r.InstrumentSensitivity.InputUnits.Name,
							OutputUnits: r.InstrumentSensitivity.OutputUnits.Name,
							Start:       correction.Start,
							End:         correction.End,
						})
					}
					// handle instruments with a normal response configution (e.g. sensor and datalogger)
				default:
					sensor, err := resp.LookupBase(settings.resp, collection.Component.Response)
					if err != nil {
						log.Fatal(err)
					}
					if err := pair.SetSensor(sensor); err != nil {
						log.Fatal(err)
					}

					datalogger, err := resp.LookupBase(settings.resp, collection.Channel.Response)
					if err != nil {
						log.Fatal(err)
					}
					if err := pair.SetDatalogger(datalogger); err != nil {
						log.Fatal(err)
					}

					r, err := pair.ResponseType()
					if err != nil {
						log.Fatal(err)
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
						configs = append(configs, Config{
							Srcname:     strings.Join([]string{external, site.Station, site.Location, collection.Code()}, "_"),
							Network:     external,
							Station:     site.Station,
							Location:    site.Location,
							Channel:     collection.Code(),
							ScaleFactor: factor,
							ScaleBias:   bias,
							InputUnits:  r.InstrumentPolynomial.InputUnits.Name,
							OutputUnits: r.InstrumentPolynomial.OutputUnits.Name,
							Start:       correction.Start,
							End:         correction.End,
						})
					}
				}
			}
		}
	}

	sort.Slice(configs, func(i, j int) bool {
		return configs[i].Less(configs[j])
	})

	switch {
	case settings.output != "":
		// need to have a base directory
		if err := os.MkdirAll(filepath.Dir(settings.output), 0700); err != nil {
			log.Fatalf("unable to make output directory to %q: %v", filepath.Dir(settings.output), err)
		}
		// output file has been given
		file, err := os.Create(settings.output)
		if err != nil {
			log.Fatalf("unable to create output file %q: %v", settings.output, err)
		}
		defer file.Close()

		if err := Encode(file, configs); err != nil {
			log.Fatalf("unable to write output to %q: %v", settings.output, err)
		}
	default:
		if err := Encode(os.Stdout, configs); err != nil {
			log.Fatalf("unable to write output: %v", err)
		}
	}
}
