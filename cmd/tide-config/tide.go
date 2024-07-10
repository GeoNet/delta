package main

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
	"time"

	"github.com/GeoNet/delta/meta"
)

// Tide represents a model of the tidal parameters at a recording station
// together with the installed datalogger pairs, allowing for multiple sensors
// at the same time but different locations.
type Tide struct {
	Station      meta.Station
	Gauge        meta.Gauge
	Constituents []meta.Constituent
	Installs     map[string][]meta.Collection
}

// Parse processes a template using Tide details.
func Parse(set *meta.Set, config []byte) ([]byte, error) {

	// build the set of known tidal details
	tides := make(map[string]Tide)

	// find the tidal constituents
	constituents := set.Constituents()

	// for each gauge site
	for _, gauge := range set.Gauges() {

		// extra the constituents
		var list []meta.Constituent
		for _, c := range constituents {
			if c.Gauge != gauge.Code {
				continue
			}
			list = append(list, c)
		}

		// find the associated tidal station
		station, ok := set.Station(gauge.Code)
		if !ok {
			continue
		}

		// find all installed sensors
		installs := make(map[string][]meta.Collection)
		for _, site := range set.Sites() {
			if site.Station != station.Code {
				continue
			}
			installs[site.Location] = set.Collections(site)
		}

		// remember this tide
		tides[gauge.Code] = Tide{
			Gauge:        gauge,
			Station:      station,
			Constituents: list,
			Installs:     installs,
		}
	}

	// parse the template
	tmpl, err := template.New("config").Funcs(
		template.FuncMap{
			"add": func(a, b float64) float64 {
				return a + b
			},
			"sub": func(a, b float64) float64 {
				return a - b
			},
			"mult": func(a, b float64) float64 {
				return a * b
			},
			"div": func(a, b float64) float64 {
				if b != 0.0 {
					return a / b
				}
				return 0.0
			},
			"lower": func(str string) string {
				return strings.ToLower(str)
			},
			"upper": func(str string) string {
				return strings.ToUpper(str)
			},
			"now": func() time.Time {
				return time.Now()
			},
			"before": func(a, b time.Time) bool {
				return a.Before(b)
			},
			"after": func(a, b time.Time) bool {
				return a.After(b)
			},
			"dict": func(values ...interface{}) (map[string]interface{}, error) {
				if len(values)%2 != 0 {
					return nil, fmt.Errorf("invalid dict call")
				}
				dict := make(map[string]interface{}, len(values)/2)
				for i := 0; i < len(values); i += 2 {
					key, ok := values[i].(string)
					if !ok {
						return nil, fmt.Errorf("dict keys must be strings")
					}
					dict[key] = values[i+1]
				}
				return dict, nil
			},
			"array": func(values ...interface{}) []interface{} {
				var array []interface{}
				for i := 0; i < len(values); i++ {
					array = append(array, values[i])
				}
				return array
			},
		},
	).Parse(string(config))
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, tides); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
