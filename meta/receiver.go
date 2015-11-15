package meta

import (
	"sort"
	"time"
)

type DeployedReceiver struct {
	Make      string    `csv:"Receiver Make"`
	Model     string    `csv:"Receiver Model"`
	Serial    string    `csv:"Serial Number"`
	Place     string    `csv:"Deployment Place"`
	StartTime time.Time `csv:"Installation Date"`
	EndTime   time.Time `csv:"Removal Date"`
}

type DeployedReceivers []DeployedReceiver

func (dr DeployedReceivers) Len() int      { return len(dr) }
func (dr DeployedReceivers) Swap(i, j int) { dr[i], dr[j] = dr[j], dr[i] }
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
