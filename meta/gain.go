package meta

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/GeoNet/delta/internal/expr"
)

const (
	gainStation = iota
	gainLocation
	gainSubLocation
	gainChannel
	gainScaleFactor
	gainScaleBias
	gainStart
	gainEnd
	gainLast
)

// Gain defines times where sensor installation scaling or offsets are needed, these will be applied to the
// existing values, i.e. A + BX => A + A' + (B * B') X, where A' and B' are the given bias and scaling factors.
type Gain struct {
	Span
	Scale

	Station     string
	Location    string
	SubLocation string
	Channel     string
}

// Id returns a unique string which can be used for sorting or checking.
func (g Gain) Id() string {
	return strings.Join([]string{g.Station, g.Location, g.Channel}, ":")
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
	case g.SubLocation < gain.SubLocation:
		return true
	case g.SubLocation > gain.SubLocation:
		return false
	case g.Channel < gain.Channel:
		return true
	case g.Channel > gain.Channel:
		return false
	case g.Span.Start.Before(gain.Span.Start):
		return true
	default:
		return false
	}
}

// Channels returns a sorted slice of single defined components.
func (g Gain) Channels() []string {
	var comps []string
	for _, c := range g.Channel {
		comps = append(comps, string(c))
	}
	return comps
}

// Gains returns a sorted slice of single Gain entries.
func (g Gain) Gains() []Gain {
	var gains []Gain
	for _, c := range g.Channel {
		gains = append(gains, Gain{
			Span:        g.Span,
			Scale:       g.Scale,
			Station:     g.Station,
			Location:    g.Location,
			SubLocation: g.SubLocation,
			Channel:     string(c),
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
	data := [][]string{{
		"Station",
		"Location",
		"SubLocation",
		"Channel",
		"Scale Factor",
		"Scale Bias",
		"Start Date",
		"End Date",
	}}

	for _, v := range g {
		data = append(data, []string{
			strings.TrimSpace(v.Station),
			strings.TrimSpace(v.Location),
			strings.TrimSpace(v.SubLocation),
			strings.TrimSpace(v.Channel),
			strings.TrimSpace(v.factor),
			strings.TrimSpace(v.bias),
			v.Start.Format(DateTimeFormat),
			v.End.Format(DateTimeFormat),
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
	var gains []Gain
	if len(data) > 1 {
		for _, d := range data[1:] {
			if len(d) != gainLast {
				return fmt.Errorf("incorrect number of installed gain fields")
			}

			factor, err := g.toFloat64(d[gainScaleFactor], 1.0)
			if err != nil {
				return err
			}

			bias, err := g.toFloat64(d[gainScaleBias], 0.0)
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
				Station:     strings.TrimSpace(d[gainStation]),
				Location:    strings.TrimSpace(d[gainLocation]),
				SubLocation: strings.TrimSpace(d[gainSubLocation]),
				Channel:     strings.TrimSpace(d[gainChannel]),
			})
		}

		*g = GainList(gains)
	}

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
