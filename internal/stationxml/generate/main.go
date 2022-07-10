package main

import (
	"flag"
	"log"
)

func main() {

	var input string
	flag.StringVar(&input, "input", "", "input schema file")

	var remote string
	flag.StringVar(&remote, "schema", "", "schema service endpoint to download from")

	var insecure bool
	flag.BoolVar(&insecure, "insecure", false, "whether the remote site has certificate issues, use with caution")

	var ns string
	flag.StringVar(&ns, "ns", "http://www.fdsn.org/xml/station/1", "schema namespace to process")

	flag.Parse()

	schemas := NewSchemas(ns)

	switch {
	case remote != "":
		if err := schemas.Download(remote, insecure); err != nil {
			log.Fatalf("unable to download schema from %s: %v", remote, err)
		}
	case input != "":
		if err := schemas.ReadFile(input); err != nil {
			log.Fatalf("unable to read schema from input file %s: %v", input, err)
		}
	default:
		log.Fatal("no schema source found, needs either an input file or remote url")
	}
}
