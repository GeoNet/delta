package meta

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/GeoNet/delta/internal/expr"
)

const (
	preampStation = iota
	preampLocation
	preampSubsource
	preampGain
	preampStart
	preampEnd
	preampLast
)

// Preamp describes when a datalogger is using an analogue pre-amplification gain setting to boost the input signal.
type Preamp struct {
	Span

	Station   string
	Location  string
	Subsource string
	Gain      float64

	gain string
}

// Id returns a unique string which can be used for sorting or checking.
func (g Preamp) Id() string {
	return strings.Join([]string{g.Station, g.Location, g.Subsource}, ":")
}

// Less returns whether one Preamp sorts before another.
func (g Preamp) Less(preamp Preamp) bool {
	switch {
	case g.Station < preamp.Station:
		return true
	case g.Station > preamp.Station:
		return false
	case g.Location < preamp.Location:
		return true
	case g.Location > preamp.Location:
		return false
	case g.Subsource < preamp.Subsource:
		return true
	case g.Subsource > preamp.Subsource:
		return false
	case g.Gain < preamp.Gain:
		return true
	case g.Gain > preamp.Gain:
		return false
	case g.Span.Start.Before(preamp.Span.Start):
		return true
	default:
		return false
	}
}

// Subsources returns a sorted slice of single byte defined components which allows unpacking multiple subsources.
func (g Preamp) Subsources() []string {
	var comps []string
	for _, c := range g.Subsource {
		comps = append(comps, string(c))
	}
	return comps
}

// Preamps returns a sorted slice of single Preamp entries by unpacking multiple subsources if present.
func (g Preamp) Preamps() []Preamp {
	var preamps []Preamp
	for _, c := range g.Subsources() {
		preamps = append(preamps, Preamp{
			Span:      g.Span,
			Gain:      g.Gain,
			Station:   g.Station,
			Location:  g.Location,
			Subsource: string(c),
		})
	}

	sort.Slice(preamps, func(i, j int) bool {
		return preamps[i].Less(preamps[j])
	})

	return preamps
}

type PreampList []Preamp

func (g PreampList) Len() int           { return len(g) }
func (g PreampList) Swap(i, j int)      { g[i], g[j] = g[j], g[i] }
func (g PreampList) Less(i, j int) bool { return g[i].Less(g[j]) }

func (g PreampList) encode() [][]string {
	data := [][]string{{
		"Station",
		"Location",
		"Subsource",
		"Gain",
		"Start Date",
		"End Date",
	}}

	for _, v := range g {
		data = append(data, []string{
			strings.TrimSpace(v.Station),
			strings.TrimSpace(v.Location),
			strings.TrimSpace(v.Subsource),
			v.gain,
			v.Start.Format(DateTimeFormat),
			v.End.Format(DateTimeFormat),
		})
	}

	return data
}

// toFloat64 is used in decoding to allow mathematical expressions as well as actual floating point values,
// if the string parameter is empty the default value will be returned.
func (g *PreampList) toFloat64(str string, def float64) (float64, error) {
	switch s := strings.TrimSpace(str); {
	case s != "":
		return expr.ToFloat64(s)
	default:
		return def, nil
	}
}

func (g *PreampList) decode(data [][]string) error {
	var preamps []Preamp

	// needs more than a comment line
	if !(len(data) > 1) {
		return nil
	}

	for _, d := range data[1:] {
		if len(d) != preampLast {
			return fmt.Errorf("incorrect number of installed preamp fields")
		}

		gain, err := g.toFloat64(d[preampGain], 1.0)
		if err != nil {
			return err
		}

		start, err := time.Parse(DateTimeFormat, d[preampStart])
		if err != nil {
			return err
		}

		end, err := time.Parse(DateTimeFormat, d[preampEnd])
		if err != nil {
			return err
		}

		preamps = append(preamps, Preamp{
			Span: Span{
				Start: start,
				End:   end,
			},
			Gain:      gain,
			Station:   strings.TrimSpace(d[preampStation]),
			Location:  strings.TrimSpace(d[preampLocation]),
			Subsource: strings.TrimSpace(d[preampSubsource]),

			gain: strings.TrimSpace(d[preampGain]),
		})
	}

	*g = PreampList(preamps)

	return nil
}

func LoadPreamps(path string) ([]Preamp, error) {
	var g []Preamp

	if err := LoadList(path, (*PreampList)(&g)); err != nil {
		return nil, err
	}

	sort.Sort(PreampList(g))

	return g, nil
}
