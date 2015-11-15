package meta

import (
	"sort"
	"time"
)

type InstalledAntenna struct {
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

type InstalledAntennas []InstalledAntenna

func (ia InstalledAntennas) Len() int      { return len(ia) }
func (ia InstalledAntennas) Swap(i, j int) { ia[i], ia[j] = ia[j], ia[i] }
func (ia InstalledAntennas) Less(i, j int) bool {
	switch {
	case ia[i].Make < ia[j].Make:
		return true
	case ia[i].Make > ia[j].Make:
		return false
	case ia[i].Model < ia[j].Model:
		return true
	case ia[i].Model > ia[j].Model:
		return false
	case Serial(ia[i].Serial).Less(Serial(ia[j].Serial)):
		return true
	case Serial(ia[i].Serial).Greater(Serial(ia[j].Serial)):
		return false
	case ia[i].StartTime.Before(ia[j].StartTime):
		return true
	default:
		return false
	}
}

func (ia InstalledAntennas) List()      {}
func (ia InstalledAntennas) Sort() List { sort.Sort(ia); return ia }
