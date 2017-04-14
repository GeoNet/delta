package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var aliases = map[string]string{
	"DECF": "RR#1",
	"CSBF": "RR#2",
	"BCOF": "FDLS",
	"TLED": "RR2",
	"COLD": "RR3",
	"GODS": "D12C",
	"DCZ":  "DCSM",
	"SECF": "SECS",
	"JCWJ": "JACS",
	"KVSD": "RR5",
	"PARS": "D11C",
}

func main() {

	var base string
	flag.StringVar(&base, "base", "../..", "delta base files")

	var output string
	flag.StringVar(&output, "output", "", "output altus xml file")

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

	sites, err := buildSites(base)
	if err != nil {
		fmt.Fprintf(os.Stderr, "problem loading sites %s: %v\n", base, err)
		os.Exit(1)
	}

	res, err := encodeSites(sites)
	if err != nil {
		log.Fatalf("error: unable to encode xml: %v", err)
	}

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
