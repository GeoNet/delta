package meta

import (
	"sort"
	"time"
)

type InstalledMetSensor struct {
	Make      string    `csv:"Met Sensor Make"`
	Model     string    `csv:"Met Sensor Model"`
	Serial    string    `csv:"Serial Number"`
	Mark      string    `csv:"Mark"`
	Comment   string    `csv:"IMS Comment"`
	Latitude  float64   `csv:"Latitude"`
	Longitude float64   `csv:"Longitude"`
	Height    float64   `csv:"Height"`
	Datum     string    `csv:"Datum"`
	StartTime time.Time `csv:"Installation Date"`
	EndTime   time.Time `csv:"Removal Date"`
}

type InstalledMetSensors []InstalledMetSensor

func (is InstalledMetSensors) Len() int      { return len(is) }
func (is InstalledMetSensors) Swap(i, j int) { is[i], is[j] = is[j], is[i] }
func (is InstalledMetSensors) Less(i, j int) bool {
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

func (is InstalledMetSensors) List()      {}
func (is InstalledMetSensors) Sort() List { sort.Sort(is); return is }
