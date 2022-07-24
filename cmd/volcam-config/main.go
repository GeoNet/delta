package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"unicode"

	"github.com/GeoNet/delta"
)

// will form the id and title for each camera
var titles = map[string]string{
	"ngauruhoe":       "Ngauruhoe",
	"whiteisland":     "White Island",
	"ruapehu":         "Ruapehu",
	"taranakiegmont":  "Taranaki/Egmont",
	"kermadecislands": "Kermadec Islands",
	"tongariro":       "Tongariro",
}

// will be used to find volcanoes from captions
var keywords = map[string][]string{
	"ngauruhoe":       {"Ngauruhoe"},
	"whiteisland":     {"White Island"},
	"ruapehu":         {"Ruapehu"},
	"taranakiegmont":  {"Taranaki"},
	"kermadecislands": {"Raoul", "Kermadec Islands"},
	"tongariro":       {"Tongariro"},
}

var letters = func(c rune) bool {
	return !unicode.IsLetter(c)
}

func main() {

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Build a volcam configuration file\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  %s [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "General Options:\n")
		fmt.Fprintf(os.Stderr, "\n")
		flag.PrintDefaults()
	}

	var base string
	flag.StringVar(&base, "base", "", "delta config files, will prefix camera and mount file names")

	var output string
	flag.StringVar(&output, "output", "", "output config file")

	flag.Parse()

	set, err := delta.NewBase(base)
	if err != nil {
		log.Fatal(err)
	}

	var volcs []Volcam
	for _, l := range set.Views() {
		if time.Since(l.Span.End) > 0 {
			continue
		}

		for _, m := range set.Mounts() {
			if time.Since(m.Span.End) > 0 {
				continue
			}
			if m.Code != l.Mount {
				continue
			}

			for _, v := range set.InstalledCameras() {
				if time.Since(v.Span.End) > 0 {
					continue
				}

				if v.Mount != m.Code {
					continue
				}

				if v.View != l.Code {
					continue
				}

				targets := make(map[string]string)
				for k, v := range keywords {
					n, ok := titles[k]
					if !ok {
						continue
					}
					for _, w := range v {
						if strings.Contains(l.Label, w) {
							targets[k] = n
						}
						if strings.Contains(l.Description, w) {
							targets[k] = n
						}
					}
				}
				var volcanoes []Volcano
				for k, w := range targets {
					volcanoes = append(volcanoes, Volcano{
						Id:    k,
						Title: w,
					})
				}

				volcs = append(volcs, Volcam{
					Id:        strings.Join(strings.FieldsFunc(strings.ToLower(l.Label), letters), ""),
					Mount:     v.Mount,
					View:      v.View,
					Title:     l.Label,
					Latitude:  m.Latitude,
					Longitude: m.Longitude,
					Datum: func() string {
						if m.Datum == "WGS84" {
							return "EPSG:4326"
						}
						return m.Datum
					}(),
					Height: func() float64 {
						if m.Elevation == 9999 {
							return 0
						}
						return m.Elevation
					}(),
					Azimuth:   v.Azimuth,
					Ground:    v.Offset.Vertical,
					Volcanoes: volcanoes,
				})
			}
		}
	}

	switch {
	case output != "":
		if err := Volcams(volcs).EncodeFile(output); err != nil {
			log.Fatal(err)
		}
	default:
		if err := Volcams(volcs).Encode(os.Stdout); err != nil {
			log.Fatal(err)
		}
	}

}
