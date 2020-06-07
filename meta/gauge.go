package meta

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

const (
	gaugeCode = iota
	gaugeNetwork
	gaugeNumber
	gaugeAnalysisTimeZone
	gaugeAnalysisLatitude
	gaugeAnalysisLongitude
	gaugeCrex
	gaugeLast
)

type Gauge struct {
	Reference
	Point

	Number   string
	TimeZone float64
	Crex     string

	timeZone string // shadow variable to maintain formatting
}

type GaugeList []Gauge

func (g GaugeList) Len() int           { return len(g) }
func (g GaugeList) Swap(i, j int)      { g[i], g[j] = g[j], g[i] }
func (g GaugeList) Less(i, j int) bool { return g[i].Code < g[j].Code }

func (g GaugeList) encode() [][]string {
	data := [][]string{{
		"Gauge",
		"Network",
		"LINZ Number",
		"Analysis Time Zone",
		"Analysis Latitude",
		"Analysis Longitude",
		"Crex Tag",
	}}
	for _, v := range g {
		data = append(data, []string{
			strings.TrimSpace(v.Code),
			strings.TrimSpace(v.Network),
			strings.TrimSpace(v.Number),
			strings.TrimSpace(v.timeZone),
			strings.TrimSpace(v.latitude),
			strings.TrimSpace(v.longitude),
			strings.TrimSpace(v.Crex),
		})
	}
	return data
}

func (g *GaugeList) decode(data [][]string) error {
	var gauges []Gauge
	if len(data) > 1 {
		for _, d := range data[1:] {
			if len(d) != gaugeLast {
				return fmt.Errorf("incorrect number of installed gauge fields")
			}
			var err error

			var lat, lon, zone float64
			if zone, err = strconv.ParseFloat(d[gaugeAnalysisTimeZone], 64); err != nil {
				return err
			}
			if lat, err = strconv.ParseFloat(d[gaugeAnalysisLatitude], 64); err != nil {
				return err
			}
			if lon, err = strconv.ParseFloat(d[gaugeAnalysisLongitude], 64); err != nil {
				return err
			}

			gauges = append(gauges, Gauge{
				Reference: Reference{
					Code:    strings.TrimSpace(d[gaugeCode]),
					Network: strings.TrimSpace(d[gaugeNetwork]),
				},
				Number: strings.TrimSpace(d[gaugeNumber]),
				Point: Point{
					Latitude:  lat,
					Longitude: lon,
					latitude:  strings.TrimSpace(d[gaugeAnalysisLatitude]),
					longitude: strings.TrimSpace(d[gaugeAnalysisLongitude]),
				},
				Crex:     strings.TrimSpace(d[gaugeCrex]),
				TimeZone: zone,
				timeZone: strings.TrimSpace(d[gaugeAnalysisTimeZone]),
			})
		}

		*g = GaugeList(gauges)
	}
	return nil
}

func LoadGauges(path string) ([]Gauge, error) {
	var g []Gauge

	if err := LoadList(path, (*GaugeList)(&g)); err != nil {
		return nil, err
	}

	sort.Sort(GaugeList(g))

	return g, nil
}
