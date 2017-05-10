package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func main() {

	var base string
	flag.StringVar(&base, "base", "../..", "delta base files")

	var channels string
	flag.StringVar(&channels, "channels", "[EBH][NH]Z", "match impact channels")

	var output string
	flag.StringVar(&output, "output", "", "output impact json file")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Build an impact json file from delta meta & response information\n")
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

	streams, err := buildStreams(base, channels)
	if err != nil {
		log.Fatalf("problem loading streams %s: %v\n", base, err)
	}

	res, err := json.MarshalIndent(streams, "", "  ")
	if err != nil {
		log.Fatalf("problem marshalling streams %s: %v\n", base, err)
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
