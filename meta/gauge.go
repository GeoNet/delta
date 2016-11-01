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
	gaugeTimeZone
	gaugeLatitude
	gaugeLongitude
	gaugeLast
)

type Gauge struct {
	Reference
	Point

	Number   string
	TimeZone float64
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
	}}
	for _, v := range g {
		data = append(data, []string{
			strings.TrimSpace(v.Code),
			strings.TrimSpace(v.Network),
			strings.TrimSpace(v.Number),
			strconv.FormatFloat(v.TimeZone, 'g', -1, 64),
			strconv.FormatFloat(v.Latitude, 'g', -1, 64),
			strconv.FormatFloat(v.Longitude, 'g', -1, 64),
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
			if zone, err = strconv.ParseFloat(d[gaugeTimeZone], 64); err != nil {
				return err
			}
			if lat, err = strconv.ParseFloat(d[gaugeLatitude], 64); err != nil {
				return err
			}
			if lon, err = strconv.ParseFloat(d[gaugeLongitude], 64); err != nil {
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
				},
				TimeZone: zone,
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
