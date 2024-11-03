package meta

import (
	"sort"
	"strings"
	"time"
)

const (
	timingStation = iota
	timingLocation
	timingCorrection
	timingStart
	timingEnd
	timingLast
)

var timingHeaders Header = map[string]int{
	"Station":    timingStation,
	"Location":   timingLocation,
	"Correction": timingCorrection,
	"Start Date": timingStart,
	"End Date":   timingEnd,
}

var TimingTable Table = Table{
	name:    "Timing",
	headers: timingHeaders,
	primary: []string{"Station", "Location", "Start Date"},
	native:  []string{},
	foreign: map[string][]string{
		"Site": {"Station", "Location"},
	},
	remap: map[string]string{
		"Start Date": "Start",
		"End Date":   "End",
	},
	start: "Start Date",
	end:   "End Date",
}

type Timing struct {
	Span

	Station    string
	Location   string
	Correction time.Duration

	correction string
}

func (t Timing) Less(timing Timing) bool {
	switch {
	case t.Station < timing.Station:
		return true
	case t.Station > timing.Station:
		return false
	case t.Location < timing.Location:
		return true
	case t.Location > timing.Location:
		return false
	case t.Start.Before(timing.Start):
		return true
	default:
		return false
	}
}

type TimingList []Timing

func (t TimingList) Len() int           { return len(t) }
func (t TimingList) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t TimingList) Less(i, j int) bool { return t[i].Less(t[j]) }

func (t TimingList) encode() [][]string {
	var data [][]string

	data = append(data, timingHeaders.Columns())

	for _, row := range t {
		correction := strings.TrimSpace(row.correction)

		data = append(data, []string{
			strings.TrimSpace(row.Station),
			strings.TrimSpace(row.Location),
			correction,
			row.Start.Format(DateTimeFormat),
			row.End.Format(DateTimeFormat),
		})
	}

	return data
}

func (t *TimingList) decode(data [][]string) error {
	if !(len(data) > 1) {
		return nil
	}

	var timings []Timing

	fields := timingHeaders.Fields(data[0])
	for _, row := range data[1:] {
		// reorder the entries based on field labels
		entries := fields.Remap(row)

		start, err := time.Parse(DateTimeFormat, entries[timingStart])
		if err != nil {
			return err
		}

		end, err := time.Parse(DateTimeFormat, entries[timingEnd])
		if err != nil {
			return err
		}

		var correction time.Duration
		if str := entries[timingCorrection]; str != "" {
			corr, err := time.ParseDuration(str)
			if err != nil {
				return err
			}
			correction = corr
		}

		timings = append(timings, Timing{
			Station:    strings.TrimSpace(entries[timingStation]),
			Location:   strings.TrimSpace(entries[timingLocation]),
			Correction: correction,
			Span: Span{
				Start: start,
				End:   end,
			},
			correction: strings.TrimSpace(entries[timingCorrection]),
		})
	}

	*t = TimingList(timings)

	return nil
}

func LoadTimings(path string) ([]Timing, error) {
	var timings []Timing

	if err := LoadList(path, (*TimingList)(&timings)); err != nil {
		return nil, err
	}

	sort.Sort(TimingList(timings))

	return timings, nil
}
