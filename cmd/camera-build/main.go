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
	verbose bool

	baseDir string

	networks string
	output   string
	active   bool
}

func (s Settings) Networks() []string {
	var networks []string
	for _, n := range strings.Split(s.networks, ",") {
		if s := strings.TrimSpace(strings.ToUpper(n)); s != "" {
			networks = append(networks, s)
		}
	}
	return networks
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

	flag.BoolVar(&settings.verbose, "verbose", false, "make noise")
	flag.BoolVar(&settings.active, "active", false, "only output active camera information")
	flag.StringVar(&settings.baseDir, "base", "", "optional base for custom delta files")
	flag.StringVar(&settings.networks, "networks", "", "comma separated list of networks, an empty value matches all networks")
	flag.StringVar(&settings.output, "output", "", "where to store json formatted output")

	flag.Parse()

	set, err := delta.NewBase(settings.baseDir)
	if err != nil {
		log.Fatalf("unable to build delta set for %s: %v", settings.baseDir, err)
	}

	// avoids null when marshalling an empty slice
	captions := make([]Caption, 0)

	nets := make(map[string]interface{})
	for _, n := range settings.Networks() {
		nets[n] = true
	}

	for _, mount := range set.Mounts() {
		if _, ok := nets[mount.Network]; !ok && len(nets) > 0 {
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
			log.Fatal(err)
		}
	default:
		if err := Captions(captions).Encode(os.Stdout); err != nil {
			log.Fatal(err)
		}
	}
}
