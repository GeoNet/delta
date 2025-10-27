package meta

import (
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	streamStation = iota
	streamLocation
	streamBand
	streamSource
	streamSamplingRate
	streamAxial
	streamTriggered
	streamStart
	streamEnd
	streamLast
)

var streamHeaders Header = map[string]int{
	"Station":       streamStation,
	"Location":      streamLocation,
	"Band":          streamBand,
	"Source":        streamSource,
	"Sampling Rate": streamSamplingRate,
	"Axial":         streamAxial,
	"Triggered":     streamTriggered,
	"Start Date":    streamStart,
	"End Date":      streamEnd,
}

var StreamTable Table = Table{
	name:    "Stream",
	headers: streamHeaders,
	primary: []string{"Station", "Location", "Source", "Sampling Rate", "Start Date"},
	native:  []string{"Sampling Rate"},
	foreign: map[string][]string{
		"Site": {"Station", "Location"},
	},
	remap: map[string]string{
		"Sampling Rate": "SamplingRate",
		"Start Date":    "Start",
		"End Date":      "End",
	},
	start: "Start Date",
	end:   "End Date",
}

type Stream struct {
	Span

	Station      string  `json:"station"`
	Location     string  `json:"location"`
	Band         string  `json:"band"`
	Source       string  `json:"source"`
	SamplingRate float64 `json:"sampling-rate"`
	Axial        string  `json:"axial,omitempty"`
	Triggered    bool    `json:"triggered,omitempty"`

	samplingRate string
}

func (s Stream) Less(stream Stream) bool {
	switch {
	case s.Station < stream.Station:
		return true
	case s.Station > stream.Station:
		return false
	case s.Location < stream.Location:
		return true
	case s.Location > stream.Location:
		return false
	case s.Source < stream.Source:
		return true
	case s.Source > stream.Source:
		return false
	case s.SamplingRate < stream.SamplingRate:
		return true
	case s.SamplingRate > stream.SamplingRate:
		return false
	case s.Start.Before(stream.Start):
		return true
	case s.Start.After(stream.Start):
		return false
	default:
		return false
	}
}

type StreamList []Stream

func (s StreamList) Len() int           { return len(s) }
func (s StreamList) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s StreamList) Less(i, j int) bool { return s[i].Less(s[j]) }

func (s StreamList) encode() [][]string {
	var data [][]string

	data = append(data, streamHeaders.Columns())

	for _, row := range s {
		data = append(data, []string{
			strings.TrimSpace(row.Station),
			strings.TrimSpace(row.Location),
			strings.TrimSpace(row.Band),
			strings.TrimSpace(row.Source),
			strings.TrimSpace(row.samplingRate),
			strings.TrimSpace(row.Axial),
			strings.TrimSpace(strconv.FormatBool(row.Triggered)),
			row.Start.Format(DateTimeFormat),
			row.End.Format(DateTimeFormat),
		})
	}

	return data
}

func (s *StreamList) decode(data [][]string) error {
	if !(len(data) > 1) {
		return nil
	}

	var streams []Stream

	fields := streamHeaders.Fields(data[0])
	for _, row := range data[1:] {
		d := fields.Remap(row)

		start, err := time.Parse(DateTimeFormat, d[streamStart])
		if err != nil {
			return err
		}
		end, err := time.Parse(DateTimeFormat, d[streamEnd])
		if err != nil {
			return err
		}

		rate, err := strconv.ParseFloat(d[streamSamplingRate], 64)
		if err != nil {
			return err
		}
		if rate < 0 {
			rate = -1.0 / rate
		}

		triggered, err := strconv.ParseBool(d[streamTriggered])
		if err != nil {
			return err
		}

		streams = append(streams, Stream{
			Station:      strings.TrimSpace(d[streamStation]),
			Location:     strings.TrimSpace(d[streamLocation]),
			Band:         strings.TrimSpace(d[streamBand]),
			Source:       strings.TrimSpace(d[streamSource]),
			SamplingRate: rate,
			samplingRate: strings.TrimSpace(d[streamSamplingRate]),
			Axial:        strings.TrimSpace(d[streamAxial]),
			Triggered:    triggered,
			Span: Span{
				Start: start,
				End:   end,
			},
		})
	}

	*s = StreamList(streams)

	return nil
}

func LoadStreams(path string) ([]Stream, error) {
	var s []Stream

	if err := LoadList(path, (*StreamList)(&s)); err != nil {
		return nil, err
	}

	sort.Sort(StreamList(s))

	return s, nil
}
