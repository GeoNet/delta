package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"text/template"

	"github.com/GeoNet/delta/internal/metadb"
	"github.com/GeoNet/delta/meta"
)

const tmpl = `
package tides

/*
 *  WARNING: CODE GENERATED AUTOMATICALLY.
 *
 *  To update: run "go generate" in the tide directory and
 *  commit any changes to the generated files.
 *
 *  THIS FILE SHOULD NOT BE EDITED BY HAND.
 *
 */


var _tides = map[string]Tide{
{{ range $s := .}}
"{{$s.Gauge.Code}}": {
Code: "{{$s.Gauge.Code}}",
Network: "{{$s.Gauge.Network}}",
Number: "{{$s.Gauge.Number}}",
TimeZone: {{$s.Gauge.TimeZone}},
Latitude: {{$s.Gauge.Latitude}},
Longitude: {{$s.Gauge.Longitude}},
Crex: "{{$s.Gauge.Crex}}",
Constituents: []Constituent{
{{ range $c := $s.Constituents}}
{
Name: "{{$c.Name}}",
Amplitude: {{$c.Amplitude}},
Lag: {{$c.Lag}},
},
{{end}}
},
},
{{ end }}
}

`

// Tide represents a model of the tidal parameters at a recording station
// together with the installed datalogger pairs, allowing for multiple sensors
// at the same time but different locations.
type Tide struct {
	Gauge        meta.Gauge
	Constituents []meta.Constituent
}

func main() {

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Provide an auto generated go tide file\n")
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
	flag.StringVar(&base, "base", "..", "base of delta files on disk")

	flag.Parse()

	// load delta meta helper
	db := metadb.NewMetaDB(base)

	// build the set of known tidal
	tides := make(map[string]Tide)

	// recover linz tide gauge sites
	gauges, err := db.Gauges()
	if err != nil {
		fmt.Fprintf(os.Stderr, "problem loading gauges from db %s: %v\n", base, err)
		os.Exit(1)
	}

	// for each linz gauge site
	for _, gauge := range gauges {

		// and the associated linz tidal constituents
		constituents, err := db.GaugeConstituents(gauge.Code)
		if err != nil {
			fmt.Fprintf(os.Stderr, "problem loading constituents from db %s [%s]: %v\n", base, gauge.Code, err)
			os.Exit(1)
		}

		// remember this tide
		tides[gauge.Code] = Tide{
			Gauge: gauge,
			//Station:      *station,
			Constituents: constituents,
		}
	}

	tmpl, err := template.New("config").Parse(string(tmpl))
	if err != nil {
		fmt.Fprintf(os.Stderr, "problem compiling template: %v\n", err)
		os.Exit(1)
	}

	var res bytes.Buffer
	if err := tmpl.Execute(&res, tides); err != nil {
		fmt.Fprintf(os.Stderr, "problem parsing template: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintln(os.Stdout, res.String())
}
