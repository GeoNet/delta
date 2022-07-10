package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
)

func main() {

	var name string
	flag.StringVar(&name, "name", "stationxml", "package name")

	var datetime string
	flag.StringVar(&datetime, "datetime", "date_time.go", "file name for extra date time go code")

	var future bool
	flag.BoolVar(&future, "future", false, "output dates in the future")

	var format string
	flag.StringVar(&format, "format", "2006-01-02T15:04:05Z", "provide date time format to encode as")

	var input string
	flag.StringVar(&input, "input", "", "input schema file")

	var remote string
	flag.StringVar(&remote, "schema", "", "schema service endpoint to download from")

	var insecure bool
	flag.BoolVar(&insecure, "insecure", false, "whether the remote site has certificate issues, use with caution")

	var ns string
	flag.StringVar(&ns, "ns", "http://www.fdsn.org/xml/station/1", "schema namespace to process")

	var output string
	flag.StringVar(&output, "output", "output", "output dir")

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

	if err := os.MkdirAll(output, 0755); err != nil {
		log.Fatal(err)
	}

	settings := Datetime{
		Package: name,
		Format:  format,
		Future:  future,
	}

	path := filepath.Join(output, datetime)
	log.Printf("rendering datetime => %s", path)
	if err := settings.RenderFile(path); err != nil {
		log.Fatalf("unable to parse date time type: %v", err)
	}
}
