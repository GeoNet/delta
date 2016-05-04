package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
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
		fmt.Fprintf(os.Stderr, "error: unable to create directory %s: %v\n", output, err)
		os.Exit(-1)
	}

	pod := NewPod(output)

	if err := pod.Header(); err != nil {
		fmt.Fprintf(os.Stderr, "error: unable to build POD header file: %v\n", err)
		os.Exit(-1)
	}

	for _, f := range flag.Args() {
		if verbose {
			fmt.Fprintf(os.Stderr, "processing StationXML file: %s\n", f)
		}

		x, err := ioutil.ReadFile(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "unable to read StationXML file: %s [%v]\n", f, err)
			os.Exit(1)
		}

		var s stationxml.FDSNStationXML
		if err := xml.Unmarshal(x, &s); err != nil {
			fmt.Fprintf(os.Stderr, "unable to decode StationXML file: %s [%v]\n", f, err)
			os.Exit(1)
		}

		for i, _ := range s.Networks {
			if err := pod.Network(&s.Networks[i]); err != nil {
				fmt.Fprintf(os.Stderr, "unable to build POD files: [%v]\n", err)
				os.Exit(1)
			}
		}
	}
}
