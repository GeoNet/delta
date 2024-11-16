package meta

import (
	"sort"
	"strings"
	"time"

	"github.com/GeoNet/delta/internal/expr"
)

const (
	gainStation = iota
	gainLocation
	gainSublocation
	gainSubsource
	gainScaleFactor
	gainScaleBias
	gainAbsoluteBias
	gainStart
	gainEnd
	gainLast
)

var gainHeaders Header = map[string]int{
	"Station":       gainStation,
	"Location":      gainLocation,
	"Sublocation":   gainSublocation,
	"Subsource":     gainSubsource,
	"Scale Factor":  gainScaleFactor,
	"Scale Bias":    gainScaleBias,
	"Absolute Bias": gainAbsoluteBias,
	"Start Date":    gainStart,
	"End Date":      gainEnd,
}

var GainTable Table = Table{
	name:     "Gain",
	headers:  gainHeaders,
	primary:  []string{"Station", "Location", "Sublocation", "Subsource", "Start Date"},
	native:   []string{"Scale Factor", "Scale Bias", "Absolute Bias"},
	foreign:  map[string]map[string]string{},
	nullable: []string{"Absolute Bias", "Scale Bias", "Scale Factor", "Sublocation", "Subsource"},
	remap: map[string]string{
		"Scale Factor":  "Factor",
		"Scale Bias":    "Bias",
		"Absolute Bias": "Absolute",
		"Start Date":    "Start",
		"End Date":      "End",
	},
	start: "Start Date",
	end:   "End Date",
}

// Gain defines times where sensor installation scaling or offsets are needed, these will be applied to the
// existing values, i.e. A * X + B => A * A' * X + B * A' + A * B' + C
// where A' and B' are the gain scale factor and bias and C is the absolute bias.
type Gain struct {
	Span
	Scale

	Station     string
	Location    string
	Sublocation string
	Subsource   string
	Absolute    float64

	absolute string
}

// Id returns a unique string which can be used for sorting or checking.
func (g Gain) Id() string {
	return strings.Join([]string{g.Station, g.Location, g.Subsource}, ":")
}

// Less returns whether one Gain sorts before another.
func (g Gain) Less(gain Gain) bool {
	switch {
	case g.Station < gain.Station:
		return true
	case g.Station > gain.Station:
		return false
	case g.Location < gain.Location:
		return true
	case g.Location > gain.Location:
		return false
	case g.Sublocation < gain.Sublocation:
		return true
	case g.Sublocation > gain.Sublocation:
		return false
	case g.Subsource < gain.Subsource:
		return true
	case g.Subsource > gain.Subsource:
		return false
	case g.Span.Start.Before(gain.Span.Start):
		return true
	default:
		return false
	}
}

// Subsources returns a sorted slice of single defined components.
func (g Gain) Subsources() []string {
	var comps []string
	for _, c := range g.Subsource {
		comps = append(comps, string(c))
	}
	return comps
}

// Gains returns a sorted slice of single Gain entries.
func (g Gain) Gains() []Gain {
	var gains []Gain
	for _, c := range g.Subsources() {
		gains = append(gains, Gain{
			Span:        g.Span,
			Scale:       g.Scale,
			Absolute:    g.Absolute,
			Station:     g.Station,
			Location:    g.Location,
			Sublocation: g.Sublocation,
			Subsource:   string(c),
			absolute:    g.absolute,
		})
	}

	sort.Slice(gains, func(i, j int) bool { return gains[i].Less(gains[j]) })

	return gains
}

type GainList []Gain

func (g GainList) Len() int           { return len(g) }
func (g GainList) Swap(i, j int)      { g[i], g[j] = g[j], g[i] }
func (g GainList) Less(i, j int) bool { return g[i].Less(g[j]) }

func (g GainList) encode() [][]string {

	var data [][]string

	data = append(data, gainHeaders.Columns())

	for _, row := range g {
		data = append(data, []string{
			strings.TrimSpace(row.Station),
			strings.TrimSpace(row.Location),
			strings.TrimSpace(row.Sublocation),
			strings.TrimSpace(row.Subsource),
			strings.TrimSpace(row.Scale.factor),
			strings.TrimSpace(row.Scale.bias),
			strings.TrimSpace(row.absolute),
			row.Start.Format(DateTimeFormat),
			row.End.Format(DateTimeFormat),
		})
	}

	return data
}

func (g *GainList) toFloat64(str string, def float64) (float64, error) {
	switch s := strings.TrimSpace(str); {
	case s != "":
		return expr.ToFloat64(s)
	default:
		return def, nil
	}
}

func (g *GainList) decode(data [][]string) error {
	if !(len(data) > 1) {
		return nil
	}

	var gains []Gain

	fields := gainHeaders.Fields(data[0])
	for _, row := range data[1:] {
		d := fields.Remap(row)

		factor, err := g.toFloat64(d[gainScaleFactor], 1.0)
		if err != nil {
			return err
		}

		bias, err := g.toFloat64(d[gainScaleBias], 0.0)
		if err != nil {
			return err
		}

		absolute, err := g.toFloat64(d[gainAbsoluteBias], 0.0)
		if err != nil {
			return err
		}

		start, err := time.Parse(DateTimeFormat, d[gainStart])
		if err != nil {
			return err
		}

		end, err := time.Parse(DateTimeFormat, d[gainEnd])
		if err != nil {
			return err
		}

		gains = append(gains, Gain{
			Span: Span{
				Start: start,
				End:   end,
			},
			Scale: Scale{
				Factor: factor,
				Bias:   bias,

				factor: strings.TrimSpace(d[gainScaleFactor]),
				bias:   strings.TrimSpace(d[gainScaleBias]),
			},
			Absolute:    absolute,
			Station:     strings.TrimSpace(d[gainStation]),
			Location:    strings.TrimSpace(d[gainLocation]),
			Sublocation: strings.TrimSpace(d[gainSublocation]),
			Subsource:   strings.TrimSpace(d[gainSubsource]),
			absolute:    strings.TrimSpace(d[gainAbsoluteBias]),
		})
	}

	*g = GainList(gains)

	return nil
}

func LoadGains(path string) ([]Gain, error) {
	var g []Gain

	if err := LoadList(path, (*GainList)(&g)); err != nil {
		return nil, err
	}

	sort.Sort(GainList(g))

	return g, nil
}
