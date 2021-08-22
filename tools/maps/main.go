package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	/*
		"bytes"
		"io/ioutil"
		"text/template"
	*/
	"time"

	"github.com/GeoNet/delta/internal/metadb"
	//"github.com/GeoNet/delta/meta"
)

func main() {

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Build a GeoJSON map\n")
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
	flag.StringVar(&base, "base", "../..", "base of delta files on disk")

	var dart bool
	flag.BoolVar(&dart, "dart", false, "build a DART map")

	var output string
	flag.StringVar(&output, "output", "", "output file")

	flag.Parse()

	// load delta meta helper
	db := metadb.NewMetaDB(base)

	// load station details
	stations, err := db.Stations()
	if err != nil {
		log.Fatal(err)
	}

	// load gnss mark details
	marks, err := db.Marks()
	if err != nil {
		log.Fatal(err)
	}
	fc := NewFeatureCollection()

	switch {
	case dart:
		fc.AddMetadata("name", "New Zealand DART Locations")
		for _, s := range stations {
			// needs to be a DART network station
			if s.Network != "TD" {
				continue
			}
			// needs to be operational
			if s.Start.After(time.Now()) {
				continue
			}

			lon, lat := s.Longitude, s.Latitude
			for lon < 0.0 {
				lon += 360.0
			}
			f := NewFeature()
			f.SetId(s.Code)
			f.AddPointGeometry(lon, lat)
			f.AddProperty("code", s.Code)
			f.AddProperty("name", s.Name)
			f.AddProperty("depth", fmt.Sprintf("%.0f m", s.Depth))
			f.AddProperty("opened", s.Start.Format(time.RFC3339))
			if s.End.Before(time.Now()) {
				f.AddProperty("closed", s.End.Format(time.RFC3339))
			}
			f.AddProperty("marker-size", "large")
			f.AddProperty("marker-color", "#00ff00")
			//f.AddProperty("marker-symbol", "ferry.svg")
			f.AddProperty("marker-symbol", strings.ToLower(s.Code[len(s.Code)-1:]))
			fc.AddFeature(*f)
		}
	default:
		fc.AddMetadata("name", "geonet delta locations")
		for _, s := range stations {
			network, err := db.Network(s.Network)
			if err != nil {
				log.Fatal(err)
			}
			lon, lat := s.Longitude, s.Latitude
			for lon < 0.0 {
				lon += 360.0
			}

			f := NewFeature()
			f.SetId(s.Code)
			f.AddPointGeometry(lon, lat)
			if network != nil {
				f.AddProperty("network", network.Description)
			}
			f.AddProperty("code", s.Code)
			f.AddProperty("name", s.Name)
			f.AddProperty("type", "Station")
			f.AddProperty("opened", s.Start.Format(time.RFC3339))
			if s.End.Before(time.Now()) {
				f.AddProperty("closed", s.End.Format(time.RFC3339))
			}
			fc.AddFeature(*f)
		}
		for _, m := range marks {
			network, err := db.Network(m.Network)
			if err != nil {
				log.Fatal(err)
			}
			lon, lat := m.Longitude, m.Latitude
			for lon < 0.0 {
				lon += 360.0
			}

			f := NewFeature()
			f.SetId(m.Code)
			f.AddPointGeometry(lon, lat)
			if network != nil {
				f.AddProperty("network", network.Description)
			}
			f.AddProperty("code", m.Code)
			f.AddProperty("name", m.Name)
			f.AddProperty("type", "Mark")
			f.AddProperty("opened", m.Start.Format(time.RFC3339))
			if m.End.Before(time.Now()) {
				f.AddProperty("closed", m.End.Format(time.RFC3339))
			}
			fc.AddFeature(*f)
		}
	}

	data, err := fc.MarshalIndent("", "  ")
	if err != nil {
		log.Fatal(err)
	}

	switch {
	case output != "":
		file, err := os.Create(output)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		if _, err := file.Write(data); err != nil {
			log.Fatal(err)
		}
	default:
		if _, err := os.Stdout.Write(data); err != nil {
			log.Fatal(err)
		}
	}
}
