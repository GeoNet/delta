package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func main() {

	var base string
	flag.StringVar(&base, "base", "../../data", "delta base files")

	var sta string
	flag.StringVar(&sta, "stationxml", "station.xml", "network StationXML file")

	var output string
	flag.StringVar(&output, "output", "", "output altus xml file")

	var config string
	flag.StringVar(&config, "config", "chart-tsunami.yaml", "input config yaml file")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Build an altus XML file from delta meta & response information\n")
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

	cfgs, err := loadConfig(config)
	if err != nil {
		log.Fatalf("problem loading config file %s: %v", config, err)
	}

	pages, err := buildPages(cfgs, base)
	if err != nil {
		log.Fatalf("problem build pages %s: %v", base, err)
	}

	res, err := encodePages(pages)
	if err != nil {
		log.Fatalf("error: unable to marshal xml: %v", err)
	}

	// output as needed ...
	switch {
	case output != "":
		if err := os.MkdirAll(filepath.Dir(output), 0755); err != nil {
			log.Fatalf("error: unable to create directory %s: %v", filepath.Dir(output), err)
		}
		if err := ioutil.WriteFile(output, res, 0644); err != nil {
			log.Fatalf("error: unable to write file %s: %v", output, err)
		}
	default:
		os.Stdout.Write(res)
	}

}
