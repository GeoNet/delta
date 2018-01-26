package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

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
	flag.StringVar(&stationRegexp, "stations", "[A-Z0-9]+", "regexp selection of stations")

	var channelRegexp string
	flag.StringVar(&channelRegexp, "channels", "[A-Z0-9]+", "regexp selection of channels")

	var networkRegexp string
	flag.StringVar(&networkRegexp, "networks", "[A-Z0-9]+", "regexp selection of internal networks")

	var externalRegexp string
	flag.StringVar(&externalRegexp, "external", "[A-Z0-9]+", "regexp selection of external networks")

	var sensorRegexp string
	flag.StringVar(&sensorRegexp, "sensors", ".*", "regexp selection of sensors")

	var dataloggerRegexp string
	flag.StringVar(&dataloggerRegexp, "dataloggers", ".*", "regexp selection of dataloggers")

	var installed bool
	flag.BoolVar(&installed, "installed", false, "set station times based on installation dates")

	var operational bool
	flag.BoolVar(&operational, "operational", false, "only output operational channels")

	var active bool
	flag.BoolVar(&active, "active", false, "only output stations with active channels")

	var offset time.Duration
	flag.DurationVar(&offset, "operational-offset", 0, "provide a recently closed window for operational only requests")

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
		fmt.Fprintf(os.Stderr, "  Use a \"!\" prefix to indicate that the regexp should not match the expression\n")
		fmt.Fprintf(os.Stderr, "\n")
	}

	flag.Parse()

	builder, err := NewBuilder(
		SetInstalled(installed),
		SetActive(active),
		SetOperational(operational, offset),
		SetNetworks(networkRegexp),
		SetExternal(externalRegexp),
		SetStations(stationRegexp),
		SetChannels(channelRegexp),
		SetSensors(sensorRegexp),
		SetDataloggers(dataloggerRegexp),
	)
	if err != nil {
		log.Fatalf("unable to make builder: %v", err)
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
