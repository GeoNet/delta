package meta

import (
	"sort"
	"strings"
	"time"

	"github.com/GeoNet/delta/internal/expr"
)

const (
	preampStation = iota
	preampLocation
	preampSubsource
	preampScaleFactor
	preampStart
	preampEnd
	preampLast
)

var preampHeaders Header = map[string]int{
	"Station":      preampStation,
	"Location":     preampLocation,
	"Subsource":    preampSubsource,
	"Scale Factor": preampScaleFactor,
	"Start Date":   preampStart,
	"End Date":     preampEnd,
}

var PreampTable Table = Table{
	name:    "Preamp",
	headers: preampHeaders,
	primary: []string{"Station", "Location", "Subsource", "Start Date"},
	native:  []string{"Scale Factor"},
	foreign: map[string]map[string]string{
		"Site": {"Station": "Station", "Location": "Location"},
	},
	nullable: []string{"Subsource"},
	remap: map[string]string{
		"Scale Factor": "ScaleFactor",
		"Start Date":   "Start",
		"End Date":     "End",
	},
	start: "Start Date",
	end:   "End Date",
}

// Preamp describes when a datalogger is using an analogue pre-amplification gain setting to boost the input signal.
type Preamp struct {
	Span

	Station     string
	Location    string
	Subsource   string
	ScaleFactor float64

	factor string
}

// Id returns a unique string which can be used for sorting or checking.
func (p Preamp) Id() string {
	return strings.Join([]string{p.Station, p.Location, p.Subsource}, ":")
}

// Less returns whether one Preamp sorts before another.
func (p Preamp) Less(preamp Preamp) bool {
	switch {
	case p.Station < preamp.Station:
		return true
	case p.Station > preamp.Station:
		return false
	case p.Location < preamp.Location:
		return true
	case p.Location > preamp.Location:
		return false
	case p.Subsource < preamp.Subsource:
		return true
	case p.Subsource > preamp.Subsource:
		return false
	case p.Span.Start.Before(preamp.Span.Start):
		return true
	default:
		return false
	}
}

// Subsources returns a sorted slice of single byte defined components which allows unpacking multiple subsources.
func (p Preamp) Subsources() []string {
	var comps []string
	for _, c := range p.Subsource {
		comps = append(comps, string(c))
	}
	return comps
}

// Preamps returns a sorted slice of single Preamp entries by unpacking multiple subsources if present.
func (p Preamp) Preamps() []Preamp {
	var preamps []Preamp
	for _, c := range p.Subsources() {
		preamps = append(preamps, Preamp{
			Span:        p.Span,
			ScaleFactor: p.ScaleFactor,
			Station:     p.Station,
			Location:    p.Location,
			Subsource:   string(c),
		})
	}

	sort.Slice(preamps, func(i, j int) bool {
		return preamps[i].Less(preamps[j])
	})

	return preamps
}

type PreampList []Preamp

func (p PreampList) Len() int           { return len(p) }
func (p PreampList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PreampList) Less(i, j int) bool { return p[i].Less(p[j]) }

func (p PreampList) encode() [][]string {

	var data [][]string

	data = append(data, preampHeaders.Columns())

	for _, row := range p {
		data = append(data, []string{
			strings.TrimSpace(row.Station),
			strings.TrimSpace(row.Location),
			strings.TrimSpace(row.Subsource),
			row.factor,
			row.Start.Format(DateTimeFormat),
			row.End.Format(DateTimeFormat),
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

func (p *PreampList) decode(data [][]string) error {
	if !(len(data) > 1) {
		return nil
	}

	var preamps []Preamp

	fields := preampHeaders.Fields(data[0])
	for _, row := range data[1:] {
		d := fields.Remap(row)

		factor, err := p.toFloat64(d[preampScaleFactor], 1.0)
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
			ScaleFactor: factor,
			Station:     strings.TrimSpace(d[preampStation]),
			Location:    strings.TrimSpace(d[preampLocation]),
			Subsource:   strings.TrimSpace(d[preampSubsource]),

			factor: strings.TrimSpace(d[preampScaleFactor]),
		})
	}

	*p = PreampList(preamps)

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
