package meta

import (
	"sort"
	"time"
)

type InstalledRadome struct {
	Make   string `csv:"Sensor Make"`
	Model  string `csv:"Sensor Model"`
	Serial string `csv:"Serial Number"`
	Mark   string `csv:"Mark Code"`
	/*
		Location  string    `csv:"Location Code"`
		Azimuth   float64   `csv:"Installation Azimuth"`
		Dip       float64   `csv:"Installation Dip"`
		Depth     float64   `csv:"Installation Depth"`
	*/
	StartTime time.Time `csv:"Installation Date"`
	EndTime   time.Time `csv:"Removal Date"`
}

type InstalledRadomes []InstalledRadome

func (ir InstalledRadomes) Len() int      { return len(ir) }
func (ir InstalledRadomes) Swap(i, j int) { ir[i], ir[j] = ir[j], ir[i] }
func (ir InstalledRadomes) Less(i, j int) bool {
	switch {
	case ir[i].Make < ir[j].Make:
		return true
	case ir[i].Make > ir[j].Make:
		return false
	case ir[i].Model < ir[j].Model:
		return true
	case ir[i].Model > ir[j].Model:
		return false
	case Serial(ir[i].Serial).Less(Serial(ir[j].Serial)):
		return true
	case Serial(ir[i].Serial).Greater(Serial(ir[j].Serial)):
		return false
	case ir[i].StartTime.Before(ir[j].StartTime):
		return true
	default:
		return false
	}
}

func (ir InstalledRadomes) List()      {}
func (ir InstalledRadomes) Sort() List { sort.Sort(ir); return ir }
