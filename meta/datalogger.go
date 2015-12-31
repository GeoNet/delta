package meta

/*
import (
	"sort"
	"time"
)
*/

type DeployedDatalogger struct {
	Install

	Place string
	Role  string
}

type DeployedDataloggers []DeployedDatalogger

func (d DeployedDataloggers) Len() int           { return len(d) }
func (d DeployedDataloggers) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d DeployedDataloggers) Less(i, j int) bool { return d[i].Install.less(d[j].Install) }

/*
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
*/
