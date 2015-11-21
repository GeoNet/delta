package meta

import (
	"sort"
)

type Stream struct {
	Response      string  `csv:"Response Name"`
	Types         string  `csv:"Stream Types"`
	Prefix        string  `csv:"Stream Prefix"`
	Channels      string  `csv:"Stream Channels"`
	SampleRate    float64 `csv:"Sample Rate"`
	Frequency     float64 `csv:"Response Frequency"`
	StorageFormat string  `csv:"Storage Format"`
	ClockDrift    float64 `csv:"Clock Drift"`
	Lookup        string  `csv:"Response Lookup"`
}

type Streams []Stream

func (s Streams) Len() int      { return len(s) }
func (s Streams) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s Streams) Less(i, j int) bool {
	switch {
	case s[i].Response < s[j].Response:
		return true
	case s[i].Response > s[j].Response:
		return false
	case s[i].SampleRate < s[j].SampleRate:
		return true
	case s[i].SampleRate > s[j].SampleRate:
		return false
	case s[i].Channels < s[j].Channels:
		return true
	case s[i].Channels > s[j].Channels:
		return false
	default:
		return false
	}
}

func (s Streams) List()      {}
func (s Streams) Sort() List { sort.Sort(s); return s }
