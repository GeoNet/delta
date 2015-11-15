package meta

import (
	"sort"
	"time"
)

type InstalledSensor struct {
	Make      string    `csv:"Sensor Make"`
	Model     string    `csv:"Sensor Model"`
	Serial    string    `csv:"Serial Number"`
	Station   string    `csv:"Station Code"`
	Location  string    `csv:"Location Code"`
	Azimuth   float64   `csv:"Installation Azimuth"`
	Dip       float64   `csv:"Installation Dip"`
	Depth     float64   `csv:"Installation Depth"`
	StartTime time.Time `csv:"Installation Date"`
	EndTime   time.Time `csv:"Removal Date"`
}

type InstalledSensors []InstalledSensor

func (is InstalledSensors) Len() int      { return len(is) }
func (is InstalledSensors) Swap(i, j int) { is[i], is[j] = is[j], is[i] }
func (is InstalledSensors) Less(i, j int) bool {
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

func (is InstalledSensors) List()      {}
func (is InstalledSensors) Sort() List { sort.Sort(is); return is }
