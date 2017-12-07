package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/ozym/fdsn/stationxml"
)

func main() {

	var verbose bool
	flag.BoolVar(&verbose, "verbose", false, "make noise")

	var output string
	flag.StringVar(&output, "output", "output", "output POD header directory")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Build POD header files from StationXML file(s)\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  %s [options] <stationxml> ...\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fmt.Fprintf(os.Stderr, "\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")
	}

	flag.Parse()

	if err := os.MkdirAll(output, 0755); err != nil {
		log.Fatalf("error: unable to create directory %s: %v", output, err)
	}

	pod := NewPod(output)

	if err := pod.Header(); err != nil {
		log.Fatalf("error: unable to build POD header file: %v", err)
	}

	for _, f := range flag.Args() {
		if verbose {
			log.Printf("processing StationXML file: %s", f)
		}

		x, err := ioutil.ReadFile(f)
		if err != nil {
			log.Fatalf("unable to read StationXML file: %s [%v]", f, err)
		}

		var s stationxml.FDSNStationXML
		if err := xml.Unmarshal(x, &s); err != nil {
			log.Fatalf("unable to decode StationXML file: %s [%v]", f, err)
		}

		for i := range s.Networks {
			if err := pod.Network(&s.Networks[i]); err != nil {
				log.Fatalf("unable to build POD files: [%v]", err)
			}
		}
	}
}
