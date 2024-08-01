package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/GeoNet/delta"
	"github.com/GeoNet/delta/meta"
	"github.com/GeoNet/delta/resp"
)

// Settings holds the application configuration.
type Settings struct {
	base      string // optional delta base directory
	resp      string // optional delta response directory
	network   string // geomag network code
	stations  string // geomag station code override
	locations string // geomag location codes
	output    string // optional output file
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
	flag.StringVar(&settings.stations, "stations", "SMHS", "geomag station code override")
	flag.StringVar(&settings.locations, "locations", "50,51", "geomag location codes")
	flag.StringVar(&settings.output, "output", "", "output geomag configuration file")

	flag.Parse()

	// add extra stations to process
	stations := make(map[string]interface{})
	for _, s := range strings.Split(settings.stations, ",") {
		if s = strings.TrimSpace(s); s != "" {
			stations[s] = true
		}
	}

	// restrict locations to ones used for geomag
	locations := make(map[string]interface{})
	for _, s := range strings.Split(settings.locations, ",") {
		if s = strings.TrimSpace(s); s != "" {
			locations[s] = true
		}
	}

	// delta details
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

	var sites []meta.Site
	for _, stn := range set.Stations() {
		// must match the network code
		if stn.Network != settings.network {
			continue
		}
		for _, site := range set.Sites() {
			// must match the station code
			if site.Station != stn.Code {
				continue
			}
			// must match the location code
			if _, ok := locations[site.Location]; !ok {
				continue
			}

			sites = append(sites, site)
		}
	}

	for s := range stations {
		stn, ok := set.Station(s)
		if !ok {
			continue
		}

		for _, site := range set.Sites() {
			// must match the station code
			if site.Station != stn.Code {
				continue
			}
			// must match the location code
			if _, ok := locations[site.Location]; !ok {
				continue
			}

			sites = append(sites, site)
		}
	}

	// check each site, skip any that don't match the network
	for _, site := range sites {

		stn, ok := set.Station(site.Station)
		if !ok {
			continue
		}

		net, ok := set.Network(stn.Network)
		if !ok {
			continue
		}

		// examine the collection of information for each site
		for _, collection := range set.Collections(site) {

			// find any corrections that might be needed, e.g. gain or calibration
			for _, correction := range set.Corrections(collection) {

				pair := resp.NewInstrumentResponse()

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
				case collection.Component.SamplingRate != 0:

					// handle instruments with a single response configution (usually no datalogger)
					derived, err := resp.LookupBase(settings.resp, collection.Component.Response)
					if err != nil {
						log.Fatalf("unable to find response %q: %v", collection.Component.Response, err)
					}

					// generate the derived response
					r, err := pair.Derived(derived)
					if err != nil {
						log.Fatalf("unable to find derived response %q: %v", collection.Component.Response, err)
					}

					if r.InstrumentSensitivity != nil {
						configs = append(configs, Config{
							Srcname:     strings.Join([]string{net.External, site.Station, site.Location, collection.Code()}, "_"),
							Network:     net.External,
							Station:     site.Station,
							Location:    site.Location,
							Channel:     collection.Code(),
							Latitude:    site.Latitude,
							Longitude:   site.Longitude,
							Elevation:   site.Elevation,
							ScaleFactor: 1.0 / r.InstrumentSensitivity.Value,
							ScaleBias:   0.0,
							InputUnits:  r.InstrumentSensitivity.InputUnits.Name,
							OutputUnits: r.InstrumentSensitivity.OutputUnits.Name,
							Start:       correction.Start,
							End:         correction.End,
						})
					}
				default:
					// handle instruments with a normal response configution (e.g. sensor and datalogger)
					sensor, err := resp.LookupBase(settings.resp, collection.Component.Response)
					if err != nil {
						log.Fatalf("unable to find response %q: %v", collection.Component.Response, err)
					}
					if err := pair.SetSensor(sensor); err != nil {
						log.Fatalf("unable to set sensor response %q: %v", collection.Component.Response, err)
					}

					datalogger, err := resp.LookupBase(settings.resp, collection.Channel.Response)
					if err != nil {
						log.Fatalf("unable to find response %q: %v", collection.Channel.Response, err)
					}
					if err := pair.SetDatalogger(datalogger); err != nil {
						log.Fatalf("unable to set datalogger response %q: %v", collection.Component.Response, err)
					}

					r, err := pair.ResponseType()
					if err != nil {
						log.Fatalf("unable to find response type: %v", err)
					}

					if r.InstrumentSensitivity != nil {
						configs = append(configs, Config{
							Srcname:     strings.Join([]string{net.External, site.Station, site.Location, collection.Code()}, "_"),
							Network:     net.External,
							Station:     site.Station,
							Location:    site.Location,
							Channel:     collection.Code(),
							Latitude:    site.Latitude,
							Longitude:   site.Longitude,
							Elevation:   site.Elevation,
							ScaleFactor: 1.0 / r.InstrumentSensitivity.Value,
							ScaleBias:   0.0,
							InputUnits:  r.InstrumentSensitivity.InputUnits.Name,
							OutputUnits: r.InstrumentSensitivity.OutputUnits.Name,
							Start:       correction.Start,
							End:         correction.End,
						})
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
							Srcname:     strings.Join([]string{net.External, site.Station, site.Location, collection.Code()}, "_"),
							Network:     net.External,
							Station:     site.Station,
							Location:    site.Location,
							Channel:     collection.Code(),
							Latitude:    site.Latitude,
							Longitude:   site.Longitude,
							Elevation:   site.Elevation,
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
