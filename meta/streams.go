package meta

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	streamStation = iota
	streamLocation
	streamSamplingRate
	streamAxial
	streamReversed
	streamTriggered
	streamStart
	streamEnd
	streamLast
)

type Stream struct {
	Span

	Station      string
	Location     string
	SamplingRate float64
	Axial        bool
	Reversed     bool
	Triggered    bool
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
	data := [][]string{{
		"Station",
		"Location",
		"Sampling Rate",
		"Axial",
		"Reversed",
		"Triggered",
		"Start Date",
		"End Date",
	}}
	for _, v := range s {
		data = append(data, []string{
			strings.TrimSpace(v.Station),
			strings.TrimSpace(v.Location),
			strings.TrimSpace(strconv.FormatFloat(v.SamplingRate, 'g', -1, 64)),
			strings.TrimSpace(strconv.FormatBool(v.Axial)),
			strings.TrimSpace(strconv.FormatBool(v.Reversed)),
			strings.TrimSpace(strconv.FormatBool(v.Triggered)),
			v.Start.Format(DateTimeFormat),
			v.End.Format(DateTimeFormat),
		})
	}
	return data
}

func (c *StreamList) decode(data [][]string) error {
	var streams []Stream
	if len(data) > 1 {
		for _, v := range data[1:] {
			if len(v) != streamLast {
				return fmt.Errorf("incorrect number of installed stream fields")
			}
			var err error

			var start, end time.Time
			if start, err = time.Parse(DateTimeFormat, v[streamStart]); err != nil {
				return err
			}
			if end, err = time.Parse(DateTimeFormat, v[streamEnd]); err != nil {
				return err
			}

			var rate float64
			if rate, err = strconv.ParseFloat(v[streamSamplingRate], 64); err != nil {
				return err
			}

			var axial, reversed, triggered bool
			if axial, err = strconv.ParseBool(v[streamAxial]); err != nil {
				return err
			}
			if reversed, err = strconv.ParseBool(v[streamReversed]); err != nil {
				return err
			}
			if triggered, err = strconv.ParseBool(v[streamTriggered]); err != nil {
				return err
			}

			streams = append(streams, Stream{
				Station:      strings.TrimSpace(v[streamStation]),
				Location:     strings.TrimSpace(v[streamLocation]),
				SamplingRate: rate,
				Axial:        axial,
				Reversed:     reversed,
				Triggered:    triggered,
				Span: Span{
					Start: start,
					End:   end,
				},
			})
		}

		*c = StreamList(streams)
	}
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
