package meta

import (
	//	"strconv"
	"time"
)

const DateTimeFormat = "2006-01-02T15:04:05Z"

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
	Height float64
	North  float64
	East   float64
}

type Scale struct {
	Factor float64
	Bias   float64
}

type Span struct {
	Start time.Time
	End   time.Time
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
