package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/GeoNet/delta/internal/stationxml"

	"github.com/google/go-cmp/cmp"
)

func decode(version string, path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	res, err := stationxml.Decode(version, data)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func main() {

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Compare two StationXML files at the molecular level, or verify reading one\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  %s [options] <file> [<file>]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fmt.Fprintf(os.Stderr, "\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")
	}

	var version string
	flag.StringVar(&version, "version", "1.2", "decode a specific StationXML version")

	flag.Parse()

	switch n := flag.NArg(); {
	case n > 1:
		first, err := decode(version, flag.Args()[0])
		if err != nil {
			log.Fatal(err)
		}

		second, err := decode(version, flag.Args()[1])
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(cmp.Diff(string(first), string(second)))
	case n > 0:
		if _, err := decode(version, flag.Args()[0]); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("requires two files to compare, or one to verify")
	}
}
