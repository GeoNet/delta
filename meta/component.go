package meta

import (
	"sort"
)

type Component struct {
	Model   string  `csv:"Sensor Model"`
	Pin     int32   `csv:"Sensor Pin"`
	Azimuth float64 `csv:"Installation Azimuth"`
	Dip     float64 `csv:"Installation Dip"`
}

type Components []Component

func (c Components) Len() int      { return len(c) }
func (c Components) Swap(i, j int) { c[i], c[j] = c[j], c[i] }
func (c Components) Less(i, j int) bool {
	switch {
	case c[i].Model < c[j].Model:
		return true
	case c[i].Model > c[j].Model:
		return false
	case c[i].Pin < c[j].Pin:
		return true
	case c[i].Pin > c[j].Pin:
		return false
	default:
		return false
	}
}

func (c Components) List()      {}
func (c Components) Sort() List { sort.Sort(c); return c }
