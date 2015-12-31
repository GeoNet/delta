package meta

/*

import (
	"sort"
	"time"
)
*/

type InstalledGauge struct {
	Install
	Offset
	Orientation

	StationCode  string
	LocationCode string
	Cable        float64
}

type InstalledGauges []InstalledGauge

func (g InstalledGauges) Len() int      { return len(g) }
func (g InstalledGauges) Swap(i, j int) { g[i], g[j] = g[j], g[i] }
func (g InstalledGauges) less(i, j int) { g[i].Install.less(g[j].Install) }

/*
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
*/
