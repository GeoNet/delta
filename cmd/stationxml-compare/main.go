package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/google/go-cmp/cmp"

	"github.com/GeoNet/delta/internal/stationxml"
)

func decode(version, path string) ([]byte, error) {
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
		fmt.Fprintf(os.Stderr, "Compare two StationXML files at the molecular level\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  %s [options] <file> <file>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fmt.Fprintf(os.Stderr, "\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")
	}

	var version string
	flag.StringVar(&version, "version", "", "decode a specific StationXML version")

	flag.Parse()

	if flag.NArg() != 2 {
		log.Fatalf("requires two files to compare")
	}

	first, err := decode(version, flag.Args()[0])
	if err != nil {
		log.Fatal(err)
	}

	second, err := decode(version, flag.Args()[1])
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(cmp.Diff(string(first), string(second)))
}
