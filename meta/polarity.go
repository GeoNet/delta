package meta

import (
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

var polarityHeaders Header = map[string]int{
	"Station":     polarityStation,
	"Location":    polarityLocation,
	"Sublocation": polaritySublocation,
	"Subsource":   polaritySubsource,
	"Primary":     polarityPrimary,
	"Reversed":    polarityReversed,
	"Method":      polarityMethod,
	"Citation":    polarityCitation,
	"Start Date":  polarityStart,
	"End Date":    polarityEnd,
}

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
	var data [][]string

	data = append(data, polarityHeaders.Columns())

	for _, row := range p {
		primary := strings.TrimSpace(row.primary)
		if primary == "" && row.Primary {
			primary = strconv.FormatBool(row.Primary)
		}

		reversed := strings.TrimSpace(row.reversed)
		if reversed == "" && row.Reversed {
			reversed = strconv.FormatBool(row.Reversed)
		}

		data = append(data, []string{
			strings.TrimSpace(row.Station),
			strings.TrimSpace(row.Location),
			strings.TrimSpace(row.Sublocation),
			strings.TrimSpace(row.Subsource),
			primary,
			reversed,
			strings.TrimSpace(row.Method),
			strings.TrimSpace(row.Citation),
			row.Start.Format(DateTimeFormat),
			row.End.Format(DateTimeFormat),
		})
	}

	return data
}

func (p *PolarityList) decode(data [][]string) error {
	if !(len(data) > 1) {
		return nil
	}

	var polarities []Polarity

	fields := polarityHeaders.Fields(data[0])
	for _, row := range data[1:] {
		d := fields.Remap(row)

		start, err := time.Parse(DateTimeFormat, d[polarityStart])
		if err != nil {
			return err
		}

		end, err := time.Parse(DateTimeFormat, d[polarityEnd])
		if err != nil {
			return err
		}

		var primary bool
		if s := d[polarityPrimary]; s != "" {
			v, err := strconv.ParseBool(s)
			if err != nil {
				return err
			}
			primary = v
		}

		var reversed bool
		if s := d[polarityReversed]; s != "" {
			v, err := strconv.ParseBool(s)
			if err != nil {
				return err
			}
			reversed = v
		}

		polarities = append(polarities, Polarity{
			Station:     strings.TrimSpace(d[polarityStation]),
			Location:    strings.TrimSpace(d[polarityLocation]),
			Sublocation: strings.TrimSpace(d[polaritySublocation]),
			Subsource:   strings.TrimSpace(d[polaritySubsource]),
			Reversed:    reversed,
			Primary:     primary,
			Method:      strings.TrimSpace(d[polarityMethod]),
			Citation:    strings.TrimSpace(d[polarityCitation]),
			Span: Span{
				Start: start,
				End:   end,
			},
			primary:  strings.TrimSpace(d[polarityPrimary]),
			reversed: strings.TrimSpace(d[polarityReversed]),
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
