package meta

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	sampleCode = iota
	sampleNetwork
	sampleName
	sampleLatitude
	sampleLongitude
	sampleElevation
	sampleDepth
	sampleDatum
	sampleStart
	sampleEnd
)

var sampleHeaders Header = map[string]int{
	"Station":    sampleCode,
	"Network":    sampleNetwork,
	"Name":       sampleName,
	"Latitude":   sampleLatitude,
	"Longitude":  sampleLongitude,
	"Elevation":  sampleElevation,
	"Depth":      sampleDepth,
	"Datum":      sampleDatum,
	"Start Date": sampleStart,
	"End Date":   sampleEnd,
}

var SampleTable Table = Table{
	name:    "Sample",
	headers: sampleHeaders,
	primary: []string{"Station", "Start Date"},
	native:  []string{"Latitude", "Longitude", "Elevation", "Depth"},
	foreign: map[string][]string{
		"Network": {"Network"},
	},
	remap: map[string]string{
		"Station":    "Code",
		"Start Date": "Start",
		"End Date":   "End",
	},
	start: "Start Date",
	end:   "End Date",
}

// Sample represents the location and time span of where data was manually collected.
type Sample struct {
	Reference
	Position
	Span
}

// Id is a shorthand reference to the sample for debugging or testing.
func (s Sample) Id() string {
	return fmt.Sprintf("%s_%s:%s", s.Network, s.Code, s.Start.Format(DateTimeFormat))
}

// Less allows samples to be sorted.
func (s Sample) Less(sample Sample) bool {
	switch {
	case s.Code < sample.Code:
		return true
	case s.Code > sample.Code:
		return false
	case s.Network < sample.Network:
		return true
	case s.Network > sample.Network:
		return false
	case s.Start.Before(sample.Start):
		return true
	default:
		return false
	}
}

// Overlaps allows samples to be tested.
func (s Sample) Overlaps(sample Sample) bool {
	switch {
	case s.Code != sample.Code:
		return false
	case s.Network != sample.Network:
		return false
	case !s.Span.Overlaps(sample.Span):
		return false
	default:
		return true
	}
}

// SampleList is a slice of Samples and generally maps the associated file content.
type SampleList []Sample

func (s SampleList) Len() int           { return len(s) }
func (s SampleList) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s SampleList) Less(i, j int) bool { return s[i].Less(s[j]) }

func (s SampleList) encode() [][]string {
	var data [][]string

	data = append(data, sampleHeaders.Columns())
	for _, row := range s {
		data = append(data, []string{
			strings.TrimSpace(row.Code),
			strings.TrimSpace(row.Network),
			strings.TrimSpace(row.Name),
			strings.TrimSpace(row.latitude),
			strings.TrimSpace(row.longitude),
			strings.TrimSpace(row.elevation),
			strings.TrimSpace(row.depth),
			strings.TrimSpace(row.Datum),
			row.Start.Format(DateTimeFormat),
			row.End.Format(DateTimeFormat),
		})
	}

	return data
}

func (s *SampleList) decode(data [][]string) error {
	var samples []Sample

	if !(len(data) > 1) {
		return nil
	}

	fields := sampleHeaders.Fields(data[0])
	for _, row := range data[1:] {
		d := fields.Remap(row)

		lat, err := strconv.ParseFloat(d[sampleLatitude], 64)
		if err != nil {
			return err
		}

		lon, err := strconv.ParseFloat(d[sampleLongitude], 64)
		if err != nil {
			return err
		}

		elev, err := ParseFloat64(d[sampleElevation])
		if err != nil {
			return err
		}

		depth, err := ParseFloat64(d[sampleDepth])
		if err != nil {
			return err
		}

		start, err := time.Parse(DateTimeFormat, d[sampleStart])
		if err != nil {
			return err
		}

		end, err := time.Parse(DateTimeFormat, d[sampleEnd])
		if err != nil {
			return err
		}

		samples = append(samples, Sample{
			Reference: Reference{
				Code:    strings.TrimSpace(d[sampleCode]),
				Network: strings.TrimSpace(d[sampleNetwork]),
				Name:    strings.TrimSpace(d[sampleName]),
			},
			Span: Span{
				Start: start,
				End:   end,
			},
			Position: Position{
				Latitude:  lat,
				Longitude: lon,
				Elevation: elev,
				Datum:     strings.TrimSpace(d[sampleDatum]),
				Depth:     depth,

				latitude:  strings.TrimSpace(d[sampleLatitude]),
				longitude: strings.TrimSpace(d[sampleLongitude]),
				elevation: strings.TrimSpace(d[sampleElevation]),
				depth:     strings.TrimSpace(d[sampleDepth]),
			},
		})
	}

	*s = SampleList(samples)

	return nil
}

func LoadSamples(path string) ([]Sample, error) {
	var s []Sample

	if err := LoadList(path, (*SampleList)(&s)); err != nil {
		return nil, err
	}

	sort.Sort(SampleList(s))

	return s, nil
}
