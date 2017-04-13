package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/ozym/fdsn/stationxml"
)

func main() {

	var source string
	flag.StringVar(&source, "source", "GeoNet", "stationxml source")

	var sender string
	flag.StringVar(&sender, "sender", "WEL(GNS_Test)", "stationxml sender")

	var module string
	flag.StringVar(&module, "module", "Delta", "stationxml module")

	var output string
	flag.StringVar(&output, "output", "-", "output xml file")

	var base string
	flag.StringVar(&base, "base", "../..", "base of delta files on disk")

	var network string
	flag.StringVar(&network, "network", "../../network", "base network directory")

	var install string
	flag.StringVar(&install, "install", "../../install", "base installs directory")

	var stationRegexp string
	flag.StringVar(&stationRegexp, "stations", "[A-Z0-9]+", "regex selection of stations")

	var stationList string
	flag.StringVar(&stationList, "station-list", "", "regex selection of stations from file")

	var channelRegexp string
	flag.StringVar(&channelRegexp, "channels", "[A-Z0-9]+", "regex selection of channels")

	var channelList string
	flag.StringVar(&channelList, "channel-list", "", "regex selection of channels from file")

	var networkRegexp string
	flag.StringVar(&networkRegexp, "networks", "[A-Z0-9]+", "regex selection of networks")

	var networkList string
	flag.StringVar(&networkList, "network-list", "", "regex selection of networks from file")

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

	flag.Parse()

	builder := Build{
		Networks: func() *regexp.Regexp {
			re, err := Matcher(networkList, networkRegexp)
			if err != nil {
				log.Fatalf("unable to compile network matcher: %v", err)
			}
			return re
		}(),
		Stations: func() *regexp.Regexp {
			re, err := Matcher(stationList, stationRegexp)
			if err != nil {
				log.Fatalf("unable to compile network matcher: %v", err)
			}
			return re
		}(),
		Channels: func() *regexp.Regexp {
			re, err := Matcher(channelList, channelRegexp)
			if err != nil {
				log.Fatalf("unable to compile network matcher: %v", err)
			}
			return re
		}(),
	}

	// build a representation of the network
	networks, err := builder.Construct(base)
	if err != nil {
		log.Fatalf("error: unable to build networks list: %v", err)
	}

	// render station xml
	root := stationxml.NewFDSNStationXML(source, sender, module, "", networks)
	if ok := root.IsValid(); ok != nil {
		log.Fatalf("error: invalid stationxml file")
	}

	// marshal into xml
	res, err := root.Marshal()
	if err != nil {
		log.Fatalf("error: unable to marshal stationxml: %v", err)
	}

	// output as needed ...
	switch output {
	case "-":
		fmt.Fprintln(os.Stdout, string(res))
	default:
		if err := os.MkdirAll(filepath.Dir(output), 0755); err != nil {
			log.Fatalf("error: unable to create directory %s: %v", filepath.Dir(output), err)
		}
		if err := ioutil.WriteFile(output, res, 0644); err != nil {
			log.Fatalf("error: unable to write file %s: %v", output, err)
		}
	}
}
