package meta

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	polarityStation = iota
	polarityLocation
	polaritySublocation
	polaritySubsource
	polarityPrimary
	polarityReversed
	polarityMethod
	polarityCitation
	polarityStart
	polarityEnd
	polarityLast
)

var ErrInvalidPolarityMethod = errors.New("invalid polarity method")

const (
	PolarityUnknown PolarityMethod = iota
	PolarityCompass
)

// PolarityMethod lists the various mechanisms for determining signal polarity.
type PolarityMethod int

// String implements the Stringer interface.
func (p PolarityMethod) String() string {
	switch p {
	case PolarityUnknown:
		return "unknown"
	case PolarityCompass:
		return "compass"
	default:
		return "unknown"
	}
}

// UnmarshalText implements the TextUnmarshaler interface.
func (p *PolarityMethod) UnmarshalText(data []byte) error {
	switch strings.ToLower(string(data)) {
	case "unknown":
		*p = PolarityUnknown
	case "compass":
		*p = PolarityCompass
	default:
		return ErrInvalidPolarityMethod
	}
	return nil
}

// MarshalText implements the TextMarshaler interface.
func (p PolarityMethod) MarshalText() ([]byte, error) {
	return []byte(p.String()), nil
}

// Polarity defines times where the sensor or datalogger installation results in a signal may be opposite polarity to that intended.
type Polarity struct {
	Span

	Station     string
	Location    string
	Sublocation string
	Subsource   string
	Primary     bool
	Reversed    bool
	Method      PolarityMethod
	Citation    string

	method   string
	primary  string
	reversed string
}

// Id returns a unique string which can be used for sorting or checking.
func (g Polarity) Id() string {
	return strings.Join([]string{g.Station, g.Location, g.Sublocation, g.Subsource}, ":")
}

// Less returns whether one Polarity sorts before another.
func (g Polarity) Less(polarity Polarity) bool {
	switch {
	case g.Station < polarity.Station:
		return true
	case g.Station > polarity.Station:
		return false
	case g.Location < polarity.Location:
		return true
	case g.Location > polarity.Location:
		return false
	case g.Sublocation < polarity.Sublocation:
		return true
	case g.Sublocation > polarity.Sublocation:
		return false
	case g.Subsource < polarity.Subsource:
		return true
	case g.Subsource > polarity.Subsource:
		return false
	case g.Span.Start.Before(polarity.Span.Start):
		return true
	default:
		return false
	}
}

// Components returns a sorted slice of single defined components.
func (g Polarity) Components() []string {
	switch {
	case g.Subsource == "":
		return []string{""}
	default:
		var comps []string
		for _, c := range g.Subsource {
			comps = append(comps, string(c))
		}
		return comps
	}
}

// Polarities returns a sorted slice of single Polarity entries.
func (g Polarity) Polarities() []Polarity {
	var polarities []Polarity
	for _, c := range g.Components() {
		polarities = append(polarities, Polarity{
			Span:        g.Span,
			Reversed:    g.Reversed,
			Station:     g.Station,
			Location:    g.Location,
			Sublocation: g.Sublocation,
			Subsource:   c,
		})
	}

	sort.Slice(polarities, func(i, j int) bool { return polarities[i].Less(polarities[j]) })

	return polarities
}

type PolarityList []Polarity

func (g PolarityList) Len() int           { return len(g) }
func (g PolarityList) Swap(i, j int)      { g[i], g[j] = g[j], g[i] }
func (g PolarityList) Less(i, j int) bool { return g[i].Less(g[j]) }

func (g PolarityList) encode() [][]string {
	data := [][]string{{
		"Station",
		"Location",
		"Sublocation",
		"Subsource",
		"Primary",
		"Reversed",
		"Method",
		"Citation",
		"Start Date",
		"End Date",
	}}

	for _, v := range g {
		data = append(data, []string{
			strings.TrimSpace(v.Station),
			strings.TrimSpace(v.Location),
			strings.TrimSpace(v.Sublocation),
			strings.TrimSpace(v.Subsource),
			strings.TrimSpace(v.primary),
			strings.TrimSpace(v.reversed),
			strings.TrimSpace(v.method),
			strings.TrimSpace(v.Citation),
			v.Start.Format(DateTimeFormat),
			v.End.Format(DateTimeFormat),
		})
	}

	return data
}

func (g *PolarityList) decode(data [][]string) error {
	var polarities []Polarity
	if len(data) > 1 {
		for _, d := range data[1:] {
			if len(d) != polarityLast {
				return fmt.Errorf("incorrect number of installed polarity fields")
			}

			var primary bool
			if s := strings.TrimSpace(d[polarityPrimary]); s != "" {
				b, err := strconv.ParseBool(s)
				if err != nil {
					return err
				}
				primary = b
			}

			var reversed bool
			if s := strings.TrimSpace(d[polarityReversed]); s != "" {
				b, err := strconv.ParseBool(s)
				if err != nil {
					return err
				}
				reversed = b
			}

			var method PolarityMethod
			if s := strings.TrimSpace(d[polarityMethod]); s != "" {
				if err := method.UnmarshalText([]byte(s)); err != nil {
					return err
				}
			}

			start, err := time.Parse(DateTimeFormat, d[polarityStart])
			if err != nil {
				return err
			}

			end, err := time.Parse(DateTimeFormat, d[polarityEnd])
			if err != nil {
				return err
			}

			polarities = append(polarities, Polarity{
				Span: Span{
					Start: start,
					End:   end,
				},
				Station:     strings.TrimSpace(d[polarityStation]),
				Location:    strings.TrimSpace(d[polarityLocation]),
				Sublocation: strings.TrimSpace(d[polaritySublocation]),
				Subsource:   strings.TrimSpace(d[polaritySubsource]),
				Primary:     primary,
				Reversed:    reversed,
				Method:      method,
				Citation:    strings.TrimSpace(d[polarityCitation]),

				method:   strings.TrimSpace(d[polarityMethod]),
				primary:  strings.TrimSpace(d[polarityPrimary]),
				reversed: strings.TrimSpace(d[polarityReversed]),
			})
		}

		*g = PolarityList(polarities)
	}

	return nil
}

func LoadPolarities(path string) ([]Polarity, error) {
	var g []Polarity

	if err := LoadList(path, (*PolarityList)(&g)); err != nil {
		return nil, err
	}

	sort.Sort(PolarityList(g))

	return g, nil
}
