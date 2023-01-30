package meta

import (
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	qualityStation = iota
	qualityLocation
	qualityFault
	qualityStart
	qualityEnd
	qualityLast
)

var qualityHeaders Header = map[string]int{
	"Station":    qualityStation,
	"Location":   qualityLocation,
	"Fault":      qualityFault,
	"Start Date": qualityStart,
	"End Date":   qualityEnd,
}

// Quality describes when a datalogger is connected to a sensor via analogue quality (e.g. FM radio).
type Quality struct {
	Span

	Station  string
	Location string
	Fault    bool

	fault string
}

// String implements the Stringer interface.
func (q Quality) String() string {
	return strings.Join([]string{q.Station, q.Location, Format(q.Start)}, " ")
}

// Id returns a unique string which can be used for sorting or checking.
func (q Quality) Id() string {
	return strings.Join([]string{q.Station, q.Location}, ":")
}

// Less returns whether one Quality sorts before another.
func (q Quality) Less(quality Quality) bool {
	switch {
	case q.Station < quality.Station:
		return true
	case q.Station > quality.Station:
		return false
	case q.Location < quality.Location:
		return true
	case q.Location > quality.Location:
		return false
	case q.Span.Start.Before(quality.Span.Start):
		return true
	default:
		return false
	}
}

type QualityList []Quality

func (t QualityList) Len() int           { return len(t) }
func (t QualityList) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t QualityList) Less(i, j int) bool { return t[i].Less(t[j]) }

func (t QualityList) encode() [][]string {
	var data [][]string

	data = append(data, qualityHeaders.Columns())
	for _, v := range t {
		data = append(data, []string{
			strings.TrimSpace(v.Station),
			strings.TrimSpace(v.Location),
			strings.TrimSpace(v.fault),
			v.Start.Format(DateTimeFormat),
			v.End.Format(DateTimeFormat),
		})
	}

	return data
}

// toBool is used when decoding boolean strings to allow for a default value when the input string is empty.
func (q *QualityList) toBool(str string, def bool) (bool, error) {
	switch s := strings.TrimSpace(str); {
	case s != "":
		return strconv.ParseBool(s)
	default:
		return def, nil
	}
}

func (q *QualityList) decode(data [][]string) error {
	var telemetries []Quality

	// needs more than a comment line
	if !(len(data) > 1) {
		return nil
	}

	fields := qualityHeaders.Fields(data[0])
	for _, v := range data[1:] {
		d := fields.Remap(v)

		fault, err := q.toBool(d[qualityFault], false)
		if err != nil {
			return err
		}

		start, err := time.Parse(DateTimeFormat, d[qualityStart])
		if err != nil {
			return err
		}

		end, err := time.Parse(DateTimeFormat, d[qualityEnd])
		if err != nil {
			return err
		}

		telemetries = append(telemetries, Quality{
			Span: Span{
				Start: start,
				End:   end,
			},
			Fault:    fault,
			Station:  strings.TrimSpace(d[qualityStation]),
			Location: strings.TrimSpace(d[qualityLocation]),

			fault: strings.TrimSpace(d[qualityFault]),
		})
	}

	*q = QualityList(telemetries)

	return nil
}

func LoadQualities(path string) ([]Quality, error) {
	var g []Quality

	if err := LoadList(path, (*QualityList)(&g)); err != nil {
		return nil, err
	}

	sort.Sort(QualityList(g))

	return g, nil
}
