package meta

import (
	"sort"
)

type Location struct {
	Station   string  `csv:"Station Code"`
	Code      string  `csv:"Location Code"`
	Latitude  float64 `csv:"Latitude"`
	Longitude float64 `csv:"Longitude"`
	Height    float64 `csv:"Height"`
	Datum     string  `csv:"Datum"`
}

type Locations []Location

func (l Locations) Len() int      { return len(l) }
func (l Locations) Swap(i, j int) { l[i], l[j] = l[j], l[i] }
func (l Locations) Less(i, j int) bool {
	switch {
	case l[i].Station < l[j].Station:
		return true
	case l[i].Station > l[j].Station:
		return false
	case l[i].Code < l[j].Code:
		return true
	case l[i].Code > l[j].Code:
		return false
	default:
		return false
	}
}

func (l Locations) List()      {}
func (l Locations) Sort() List { sort.Sort(l); return l }
