package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

type Settings struct {
	match  regexp.Regexp
	input  string
	output string
}

func main() {

	var settings Settings

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Build a simplified instrument response file\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  %s [options] <file>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fmt.Fprintf(os.Stderr, "\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")
	}

	flag.TextVar(&settings.match, "match", regexp.MustCompile(`^[ESHBL][HN][ZNE12]$`), "provide a regexp match for input channels")
	flag.StringVar(&settings.input, "input", "", "provide an input StationXML file")
	flag.StringVar(&settings.output, "output", "", "provide an output JSON simplified response file")

	flag.Parse()

	root, err := StationXML(settings.input)
	if err != nil {
		log.Fatal(err)
	}

	var responses Responses

	// run down the stationxml structure
	for _, n := range root.Network {
		for _, s := range n.Station {
			for _, c := range s.Channel {
				if !settings.match.MatchString(c.Code) {
					continue
				}

				if c.Response == nil || c.Response.InstrumentSensitivity == nil {
					continue
				}

				srcname := strings.Join([]string{
					n.Code, s.Code, c.LocationCode, c.Code,
				}, "_")

				responses.AddStream(srcname, c)
			}
		}
	}

	switch {
	case settings.output != "":
		if err := responses.WriteFile(settings.output); err != nil {
			log.Fatal(err)
		}
	default:
		if err := responses.Write(os.Stdout); err != nil {
			log.Fatal(err)
		}
	}
}
