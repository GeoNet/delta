package meta

/*

import (
	"sort"
	"time"
)
*/

type InstalledRadome struct {
	Install

	MarkCode string
}

type InstalledRadomes []InstalledRadome

func (r InstalledRadomes) Len() int           { return len(r) }
func (r InstalledRadomes) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r InstalledRadomes) Less(i, j int) bool { return r[i].Install.less(r[j].Install) }

/*
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
*/
