package meta

import (
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	gaugeCode = iota
	gaugeNetwork
	gaugeNumber
	gaugeAnalysisTimeZone
	gaugeAnalysisLatitude
	gaugeAnalysisLongitude
	gaugeCrex
	gaugeStart
	gaugeEnd
	gaugeLast
)

var gaugeHeaders Header = map[string]int{
	"Gauge":                 gaugeCode,
	"Network":               gaugeNetwork,
	"Identification Number": gaugeNumber,
	"Analysis Time Zone":    gaugeAnalysisTimeZone,
	"Analysis Latitude":     gaugeAnalysisLatitude,
	"Analysis Longitude":    gaugeAnalysisLongitude,
	"Crex Tag":              gaugeCrex,
	"Start Date":            gaugeStart,
	"End Date":              gaugeEnd,
}

type Gauge struct {
	Span
	Reference
	Position

	Number   string
	TimeZone float64
	Crex     string

	timeZone string // shadow variable to maintain formatting
}

type GaugeList []Gauge

func (g GaugeList) Len() int      { return len(g) }
func (g GaugeList) Swap(i, j int) { g[i], g[j] = g[j], g[i] }
func (g GaugeList) Less(i, j int) bool {
	switch {
	case g[i].Code < g[j].Code:
		return true
	case g[i].Code > g[j].Code:
		return false
	case g[i].Start.Before(g[j].Start):
		return true
	default:
		return false
	}
}

func (g GaugeList) encode() [][]string {
	var data [][]string

	data = append(data, gaugeHeaders.Columns())

	for _, row := range g {
		data = append(data, []string{
			strings.TrimSpace(row.Code),
			strings.TrimSpace(row.Network),
			strings.TrimSpace(row.Number),
			strings.TrimSpace(row.timeZone),
			strings.TrimSpace(row.latitude),
			strings.TrimSpace(row.longitude),
			strings.TrimSpace(row.Crex),
			row.Start.Format(DateTimeFormat),
			row.End.Format(DateTimeFormat),
		})
	}

	return data
}

func (g *GaugeList) decode(data [][]string) error {
	if !(len(data) > 1) {
		return nil
	}

	var gauges []Gauge

	fields := gaugeHeaders.Fields(data[0])
	for _, row := range data[1:] {
		d := fields.Remap(row)

		zone, err := strconv.ParseFloat(d[gaugeAnalysisTimeZone], 64)
		if err != nil {
			return err
		}
		lat, err := strconv.ParseFloat(d[gaugeAnalysisLatitude], 64)
		if err != nil {
			return err
		}
		lon, err := strconv.ParseFloat(d[gaugeAnalysisLongitude], 64)
		if err != nil {
			return err
		}

		start, err := time.Parse(DateTimeFormat, d[gaugeStart])
		if err != nil {
			return err
		}
		end, err := time.Parse(DateTimeFormat, d[gaugeEnd])
		if err != nil {
			return err
		}

		gauges = append(gauges, Gauge{
			Span: Span{
				Start: start,
				End:   end,
			},
			Reference: Reference{
				Code:    strings.TrimSpace(d[gaugeCode]),
				Network: strings.TrimSpace(d[gaugeNetwork]),
			},
			Number: strings.TrimSpace(d[gaugeNumber]),
			Position: Position{
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
