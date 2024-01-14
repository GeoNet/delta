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
		d := fields.Remap(row)

		start, err := time.Parse(DateTimeFormat, d[timingStart])
		if err != nil {
			return err
		}

		end, err := time.Parse(DateTimeFormat, d[timingEnd])
		if err != nil {
			return err
		}

		var correction time.Duration
		if s := d[timingCorrection]; s != "" {
			v, err := time.ParseDuration(s)
			if err != nil {
				return err
			}
			correction = v
		}

		timings = append(timings, Timing{
			Station:    strings.TrimSpace(d[timingStation]),
			Location:   strings.TrimSpace(d[timingLocation]),
			Correction: correction,
			Span: Span{
				Start: start,
				End:   end,
			},
			correction: strings.TrimSpace(d[timingCorrection]),
		})
	}

	*t = TimingList(timings)

	return nil
}

func LoadTimings(path string) ([]Timing, error) {
	var s []Timing

	if err := LoadList(path, (*TimingList)(&s)); err != nil {
		return nil, err
	}

	sort.Sort(TimingList(s))

	return s, nil
}
