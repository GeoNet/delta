package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/GeoNet/delta"
	"github.com/GeoNet/delta/meta"
)

type Settings struct {
	base string // optional delta base directory
	resp string // optional delta resp directory

	primary string // order by constituent
	single  bool   // only a single stream per site

	exclude   regexp.Regexp // skip stations
	include   regexp.Regexp // include stations
	stations  regexp.Regexp // select stations
	networks  regexp.Regexp // select networks
	locations regexp.Regexp // select locations
	channels  regexp.Regexp // select channels
	skip      string        // select streams to exclude
	extra     string        // extra streams to include

	output string // optional output file
}

// MatchStation returns whether a stream with the given station and network codes should be included.
func (s Settings) MatchStation(net, stn string) bool {
	// include list
	if s.include.String() != "" && s.include.MatchString(stn) {
		return true
	}
	// exclude list
	if s.exclude.String() != "" && s.exclude.MatchString(stn) {
		return false
	}
	// must match network regexp
	if !s.networks.MatchString(net) {
		return false
	}
	// must match station regexp
	if !s.stations.MatchString(stn) {
		return false
	}
	return true
}

// MatchSite returns whether a stream with the given site code should be included.
func (s Settings) MatchSite(loc string) bool {
	return s.locations.MatchString(loc)
}

func main() {

	var settings Settings

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Build a chart plotting config file\n")
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

	flag.StringVar(&settings.primary, "primary", "M2", "add phase constituent for tsunami streams")
	flag.BoolVar(&settings.single, "single", false, "only add one stream per station")
	flag.StringVar(&settings.base, "base", "", "delta base files")
	flag.StringVar(&settings.resp, "resp", "", "base directory for response xml files on disk")
	flag.TextVar(&settings.exclude, "exclude", &regexp.Regexp{}, "station exclusion regexp")
	flag.TextVar(&settings.include, "include", &regexp.Regexp{}, "station inclusion regexp")
	flag.TextVar(&settings.networks, "networks", regexp.MustCompile(".*"), "network selection regexp")
	flag.TextVar(&settings.stations, "stations", regexp.MustCompile(".*"), "station selection regexp")
	flag.TextVar(&settings.locations, "locations", regexp.MustCompile(".*"), "location selection regexp")
	flag.TextVar(&settings.channels, "channels", regexp.MustCompile(".*"), "channel selection regexp")
	flag.StringVar(&settings.skip, "skip", "", "extra streams to exclude")
	flag.StringVar(&settings.extra, "extra", "", "extra streams to include")
	flag.StringVar(&settings.output, "output", "", "output chart configuration file")

	flag.Parse()

	// if only single streams are needed
	found := make(map[string]interface{})

	// build map of streams to skip
	skip := make(map[string]interface{})
	for _, stream := range strings.Split(settings.skip, ",") {
		if s := strings.TrimSpace(stream); s != "" {
			skip[s] = true
		}
	}

	set, err := delta.NewBase(settings.base)
	if err != nil {
		log.Fatalf("unable to create delta set: %v", err)
	}

	// these will override network selections
	networks := make(map[string][]meta.Station)
	stations := make(map[string]interface{})
	for _, s := range set.Stations() {
		if !settings.MatchStation(s.Network, s.Code) {
			continue
		}
		// station must be still open
		if time.Since(s.End) > 0 {
			continue
		}
		networks[s.Network] = append(networks[s.Network], s)
		stations[s.Code] = true
	}

	sites := make(map[string][]meta.Site)
	for _, s := range set.Sites() {
		// must match a station at least
		if _, ok := stations[s.Station]; !ok {
			continue
		}
		// must match expected locations
		if !settings.MatchSite(s.Location) {
			continue
		}
		// must still be operational
		if time.Since(s.End) > 0 {
			continue
		}
		sites[s.Station] = append(sites[s.Station], s)
	}

	// for extra sorting information
	lags := make(map[string]float64)
	for _, c := range set.Constituents() {
		if c.Name != settings.primary {
			continue
		}
		lags[c.Gauge] = c.Lag
	}

	// build streams, this forces an empty slice rather than a nil slice if empty.
	streams := make([]Stream, 0)

	// run through all known networks
	for _, net := range set.Networks() {
		// run through all desired stations for that given network
		for _, stn := range networks[net.Code] {
			// run through all sites for the given station
			for _, site := range sites[stn.Code] {

				// check in case of single site streams
				if _, ok := found[stn.Code]; ok && settings.single {
					continue
				}

				// run through each collection of component, channel and streams.
				for _, collection := range set.Collections(site) {

					// must be current for chart configuration
					if time.Since(collection.End) > 0 {
						continue
					}

					// must match the requested channel codes
					if !settings.channels.MatchString(collection.Code()) {
						continue
					}

					// the stream needs a sampling rate
					rate := collection.Stream.SamplingRate
					if !(rate > 0.0) {
						continue
					}

					// sampling period is the inverse of sampling rate.
					period := time.Duration((float64(time.Second) / rate))
					srcname := strings.Join([]string{
						net.External,
						site.Station,
						site.Location,
						collection.Code(),
					}, "_")

					if _, ok := skip[srcname]; ok {
						continue
					}

					// find the response from the xml snippets
					info, err := Response(settings.resp, collection)
					if err != nil {
						log.Fatalf("unable to find stream response %q: %v", srcname, err)
					}

					// build and store the stream details
					streams = append(streams, Stream{
						Srcname:            srcname,
						StationName:        stn.Name,
						NetworkDescription: net.Description,
						StationCode:        site.Station,
						LocationCode:       site.Location,
						ChannelCode:        collection.Code(),
						NetworkCode:        net.External,
						InternalNetwork:    net.Code,
						Latitude:           stn.Latitude,
						Longitude:          stn.Longitude,
						Elevation:          stn.Elevation,
						Depth:              stn.Depth,
						SamplingPeriod:     period,
						TidalLag:           lags[stn.Code],
						Sensitivity:        info.Sensitivity,
						Gain:               info.Gain,
						Bias:               info.Bias,
						InputUnits:         info.Input,
						OutputUnits:        info.Output,
					})

					// for single streams
					found[stn.Code] = true
				}
			}
		}
	}

	// there may be extra streams wanted, needs srcname and sampling rate.
	for _, extra := range strings.Split(settings.extra, ",") {

		// extract srcname and sampling rate
		split := strings.Split(strings.TrimSpace(extra), ":")
		if len(split) != 2 {
			continue
		}
		rate, err := strconv.ParseFloat(split[1], 64)
		if err != nil {
			log.Fatal(err)
		}
		if !(rate > 0.0) {
			continue
		}

		// sampling period is the inverse of sampling rate.
		period := time.Duration(float64(time.Second) / rate)

		// decode the srcname for the given parts, ignore early streams with the location missing.
		parts := strings.Split(strings.TrimSpace(split[0]), "_")
		if len(parts) != 4 {
			continue
		}

		// the station needs to be in delta at least.
		stn, ok := set.Station(parts[1])
		if !ok {
			continue
		}

		// encode the extra streams
		streams = append(streams, Stream{
			Srcname:         split[0],
			StationName:     stn.Name,
			StationCode:     parts[1],
			LocationCode:    parts[2],
			ChannelCode:     parts[3],
			NetworkCode:     parts[0],
			InternalNetwork: parts[0],
			Latitude:        stn.Latitude,
			Longitude:       stn.Longitude,
			Elevation:       stn.Elevation,
			SamplingPeriod:  period,
		})

	}

	switch {
	case settings.output != "":
		if err := Config(streams).WriteFile(settings.output); err != nil {
			log.Fatalf("unable to create output file %q: %v", settings.output, err)
		}
	default:
		if err := Config(streams).Write(os.Stdout); err != nil {
			log.Fatalf("unable to write output: %v", err)
		}
	}
}
