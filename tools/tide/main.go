package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/GeoNet/delta/helper/metadb"
	"github.com/GeoNet/delta/meta"
)

type Install struct {
	Sensor     meta.InstalledSensor
	Datalogger meta.DeployedDatalogger
	Start      time.Time
	End        time.Time
}

type Tide struct {
	Station      meta.Station
	Gauge        meta.Gauge
	Constituents []meta.Constituent
	Installs     map[string][]Install
}

func main() {

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Provide a tidal templating\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  %s [options] [templates ....]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "General Options:\n")
		fmt.Fprintf(os.Stderr, "\n")
		flag.PrintDefaults()
	}

	var base string
	flag.StringVar(&base, "base", "../..", "base of delta files on disk")

	var header string
	flag.StringVar(&header, "header", "", "output header before proccessing any templates")

	var footer string
	flag.StringVar(&footer, "footer", "", "output footer after all templates have been proccessed")

	flag.Parse()

	db, err := metadb.NewMetaDB(base)
	if err != nil {
		fmt.Fprintf(os.Stderr, "problem loading meta db %s: %v\n", base, err)
		os.Exit(1)
	}

	tides := make(map[string]Tide)

	gauges, err := db.Gauges()
	if err != nil {
		fmt.Fprintf(os.Stderr, "problem loading gauges from db %s: %v\n", base, err)
		os.Exit(1)
	}

	for _, gauge := range gauges {

		station, err := db.Station(gauge.Code)
		if err != nil {
			fmt.Fprintf(os.Stderr, "problem loading stations from db %s [%s]: %v\n", base, gauge.Code, err)
			os.Exit(1)
		}

		constituents, err := db.Constituents(gauge.Code)
		if err != nil {
			fmt.Fprintf(os.Stderr, "problem loading constituents from db %s [%s]: %v\n", base, gauge.Code, err)
			os.Exit(1)
		}

		sites, err := db.Sites(gauge.Code)
		if err != nil {
			fmt.Fprintf(os.Stderr, "problem loading sites from db %s [%s]: %v\n", base, gauge.Code, err)
			os.Exit(1)
		}

		installs := make(map[string][]Install)
		for _, site := range sites {
			sensors, err := db.InstalledSensors(gauge.Code, site.Location)
			if err != nil {
				fmt.Fprintf(os.Stderr, "problem loading sensors from db %s [%s]: %v\n", base, gauge.Code, err)
				os.Exit(1)
			}

			for _, sensor := range sensors {
				connections, err := db.Connections(gauge.Code, site.Location)
				if err != nil {
					fmt.Fprintf(os.Stderr, "problem loading connections from db %s [%s]: %v\n", base, gauge.Code, err)
					os.Exit(1)
				}
				for _, connection := range connections {
					if connection.Start.After(sensor.End) {
						continue
					}
					if connection.End.Before(sensor.Start) {
						continue
					}
					dataloggers, err := db.DeployedDataloggers(connection.Place, connection.Role)
					if err != nil {
						fmt.Fprintf(os.Stderr, "problem loading connections from db %s [%s]: %v\n", base, gauge.Code, err)
						os.Exit(1)
					}
					for _, datalogger := range dataloggers {
						if connection.Start.After(datalogger.End) {
							continue
						}
						if connection.End.Before(datalogger.Start) {
							continue
						}

						if _, ok := installs[site.Location]; !ok {
							installs[site.Location] = []Install{}
						}
						installs[site.Location] = append(installs[site.Location], Install{
							Sensor:     sensor,
							Datalogger: datalogger,
							Start: func() time.Time {
								if sensor.Start.After(datalogger.Start) {
									return sensor.Start
								}
								return datalogger.Start
							}(),
							End: func() time.Time {
								if sensor.End.Before(datalogger.End) {
									return sensor.End
								}
								return datalogger.End
							}(),
						})
					}
				}
			}
		}

		tides[gauge.Code] = Tide{
			Gauge:        gauge,
			Station:      station,
			Constituents: constituents,
			Installs:     installs,
		}
	}

	if header != "" {
		fmt.Fprintln(os.Stdout, header)
	}

	for _, t := range flag.Args() {
		conf, err := ioutil.ReadFile(t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "problem reading template file %s: %v\n", t, err)
			os.Exit(1)
		}

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
		).Parse(string(conf))

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
	if footer != "" {
		fmt.Fprintln(os.Stdout, footer)
	}
}
