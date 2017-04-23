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
	flag.StringVar(&base, "base", "../..", "delta base files")

	var dir string
	flag.StringVar(&dir, "dir", "/work/chart/amplitude", "base file directory")

	var config string
	flag.StringVar(&config, "config", "chart-amplitude.yaml", "input config yaml file")

	var output string
	flag.StringVar(&output, "output", "", "output altus xml file")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Build a chart amplitude XML file from delta meta & response information\n")
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
		fmt.Fprintf(os.Stderr, "problem reading config %s: %v\n", config, err)
		os.Exit(1)
	}

	amplitudes, err := buildAmplitudes(cfgs, base, dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "problem building amplitudes %s: %v\n", base, err)
		os.Exit(1)
	}

	res, err := encodeAmplitudes(amplitudes)
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
