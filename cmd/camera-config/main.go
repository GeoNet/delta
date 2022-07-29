package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/GeoNet/delta"
)

func main() {

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Build a camera caption configuration from delta meta information\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  %s [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fmt.Fprintf(os.Stderr, "\n")
		flag.PrintDefaults()
	}

	var verbose bool
	flag.BoolVar(&verbose, "verbose", false, "make noise")

	var active bool
	flag.BoolVar(&active, "active", false, "only output active camera information")

	var base string
	flag.StringVar(&base, "base", "", "optional base for custom delta files")

	var networks string
	flag.StringVar(&networks, "networks", "", "comma separated list of networks, an empty value matches all networks")

	var output string
	flag.StringVar(&output, "output", "", "where to store json formatted output")

	flag.Parse()

	nets := make(map[string]interface{})
	for _, n := range strings.Split(networks, ",") {
		if s := strings.TrimSpace(strings.ToUpper(n)); s != "" {
			nets[s] = true
		}
	}

	set, err := delta.NewBase(base)
	if err != nil {
		log.Fatalf("unable to build delta set for %s: %v", base, err)
	}

	var captions []Caption

	for _, mount := range set.Mounts() {
		if _, ok := nets[mount.Network]; len(nets) > 0 && !ok {
			continue
		}
		if t := time.Since(mount.End); active && t > 0 {
			continue
		}

		for _, view := range set.Views() {
			if view.Mount != mount.Code {
				continue
			}
			if t := time.Since(view.End); active && t > 0 {
				continue
			}
			captions = append(captions, Caption{
				Mount: view.Mount,
				View:  view.Code,
				Label: view.Label,
			})
		}
	}

	switch {
	case output != "":
		if err := Captions(captions).EncodeFile(output); err != nil {
			log.Fatal(err)
		}
	default:
		if err := Captions(captions).Encode(os.Stdout); err != nil {
			log.Fatal(err)
		}
	}
}
