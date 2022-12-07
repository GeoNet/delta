package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/GeoNet/delta"
	"github.com/GeoNet/delta/meta"

	"github.com/GeoNet/delta/internal/stationxml"
)

const ClockDrift = 0.0001

const (
	externalRe = `^(NZ)$`
	excludeRe  = `^()$`
	networkRe  = `[A-Z0-9]+`
	stationRe  = `[A-Z0-9]+`
	locationRe = `[A-Z0-9]+`
	channelRe  = `[A-Z0-9]+`
)

var created = regexp.MustCompile(`<Created>([^<]+)</Created>`)

func redacted(contents []byte) []byte {
	return created.ReplaceAll(contents, []byte("<Created>xxxxxxxxxx</Created>"))
}

// default response frequency values
var freqs Frequencies = map[string]float64{
	"V": 0.05,
	"L": 0.1,
	"B": 1.0,
	"H": 1.0,
	"S": 15.0,
	"E": 15.0,
	"":  15.0,
}

func main() {

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Build a network StationXML file from delta meta & response information\n")
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

	var verbose bool
	flag.BoolVar(&verbose, "verbose", false, "add operational info")

	var debug bool
	flag.BoolVar(&debug, "debug", false, "add extra operational info")

	var base string
	flag.StringVar(&base, "base", "", "base directory of delta files on disk")

	var resp string
	flag.StringVar(&resp, "resp", "", "base directory for response xml files on disk")

	var version string
	flag.StringVar(&version, "version", "", "create a specific StationXML version")

	var create bool
	flag.BoolVar(&create, "create", false, "add a root XML \"Created\" entry")

	var source string
	flag.StringVar(&source, "source", "GeoNet", "stationxml source")

	var sender string
	flag.StringVar(&sender, "sender", "WEL(GNS_Test)", "stationxml sender")

	var module string
	flag.StringVar(&module, "module", "Delta", "stationxml module")

	var external Matcher
	flag.TextVar(&external, "external", MustMatcher(externalRe), "regexp selection of external networks")

	var exclude Matcher
	flag.TextVar(&exclude, "exclude", MustMatcher(excludeRe), "regexp selection of networks to exclude")

	var network Matcher
	flag.TextVar(&network, "network", MustMatcher(networkRe), "regexp selection of networks")

	var station Matcher
	flag.TextVar(&station, "station", MustMatcher(stationRe), "regexp selection of stations")

	var location Matcher
	flag.TextVar(&location, "location", MustMatcher(locationRe), "regexp selection of locations")

	var channel Matcher
	flag.TextVar(&channel, "channel", MustMatcher(channelRe), "regexp selection of channels")

	var single bool
	flag.BoolVar(&single, "single", false, "produce single station xml files")

	var directory string
	flag.StringVar(&directory, "directory", "xml", "where to store station xml files")

	var plate string
	flag.StringVar(&plate, "template", "station_{{.ExternalCode}}_{{.StationCode}}.xml", "how to name the single station xml files")

	var purge bool
	flag.BoolVar(&purge, "purge", false, "remove unknown single xml files")

	var output string
	flag.StringVar(&output, "output", "", "output xml file, use \"-\" for stdout")

	var ignore string
	flag.StringVar(&ignore, "ignore", "", "list of stations to skip")

	flag.Func("freq", "response frequency (e.g B:1.0)", func(s string) error {
		freq, err := NewFrequency(s)
		if err != nil {
			return err
		}
		freqs.Set(freq.Prefix, freq.Value)
		return nil
	})

	flag.Parse()

	// set recovers the delta tables
	set, err := delta.NewBase(base)
	if err != nil {
		log.Fatal(err)
	}

	// simple skip list of stations
	skip := make(map[string]interface{})
	for _, s := range strings.Split(ignore, ",") {
		if v := strings.TrimSpace(s); v != "" {
			skip[v] = true
		}
	}

	// builder is used to manage response files
	builder := NewBuilder(resp, freqs)

	// placenames is a delta utility table to geographically name stations
	placenames := meta.PlacenameList(set.Placenames())

	// remember individual stations in case of single file output
	var singles []string

	// find the external network codes to process
	exts := make(map[string][]string)
	for _, n := range set.Networks() {
		if !external.MatchString(n.External) {
			if debug {
				log.Printf("debug: skip network %q, doesn't match external regexp %q", n.Code, n.External)
			}
			continue
		}
		if !network.MatchString(n.Code) {
			if debug {
				log.Printf("debug: skip network %q, doesn't match network regexp", n.Code)
			}
			continue
		}
		if exclude.MatchString(n.Code) {
			if debug {
				log.Printf("debug: skip network %q, matches exclude regexp", n.Code)
			}
			continue
		}
		if debug {
			log.Printf("debug: add network %q to external network %q", n.Code, n.External)
		}
		exts[n.External] = append(exts[n.External], n.Code)
	}

	// find a map of stations that match
	stns := make(map[string]meta.Station)
	for _, s := range set.Stations() {
		if !station.MatchString(s.Code) {
			if debug {
				log.Printf("debug: skip station %q, doesn't match station regexp", s.Code)
			}
			continue
		}
		if !network.MatchString(s.Network) {
			if debug {
				log.Printf("debug: skip station %q, doesn't match network regexp", s.Code)
			}
			continue
		}
		if _, ok := skip[s.Code]; ok {
			if debug {
				log.Printf("debug: skip station %q, matches skip list", s.Code)
			}
			continue
		}
		if debug {
			log.Printf("debug: add station %q from network %q", s.Code, s.Network)
		}
		stns[s.Code] = s
	}

	// the top level stationxml networks are based on meta networks and their external codes
	var externals []stationxml.External
	for n, codes := range exts {

		// external network details
		ext, ok := set.Network(n)
		if !ok {
			if debug {
				log.Printf("debug: skip missing external network %q", n)
			}
			continue
		}

		// networks are gathered, but are mainly used for thier properties, e.g. restrictions
		var networks []stationxml.Network
		for _, lookup := range codes {

			net, ok := set.Network(lookup)
			if !ok {
				if debug {
					log.Printf("debug: skip missing network %q", lookup)
				}
				continue
			}

			var stations []stationxml.Station
			for _, stn := range stns {
				if stn.Network != lookup {
					continue
				}
				if _, ok := skip[stn.Code]; ok {
					continue
				}

				var channels []stationxml.Channel
				for _, site := range set.Sites() {
					if site.Station != stn.Code {
						continue
					}

					if !location.MatchString(site.Location) {
						if debug {
							log.Printf("debug: skip location %q of station %q, doesn't match location regexp", site.Location, site.Station)
						}
						continue
					}

					var streams []stationxml.Stream

					// a collection joins any installed sensors with dataloggers
					for _, collection := range set.Collections(site) {
						if !channel.MatchString(collection.Code()) {
							continue
						}

						// a correction adjusts a collection for site or equipment specific settings
						for _, correction := range set.Corrections(collection) {

							// recover a response encoding for the collection and adjusted correction
							r, err := builder.Response(collection, correction)
							if err != nil {
								log.Fatal(err)
							}

							// build a stationxml shadow stream structure
							streams = append(streams, stationxml.Stream{
								Code:      collection.Code(),
								StartDate: correction.Span.Start,
								EndDate:   correction.Span.End,

								SamplingRate: collection.Stream.SamplingRate,
								Triggered:    collection.Stream.Triggered,
								Types:        collection.Component.Types,

								Vertical: collection.InstalledSensor.Vertical,
								Azimuth:  collection.Azimuth(correction.Polarity),
								Dip:      collection.Dip(correction.Polarity),

								Datalogger: stationxml.Equipment{
									Type:             collection.Channel.Type,
									Description:      strings.Fields(collection.Channel.Model)[0],
									Manufacturer:     strings.Fields(collection.Channel.Make)[0],
									Model:            collection.DeployedDatalogger.Model,
									SerialNumber:     collection.DeployedDatalogger.Serial,
									InstallationDate: collection.DeployedDatalogger.Start,
									RemovalDate:      collection.DeployedDatalogger.End,
								},
								Sensor: stationxml.Equipment{
									Type:             collection.Component.Type,
									Description:      strings.Fields(collection.Component.Model)[0],
									Manufacturer:     strings.Fields(collection.Component.Make)[0],
									Model:            collection.InstalledSensor.Model,
									SerialNumber:     collection.InstalledSensor.Serial,
									InstallationDate: collection.InstalledSensor.Start,
									RemovalDate:      collection.InstalledSensor.End,
								},

								Response: r,
							})
						}
					}

					if !(len(streams) > 0) {
						if debug {
							log.Printf("debug: skip channels for location %q of station %q, no streams found", site.Location, site.Station)
						}
						continue
					}

					// build a stationxml shadow channel structure
					channels = append(channels, stationxml.Channel{
						LocationCode: site.Location,

						Latitude:  site.Latitude,
						Longitude: site.Longitude,
						Elevation: site.Elevation,
						Survey:    site.Survey,
						Datum:     site.Datum,

						Streams: streams,
					})
				}

				// use the station name, otherwise the station code.
				name := stn.Name
				if name == "" {
					name = stn.Code
				}

				// build a stationxml shadow station structure
				stations = append(stations, stationxml.Station{
					Code:        stn.Code,
					Name:        name,
					Description: placenames.Description(stn.Latitude, stn.Longitude),

					Latitude:  stn.Latitude,
					Longitude: stn.Longitude,
					Elevation: stn.Elevation,
					Datum:     stn.Datum,

					StartDate: stn.Start,
					EndDate:   stn.End,

					CreationDate:    stn.Start,
					TerminationDate: stn.End,

					Channels: channels,
				})

				singles = append(singles, stn.Code)
			}

			if !(len(stations) > 0) {
				if debug {
					log.Printf("debug: skip networks for %q, no stations found", net.Code)
				}
				continue
			}

			// build a stationxml shadow network structure
			networks = append(networks, stationxml.Network{
				Code:        net.Code,
				Description: net.Description,
				Restricted:  net.Restricted,

				Stations: stations,
			})
		}

		// build a stationxml shadow external structure
		externals = append(externals, stationxml.External{
		if !(len(networks) > 0) {
			continue
		}

		// build a stationxml shadow external structure
		externals = append(externals, stationxml.External{
			Code:        ext.Code,
			Description: ext.Description,
			Restricted:  ext.Restricted,

			Networks: networks,
		})
	}

	// build a stationxml shadow root structure
	root := stationxml.Root{
		Source: source,
		Sender: sender,
		Module: module,
		Create: create,

		Externals: externals,
	}

	switch {
	case single:
		// for single file output, first build the file name, then extract a root shadow, and then encode it.
		tmpl, err := template.New("single").Parse(plate)
		if err != nil {
			log.Fatalf("unable to parse single xml file template: %v", err)
		}

		// keep track of files in the single directory, in case they need purging
		files := make(map[string]string)
		if err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			files[filepath.Base(path)] = path
			return nil
		}); err != nil {
			log.Fatalf("unable to walk dir %s: %v", directory, err)
		}

		var count, updated int
		for _, s := range singles {
			// build a station specific root structure
			if r, ok := root.Single(s); ok {

				var name bytes.Buffer
				if err := tmpl.Execute(&name, r); err != nil {
					log.Fatalf("unable to encode single xml filename: %s", err)
				}

				path := filepath.Join(directory, name.String())

				res, err := r.MarshalVersion(version)
				if err != nil {
					log.Fatalf("unable to encode single response %s: %v", s, err)
				}

				// remove file name from purge list
				delete(files, name.String())

				// has anything changed, other than the creation time?
				if raw, err := os.ReadFile(path); err != nil {
					if bytes.Equal(redacted(raw), redacted(res)) {
						continue
					}
				}

				if err := os.WriteFile(path, res, 0600); err != nil {
					log.Fatalf("error: unable to write file %s: %v", path, err)
				}

				updated++
			}
		}

		var purged int
		for k, v := range files {
			if !purge {
				if verbose {
					log.Printf("found extra file: %s", k)
				}
				continue
			}

			if verbose {
				log.Printf("removing extra file: %s", k)
			}

			if err := os.Remove(v); err != nil {
				log.Fatalf("unable to remove file %s: %v", k, err)
			}

			purged++
		}

		if verbose {
			log.Printf("built %d files, updated %d, removed %d", count, updated, purged)
		}

	case output == "" || output == "-":
		// using the given encoder write the stationxml to the standard output
		if err := root.Write(os.Stdout, version); err != nil {
			log.Fatalf("unable to encode response: %v", err)
		}
	default:
		// using the given encoder write the stationxml to a file
		if err := root.WriteFile(output, version); err != nil {
			log.Fatalf("unable to encode response %s: %v", output, err)
		}
	}
}
