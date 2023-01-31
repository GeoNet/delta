package meta

import (
	"sort"
	"strings"
	"time"

	"github.com/GeoNet/delta/internal/expr"
)

const (
	telemetryStation = iota
	telemetryLocation
	telemetryScaleFactor
	telemetryStart
	telemetryEnd
	telemetryLast
)

var telemetryHeaders Header = map[string]int{
	"Station":      telemetryStation,
	"Location":     telemetryLocation,
	"Scale Factor": telemetryScaleFactor,
	"Start Date":   telemetryStart,
	"End Date":     telemetryEnd,
}

// Telemetry describes when a datalogger is connected to a sensor via analogue telemetry (e.g. FM radio).
type Telemetry struct {
	Span

	Station     string
	Location    string
	ScaleFactor float64

	factor string
}

// String implements the Stringer interface.
func (t Telemetry) String() string {
	return strings.Join([]string{t.Station, t.Location, Format(t.Start)}, " ")
}

// Id returns a unique string which can be used for sorting or checking.
func (t Telemetry) Id() string {
	return strings.Join([]string{t.Station, t.Location}, ":")
}

// Less returns whether one Telemetry sorts before another.
func (t Telemetry) Less(telemetry Telemetry) bool {
	switch {
	case t.Station < telemetry.Station:
		return true
	case t.Station > telemetry.Station:
		return false
	case t.Location < telemetry.Location:
		return true
	case t.Location > telemetry.Location:
		return false
	case t.Span.Start.Before(telemetry.Span.Start):
		return true
	default:
		return false
	}
}

type TelemetryList []Telemetry

func (t TelemetryList) Len() int           { return len(t) }
func (t TelemetryList) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t TelemetryList) Less(i, j int) bool { return t[i].Less(t[j]) }

func (t TelemetryList) encode() [][]string {
	var data [][]string

	data = append(data, telemetryHeaders.Columns())
	for _, v := range t {
		data = append(data, []string{
			strings.TrimSpace(v.Station),
			strings.TrimSpace(v.Location),
			strings.TrimSpace(v.factor),
			v.Start.Format(DateTimeFormat),
			v.End.Format(DateTimeFormat),
		})
	}

	return data
}

// toFloat64 is used in decoding to allow mathematical expressions as well as actual floating point values,
// if the string parameter is empty the default value will be returned.
func (t *TelemetryList) toFloat64(str string, def float64) (float64, error) {
	switch s := strings.TrimSpace(str); {
	case s != "":
		return expr.ToFloat64(s)
	default:
		return def, nil
	}
}

func (t *TelemetryList) decode(data [][]string) error {
	var telemetries []Telemetry

	// needs more than a comment line
	if !(len(data) > 1) {
		return nil
	}

	fields := telemetryHeaders.Fields(data[0])
	for _, v := range data[1:] {
		d := fields.Remap(v)

		factor, err := t.toFloat64(d[telemetryScaleFactor], 1.0)
		if err != nil {
			return err
		}

		start, err := time.Parse(DateTimeFormat, d[telemetryStart])
		if err != nil {
			return err
		}

		end, err := time.Parse(DateTimeFormat, d[telemetryEnd])
		if err != nil {
			return err
		}

		telemetries = append(telemetries, Telemetry{
			Span: Span{
				Start: start,
				End:   end,
			},
			ScaleFactor: factor,
			Station:     strings.TrimSpace(d[telemetryStation]),
			Location:    strings.TrimSpace(d[telemetryLocation]),

			factor: strings.TrimSpace(d[telemetryScaleFactor]),
		})
	}

	*t = TelemetryList(telemetries)

	return nil
}

func LoadTelemetries(path string) ([]Telemetry, error) {
	var g []Telemetry

	if err := LoadList(path, (*TelemetryList)(&g)); err != nil {
		return nil, err
	}

	sort.Sort(TelemetryList(g))

	return g, nil
}
