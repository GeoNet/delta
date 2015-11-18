package meta

import (
	"sort"
	"time"
)

type DeployedDatalogger struct {
	Make      string    `csv:"Datalogger Make"`
	Model     string    `csv:"Datalogger Model"`
	Serial    string    `csv:"Serial Number"`
	Place     string    `csv:"Deployment Place"`
	Role      string    `csv:"Deployment Role"`
	StartTime time.Time `csv:"Installation Date"`
	EndTime   time.Time `csv:"Removal Date"`
}

type DeployedDataloggers []DeployedDatalogger

func (dd DeployedDataloggers) Len() int      { return len(dd) }
func (dd DeployedDataloggers) Swap(i, j int) { dd[i], dd[j] = dd[j], dd[i] }
func (dd DeployedDataloggers) Less(i, j int) bool {
	switch {
	case dd[i].Make < dd[j].Make:
		return true
	case dd[i].Make > dd[j].Make:
		return false
	case dd[i].Model < dd[j].Model:
		return true
	case dd[i].Model > dd[j].Model:
		return false
	case Serial(dd[i].Serial).Less(Serial(dd[j].Serial)):
		return true
	case Serial(dd[i].Serial).Greater(Serial(dd[j].Serial)):
		return false
	case dd[i].StartTime.Before(dd[j].StartTime):
		return true
	default:
		return false
	}
}

func (dd DeployedDataloggers) List()      {}
func (dd DeployedDataloggers) Sort() List { sort.Sort(dd); return dd }
