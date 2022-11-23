package meta

import (
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

type Polarity struct {
	Span

	Station     string
	Location    string
	Sublocation string
	Subsource   string
	Primary     bool
	Reversed    bool
	Method      string
	Citation    string

	primary  string
	reversed string
}

func (p Polarity) Less(polarity Polarity) bool {
	switch {
	case p.Station < polarity.Station:
		return true
	case p.Station > polarity.Station:
		return false
	case p.Location < polarity.Location:
		return true
	case p.Location > polarity.Location:
		return false
	case p.Sublocation < polarity.Sublocation:
		return true
	case p.Sublocation > polarity.Sublocation:
		return false
	case p.Subsource < polarity.Subsource:
		return true
	case p.Subsource > polarity.Subsource:
		return false
	case p.Start.Before(polarity.Start):
		return true
	default:
		return false
	}
}

type PolarityList []Polarity

func (p PolarityList) Len() int           { return len(p) }
func (p PolarityList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PolarityList) Less(i, j int) bool { return p[i].Less(p[j]) }

func (p PolarityList) encode() [][]string {
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
	for _, v := range p {
		primary := strings.TrimSpace(v.primary)
		if primary == "" && v.Primary {
			primary = strconv.FormatBool(v.Primary)
		}

		reversed := strings.TrimSpace(v.reversed)
		if reversed == "" && v.Reversed {
			reversed = strconv.FormatBool(v.Reversed)
		}

		data = append(data, []string{
			strings.TrimSpace(v.Station),
			strings.TrimSpace(v.Location),
			strings.TrimSpace(v.Sublocation),
			strings.TrimSpace(v.Subsource),
			primary,
			reversed,
			strings.TrimSpace(v.Method),
			strings.TrimSpace(v.Citation),
			v.Start.Format(DateTimeFormat),
			v.End.Format(DateTimeFormat),
		})
	}
	return data
}

func (p *PolarityList) decode(data [][]string) error {
	if !(len(data) > 1) {
		return nil
	}

	var polarities []Polarity
	for _, v := range data[1:] {
		if len(v) != polarityLast {
			return fmt.Errorf("incorrect number of installed polarity fields")
		}

		start, err := time.Parse(DateTimeFormat, v[polarityStart])
		if err != nil {
			return err
		}

		end, err := time.Parse(DateTimeFormat, v[polarityEnd])
		if err != nil {
			return err
		}

		var primary bool
		if s := v[polarityPrimary]; s != "" {
			v, err := strconv.ParseBool(s)
			if err != nil {
				return err
			}
			primary = v
		}

		var reversed bool
		if s := v[polarityReversed]; s != "" {
			v, err := strconv.ParseBool(s)
			if err != nil {
				return err
			}
			reversed = v
		}

		polarities = append(polarities, Polarity{
			Station:     strings.TrimSpace(v[polarityStation]),
			Location:    strings.TrimSpace(v[polarityLocation]),
			Sublocation: strings.TrimSpace(v[polaritySublocation]),
			Subsource:   strings.TrimSpace(v[polaritySubsource]),
			Reversed:    reversed,
			Primary:     primary,
			Method:      strings.TrimSpace(v[polarityMethod]),
			Citation:    strings.TrimSpace(v[polarityCitation]),
			Span: Span{
				Start: start,
				End:   end,
			},
			primary:  strings.TrimSpace(v[polarityPrimary]),
			reversed: strings.TrimSpace(v[polarityReversed]),
		})
	}

	*p = PolarityList(polarities)

	return nil
}

func LoadPolarities(path string) ([]Polarity, error) {
	var s []Polarity

	if err := LoadList(path, (*PolarityList)(&s)); err != nil {
		return nil, err
	}

	sort.Sort(PolarityList(s))

	return s, nil
}
