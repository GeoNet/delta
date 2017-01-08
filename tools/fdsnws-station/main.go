package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/ozym/fdsn/stationxml"
)

func matchString(str string, patterns []string) bool {
	if !(len(patterns) > 0) {
		return true
	}
	for _, pattern := range patterns {
		if hit, err := filepath.Match(pattern, str); err == nil && hit {
			return true
		}
	}

	return false
}

func buildPatterns(u *url.URL, params ...string) []string {
	var patterns []string
	for _, p := range params {
		for _, x := range u.Query()[p] {
			for _, y := range strings.Split(x, ",") {
				patterns = append(patterns, strings.TrimSpace(y))
			}
		}
	}
	return patterns
}

func buildLevelMatcher(level string) func(string) bool {
	switch level {
	case "network":
		return func(s string) bool {
			switch s {
			case "network":
				return true
			default:
				return false
			}
		}
	case "station":
		return func(s string) bool {
			switch s {
			case "network", "station":
				return true
			default:
				return false
			}
		}
	case "channel":
		return func(s string) bool {
			switch s {
			case "network", "station", "channel":
				return true
			default:
				return false
			}
		}
	default:
		return func(s string) bool {
			return true
		}
	}
}

func handleStationXML(x *stationxml.FDSNStationXML) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml; charset=utf-8")

		levelMatcher := buildLevelMatcher(strings.ToLower(r.URL.Query().Get("level")))

		var networks []stationxml.Network
		for _, n := range x.Networks {
			if levelMatcher("network") {
				if ok := matchString(n.Code, buildPatterns(r.URL, "network", "net")); !ok {
					continue
				}
				var stations []stationxml.Station
				if levelMatcher("station") {
					var channels []stationxml.Channel
					for _, s := range n.Stations {
						if ok := matchString(s.Code, buildPatterns(r.URL, "station", "sta")); !ok {
							continue
						}
						if levelMatcher("channel") {
							for _, c := range s.Channels {
								if ok := matchString(c.LocationCode, buildPatterns(r.URL, "location", "loc")); !ok {
									continue
								}
								if ok := matchString(c.Code, buildPatterns(r.URL, "channel", "cha")); !ok {
									continue
								}
								channels = append(channels, stationxml.Channel{
									BaseNode:           c.BaseNode,
									ExternalReferences: c.ExternalReferences,
									LocationCode:       c.LocationCode,
									Latitude:           c.Latitude,
									Longitude:          c.Longitude,
									Elevation:          c.Elevation,
									Depth:              c.Depth,
									Azimuth:            c.Azimuth,
									Dip:                c.Dip,
									Types:              c.Types,
									SampleRateGroup:    c.SampleRateGroup,
									StorageFormat:      c.StorageFormat,
									ClockDrift:         c.ClockDrift,
									CalibrationUnits:   c.CalibrationUnits,
									Sensor:             c.Sensor,
									PreAmplifier:       c.PreAmplifier,
									DataLogger:         c.DataLogger,
									Equipment:          c.Equipment,
									Response: func() *stationxml.Response {
										if levelMatcher("response") {
											return c.Response
										}
										return nil
									}(),
								})
							}
						}
						stations = append(stations, stationxml.Station{
							BaseNode:           s.BaseNode,
							Latitude:           s.Latitude,
							Longitude:          s.Longitude,
							Elevation:          s.Elevation,
							Site:               s.Site,
							Vault:              s.Vault,
							Geology:            s.Geology,
							Equipments:         s.Equipments,
							Operators:          s.Operators,
							CreationDate:       s.CreationDate,
							TerminationDate:    s.TerminationDate,
							ExternalReferences: s.ExternalReferences,
							Channels:           channels,
						})
					}
				}
				networks = append(networks, stationxml.Network{
					BaseNode: n.BaseNode,
					Stations: stations,
				})
			}
		}

		// render station xml
		root := stationxml.NewFDSNStationXML(x.Source, x.Sender, x.Module, "", networks)
		if ok := root.IsValid(); ok != nil {
			http.Error(w, "invalid stationxml file", http.StatusInternalServerError)
		}

		// marshal into xml
		res, err := root.Marshal()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		// done
		if _, err := w.Write(res); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func main() {

	var server string
	flag.StringVar(&server, "server", ":8899", "where to run the server: [:8899]")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Build a fdsn-station web service\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  %s [options] <station-xml-file>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fmt.Fprintf(os.Stderr, "\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")
	}

	flag.Parse()

	if !(flag.NArg() > 0) {
		log.Fatal("No station xml file given")
	}

	var root stationxml.FDSNStationXML

	xmlFile, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatal("Error opening file:", err)
	}
	defer xmlFile.Close()

	b, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		log.Fatal("Error reading file:", err)
	}

	if err := xml.Unmarshal(b, &root); err != nil {
		log.Fatal("Error loading xml file:", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", handleStationXML(&root))

	log.Printf("listening on: %s\n", server)
	if err := http.ListenAndServe(server, mux); err != nil {
		log.Fatal(err)
	}

}
