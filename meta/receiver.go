package meta

/*

import (
	"sort"
	"time"
)
*/

type DeployedReceiver struct {
	Install

	Place string
}

type DeployedReceivers []DeployedReceiver

func (r DeployedReceivers) Len() int           { return len(r) }
func (r DeployedReceivers) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r DeployedReceivers) Less(i, j int) bool { return r[i].Install.less(r[j].Install) }

/*
func (dr DeployedReceivers) Less(i, j int) bool {
	switch {
	case dr[i].Make < dr[j].Make:
		return true
	case dr[i].Make > dr[j].Make:
		return false
	case dr[i].Model < dr[j].Model:
		return true
	case dr[i].Model > dr[j].Model:
		return false
	case Serial(dr[i].Serial).Less(Serial(dr[j].Serial)):
		return true
	case Serial(dr[i].Serial).Greater(Serial(dr[j].Serial)):
		return false
	case dr[i].StartTime.Before(dr[j].StartTime):
		return true
	default:
		return false
	}
}

func (dr DeployedReceivers) List()      {}
func (dr DeployedReceivers) Sort() List { sort.Sort(dr); return dr }
*/
