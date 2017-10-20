package meta

import (
	"strconv"
	"strings"
	"time"
)

const DateTimeFormat = "2006-01-02T15:04:05Z"

type Compare int

const (
	EqualTo Compare = iota
	LessThan
	GreaterThan
)

type Reference struct {
	Code    string
	Network string
	Name    string
}

type Point struct {
	Latitude  float64
	Longitude float64
	Elevation float64
	Datum     string
}

type Orientation struct {
	Dip     float64
	Azimuth float64
}

type Offset struct {
	Vertical float64
	North    float64
	East     float64
}

type Scale struct {
	Factor float64
	Bias   float64
}

type Span struct {
	Start time.Time
	End   time.Time
}

type Range struct {
	Value   float64
	Compare Compare
}

func NewRange(s string) (Range, error) {
	switch {
	case strings.HasPrefix(s, "<"):
		v, err := strconv.ParseFloat(s[1:], 64)
		if err != nil {
			return Range{}, err
		}
		return Range{
			Value:   v,
			Compare: LessThan,
		}, nil
	case strings.HasPrefix(s, ">"):
		v, err := strconv.ParseFloat(s[1:], 64)
		if err != nil {
			return Range{}, err
		}
		return Range{
			Value:   v,
			Compare: GreaterThan,
		}, nil
	default:
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return Range{}, err
		}
		return Range{
			Value: v,
		}, nil
	}
}

func (r Range) String() string {
	switch r.Compare {
	case LessThan:
		return "<" + strconv.FormatFloat(r.Value, 'g', -1, 64)
	case GreaterThan:
		return ">" + strconv.FormatFloat(r.Value, 'g', -1, 64)
	default:
		return strconv.FormatFloat(r.Value, 'g', -1, 64)
	}
}

type Equipment struct {
	Make   string
	Model  string
	Serial string
}

func (e Equipment) String() string {
	return e.Make + " " + e.Model + " [" + e.Serial + "]"
}

func (e Equipment) Less(eq Equipment) bool {

	switch {
	case e.Make < eq.Make:
		return true
	case e.Make > eq.Make:
		return false
	case e.Model < eq.Model:
		return true
	case e.Model > eq.Model:
		return false
	}

	// too many edge cases depending on the original slice order
	/*
		if a, err := strconv.Atoi(e.Serial); err == nil {
			if b, err := strconv.Atoi(eq.Serial); err == nil {
				return a < b
			}
		}
	*/

	return e.Serial < eq.Serial
}

type Install struct {
	Equipment
	Span
}

func (i Install) less(in Install) bool {
	switch {
	case i.Equipment.Less(in.Equipment):
		return true
	case in.Equipment.Less(i.Equipment):
		return false
	default:
		return i.Start.Before(in.Start)
	}
}
