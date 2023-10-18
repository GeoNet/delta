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

type Settings struct {
	base     string // optional base directory for delta files
	output   string // optional output file name
	networks string // networks to configure
	active   bool   // only output active camera details
}

func main() {

	var settings Settings

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

	flag.BoolVar(&settings.active, "active", false, "only output active camera information")
	flag.StringVar(&settings.base, "base", "", "optional base directory for delta files")
	flag.StringVar(&settings.networks, "networks", "", "comma separated list of networks, an empty value matches all networks")
	flag.StringVar(&settings.output, "output", "", "where to store json formatted output")

	flag.Parse()

	nets := make(map[string]interface{})
	for _, n := range strings.Split(settings.networks, ",") {
		if s := strings.TrimSpace(strings.ToUpper(n)); s != "" {
			nets[s] = true
		}
	}

	set, err := delta.NewBase(settings.base)
	if err != nil {
		log.Fatalf("unable to build delta set for %s: %v", settings.base, err)
	}

	// avoids null when marshalling an empty slice
	captions := make([]Caption, 0)

	for _, mount := range set.Mounts() {
		if _, ok := nets[mount.Network]; len(nets) > 0 && !ok {
			continue
		}
		if t := time.Since(mount.End); settings.active && t > 0 {
			continue
		}

		for _, view := range set.Views() {
			if view.Mount != mount.Code {
				continue
			}
			if t := time.Since(view.End); settings.active && t > 0 {
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
	case settings.output != "":
		if err := Captions(captions).EncodeFile(settings.output); err != nil {
			log.Fatalf("unable to write to output file %q: %v", settings.output, err)
		}
	default:
		if err := Captions(captions).Encode(os.Stdout); err != nil {
			log.Fatalf("unable to write output: %v", err)
		}
	}
}
