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
	for _, view := range set.Views() {
		if time.Since(view.Span.End) > 0 {
			continue
		}

		for _, mount := range set.Mounts() {
			if time.Since(mount.Span.End) > 0 {
				continue
			}
			if mount.Code != view.Mount {
				continue
			}

			for _, camera := range set.InstalledCameras() {
				if time.Since(camera.Span.End) > 0 {
					continue
				}

				if camera.Mount != mount.Code {
					continue
				}

				if camera.View != view.Code {
					continue
				}

				targets := make(map[string]string)
				for id, vals := range keywords {
					target, ok := titles[id]
					if !ok {
						continue
					}
					for _, val := range vals {
						if strings.Contains(view.Label, val) {
							targets[id] = target
						}
						if strings.Contains(view.Description, val) {
							targets[id] = target
						}
					}
				}

				var volcanoes []Volcano
				for id, title := range targets {
					volcanoes = append(volcanoes, Volcano{
						Id:    id,
						Title: title,
					})
				}

				volcs = append(volcs, Volcam{
					Id:        strings.Join(strings.FieldsFunc(strings.ToLower(view.Label), letters), ""),
					Mount:     camera.Mount,
					View:      camera.View,
					Title:     view.Label,
					Latitude:  mount.Latitude,
					Longitude: mount.Longitude,
					Datum: func() string {
						if mount.Datum == "WGS84" {
							return "EPSG:4326"
						}
						return mount.Datum
					}(),
					Height: func() float64 {
						if mount.Elevation == 9999 {
							return 0
						}
						return mount.Elevation
					}(),
					Azimuth:   camera.Azimuth,
					Ground:    camera.Offset.Vertical,
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
