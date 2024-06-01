package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"text/template"

	delta "github.com/GeoNet/delta"
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
	set, err := delta.NewBase(base)
	if err != nil {
		log.Fatal(err)
	}

	constituents := set.Constituents()

	// build the set of known tidal
	tides := make(map[string]Tide)

	// for each gauge site
	for _, gauge := range set.Gauges() {
		var list []meta.Constituent
		for _, c := range constituents {
			if c.Gauge != gauge.Code {
				continue
			}
			list = append(list, c)
		}

		// remember this tide
		tides[gauge.Code] = Tide{
			Gauge:        gauge,
			Constituents: list,
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
