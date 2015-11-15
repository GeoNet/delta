package meta

import (
	"sort"
	"time"
)

type InstalledGauge struct {
	Make      string    `csv:"Gauge Make"`
	Model     string    `csv:"Gauge Model"`
	Serial    string    `csv:"Serial Number"`
	Station   string    `csv:"Station Code"`
	Location  string    `csv:"Location Code"`
	Dip       float64   `csv:"Installation Dip"`
	Vertical  float64   `csv:"Vertical Offset"`
	North     float64   `csv:"Offset North"`
	East      float64   `csv:"Offset East"`
	Cable     float64   `csv:"Cable Length"`
	StartTime time.Time `csv:"Installation Date"`
	EndTime   time.Time `csv:"Removal Date"`
}

type InstalledGauges []InstalledGauge

func (is InstalledGauges) Len() int      { return len(is) }
func (is InstalledGauges) Swap(i, j int) { is[i], is[j] = is[j], is[i] }
func (is InstalledGauges) Less(i, j int) bool {
	switch {
	case is[i].Make < is[j].Make:
		return true
	case is[i].Make > is[j].Make:
		return false
	case is[i].Model < is[j].Model:
		return true
	case is[i].Model > is[j].Model:
		return false
	case Serial(is[i].Serial).Less(Serial(is[j].Serial)):
		return true
	case Serial(is[i].Serial).Greater(Serial(is[j].Serial)):
		return false
	case is[i].StartTime.Before(is[j].StartTime):
		return true
	default:
		return false
	}
}

func (is InstalledGauges) List()      {}
func (is InstalledGauges) Sort() List { sort.Sort(is); return is }
