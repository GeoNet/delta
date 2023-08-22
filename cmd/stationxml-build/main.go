package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/fs"
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

// update a file if it is new or its redacted contents differ.
func updateFile(path string, contents []byte, mode fs.FileMode) error {
	// write the file if it doesn't exist.
	if _, err := os.Stat(path); err != nil {
		return os.WriteFile(path, contents, mode)
	}
	// write the file if reading the existing one fails
	existing, err := os.ReadFile(path)
	if err != nil {
		return os.WriteFile(path, contents, mode)
	}
	// write the file if the redacted contents differ
	if !bytes.Equal(redacted(existing), redacted(contents)) {
		return os.WriteFile(path, contents, mode)
	}
	// skip as they have the same contents
	return nil
}

type Settings struct {
	verbose bool // output operational info
	debug   bool // output more operational info

	base string // base directory of delta files on disk
	resp string // base directory for response xml files on disk

	version     string // create a specific StationXML version
	create      bool   // add a root XML Created entry
	corrections bool   // add calculated and applied response delays and corrections

	source string // stationxml source
	sender string // stationxml sender
	module string // stationxml module

	external Matcher // regexp selection of external networks
	exclude  Matcher // regexp selection of networks to exclude
	network  Matcher // regexp selection of networks
	station  Matcher // regexp selection of stations
	location Matcher // regexp selection of locations
	channel  Matcher // regexp selection of channels
	ignore   string  // list of stations to skip

	single    bool   // produce single station xml files
	directory string // where to store station xml files
	template  string // how to name the single station xml files
	purge     bool   // remove unknown single xml files

	output  string // output xml file, use - for stdout
	changed bool   // only update existing file if a change is detected
}

func main() {

	var settings Settings

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

	flag.BoolVar(&settings.verbose, "verbose", false, "add operational info")
	flag.BoolVar(&settings.debug, "debug", false, "add extra operational info")

	flag.StringVar(&settings.base, "base", "", "base directory of delta files on disk")
	flag.StringVar(&settings.resp, "resp", "", "base directory for response xml files on disk")

	flag.StringVar(&settings.version, "version", "", "create a specific StationXML version")
	flag.BoolVar(&settings.create, "create", false, "add a root XML \"Created\" entry")
	flag.BoolVar(&settings.corrections, "corrections", false, "add calculated and applied response delays and corrections")

	flag.StringVar(&settings.source, "source", "GeoNet", "stationxml source")
	flag.StringVar(&settings.sender, "sender", "WEL(GNS_Test)", "stationxml sender")
	flag.StringVar(&settings.module, "module", "Delta", "stationxml module")

	flag.TextVar(&settings.external, "external", MustMatcher(externalRe), "regexp selection of external networks")
	flag.TextVar(&settings.exclude, "exclude", MustMatcher(excludeRe), "regexp selection of networks to exclude")
	flag.TextVar(&settings.network, "network", MustMatcher(networkRe), "regexp selection of networks")
	flag.TextVar(&settings.station, "station", MustMatcher(stationRe), "regexp selection of stations")
	flag.TextVar(&settings.location, "location", MustMatcher(locationRe), "regexp selection of locations")
	flag.TextVar(&settings.channel, "channel", MustMatcher(channelRe), "regexp selection of channels")
	flag.StringVar(&settings.ignore, "ignore", "", "list of stations to skip")

	flag.BoolVar(&settings.single, "single", false, "produce single station xml files")
	flag.StringVar(&settings.directory, "directory", "xml", "where to store station xml files")
	flag.StringVar(&settings.template, "template", "station_{{.ExternalCode}}_{{.StationCode}}.xml", "how to name the single station xml files")
	flag.BoolVar(&settings.purge, "purge", false, "remove unknown single xml files")
	flag.StringVar(&settings.output, "output", "", "output xml file, use \"-\" for stdout")
	flag.BoolVar(&settings.changed, "changed", false, "only update existing file if a change is detected")

	flag.Func("freq", "response frequency (e.g B:1.0)", func(s string) error {
		freq, err := NewFrequency(s)
		if err != nil {
			return err
		}
		freqs.Set(freq.Prefix, freq.Value)
		return nil
	})

	flag.Parse()

	switch {
	case settings.changed && (settings.output == "" || settings.output == "-"):
		log.Fatalf("invalid \"changed\" option, requires an output file to be given")
	case settings.single && (settings.output != "" && settings.output != "-"):
		log.Fatalf("invalid \"single\" option, implies an empty output file should be given")
	}

	// set recovers the delta tables
	set, err := delta.NewBase(settings.base)
	if err != nil {
		log.Fatal(err)
	}

	// simple skip list of stations
	skip := make(map[string]interface{})
	for _, s := range strings.Split(settings.ignore, ",") {
		if v := strings.TrimSpace(s); v != "" {
			skip[v] = true
		}
	}

	// builder is used to manage response files
	builder := NewBuilder(settings.resp, settings.corrections, freqs)

	// placenames is a delta utility table to geographically name stations
	placenames := meta.PlacenameList(set.Placenames())

	// remember individual stations in case of single file output
	var singles []string

	// find the external network codes to process
	exts := make(map[string][]string)
	for _, n := range set.Networks() {
		if !settings.external.MatchString(n.External) {
			if settings.debug {
				log.Printf("debug: skip network %q, doesn't match external regexp %q", n.Code, n.External)
			}
			continue
		}
		if !settings.network.MatchString(n.Code) {
			if settings.debug {
				log.Printf("debug: skip network %q, doesn't match network regexp", n.Code)
			}
			continue
		}
		if settings.exclude.MatchString(n.Code) {
			if settings.debug {
				log.Printf("debug: skip network %q, matches exclude regexp", n.Code)
			}
			continue
		}
		if settings.debug {
			log.Printf("debug: add network %q to external network %q", n.Code, n.External)
		}
		exts[n.External] = append(exts[n.External], n.Code)
	}

	// keep the external network operational times constant even if not all stations are used.
	// however, only stations which sensor installations are taken into account.
	installed := make(map[string]interface{})
	for _, sensor := range set.InstalledSensors() {
		installed[sensor.Station] = true
	}

	spans := make(map[string]meta.Span)
	for _, ext := range set.Networks() {
		if ext.Code != ext.External {
			continue
		}
		var span meta.Span
		for _, net := range set.Networks() {
			if net.External != ext.Code {
				continue
			}
			for _, stn := range set.Stations() {
				if stn.Network != net.Code {
					continue
				}
				if _, ok := installed[stn.Code]; !ok {
					continue
				}
				if span.Start.IsZero() || stn.Span.Start.Before(span.Start) {
					span.Start = stn.Span.Start
				}
				if span.End.IsZero() || stn.Span.End.After(span.End) {
					span.End = stn.Span.End
				}
			}
		}
		if span.Start.IsZero() || span.End.IsZero() {
			continue
		}
		spans[ext.Code] = span
	}

	// find a map of stations that match
	stns := make(map[string]meta.Station)
	for _, s := range set.Stations() {
		if !settings.station.MatchString(s.Code) {
			if settings.debug {
				log.Printf("debug: skip station %q, doesn't match station regexp", s.Code)
			}
			continue
		}
		if !settings.network.MatchString(s.Network) {
			if settings.debug {
				log.Printf("debug: skip station %q, doesn't match network regexp", s.Code)
			}
			continue
		}
		if _, ok := skip[s.Code]; ok {
			if settings.debug {
				log.Printf("debug: skip station %q, matches skip list", s.Code)
			}
			continue
		}
		if settings.debug {
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
			if settings.debug {
				log.Printf("debug: skip missing external network %q", n)
			}
			continue
		}

		// networks are gathered, but are mainly used for thier properties, e.g. restrictions
		var networks []stationxml.Network
		for _, lookup := range codes {

			net, ok := set.Network(lookup)
			if !ok {
				if settings.debug {
					log.Printf("debug: skip missing network %q", lookup)
				}
				continue
			}

			var stations []stationxml.Station
			for _, stn := range stns {
				if stn.Network != lookup {
					continue
				}

				var channels []stationxml.Channel
				for _, site := range set.Sites() {
					if site.Station != stn.Code {
						continue
					}

					if !settings.location.MatchString(site.Location) {
						if settings.debug {
							log.Printf("debug: skip location %q of station %q, doesn't match location regexp", site.Location, site.Station)
						}
						continue
					}

					var streams []stationxml.Stream

					// a collection joins any installed sensors with dataloggers
					for _, collection := range set.Collections(site) {
						if !settings.channel.MatchString(collection.Code()) {
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
									Description:      collection.Channel.Description(),
									Manufacturer:     collection.Channel.Make,
									Model:            collection.DeployedDatalogger.Model,
									SerialNumber:     collection.DeployedDatalogger.Serial,
									InstallationDate: collection.DeployedDatalogger.Start,
									RemovalDate:      collection.DeployedDatalogger.End,
								},
								Sensor: stationxml.Equipment{
									Type:             collection.Component.Type,
									Description:      collection.Component.Description(),
									Manufacturer:     collection.Component.Make,
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
						if settings.debug {
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
				if settings.debug {
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

		// lookup the external network time span
		span := spans[ext.Code]

		// build a stationxml shadow external structure
		externals = append(externals, stationxml.External{
			Code:        ext.Code,
			Description: ext.Description,
			Restricted:  ext.Restricted,

			StartDate: span.Start,
			EndDate:   span.End,

			Networks: networks,
		})
	}

	// build a stationxml shadow root structure
	root := stationxml.Root{
		Source: settings.source,
		Sender: settings.sender,
		Module: settings.module,
		Create: settings.create,

		Externals: externals,
	}

	switch {
	case settings.single:
		// for single file output, first build the file name, then extract a root shadow, and then encode it.
		tmpl, err := template.New("single").Parse(settings.template)
		if err != nil {
			log.Fatalf("unable to parse single xml file template: %v", err)
		}

		// keep track of files in the single directory, in case they need purging
		files := make(map[string]string)
		if err := filepath.Walk(settings.directory, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			files[filepath.Base(path)] = path
			return nil
		}); err != nil {
			log.Fatalf("unable to walk dir %s: %v", settings.directory, err)
		}

		var count, updated int
		for _, s := range singles {
			// build a station specific root structure
			if r, ok := root.Single(s); ok {

				var name bytes.Buffer
				if err := tmpl.Execute(&name, r); err != nil {
					log.Fatalf("unable to encode single xml filename: %s", err)
				}

				path := filepath.Join(settings.directory, name.String())

				res, err := r.MarshalVersion(settings.version)
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
			if !settings.purge {
				if settings.verbose {
					log.Printf("found extra file: %s", k)
				}
				continue
			}

			if settings.verbose {
				log.Printf("removing extra file: %s", k)
			}

			if err := os.Remove(v); err != nil {
				log.Fatalf("unable to remove file %s: %v", k, err)
			}

			purged++
		}

		if settings.verbose {
			log.Printf("built %d files, updated %d, removed %d", count, updated, purged)
		}

	case settings.changed:
		var raw bytes.Buffer
		if err := root.Write(&raw, settings.version); err != nil {
			log.Fatalf("unable to encode response: %v", err)
		}
		if err := updateFile(settings.output, raw.Bytes(), 0600); err != nil {
			log.Fatalf("error: unable to update file %s: %v", settings.output, err)
		}
	case settings.output == "" || settings.output == "-":
		// using the given encoder write the stationxml to the standard output
		if err := root.Write(os.Stdout, settings.version); err != nil {
			log.Fatalf("unable to encode response: %v", err)
		}
	default:
		// using the given encoder write the stationxml to a file
		if err := root.WriteFile(settings.output, settings.version); err != nil {
			log.Fatalf("unable to encode response %s: %v", settings.output, err)
		}
	}
}
