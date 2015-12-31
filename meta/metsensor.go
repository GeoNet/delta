package meta

/*
import (
	"sort"
	"time"
)
*/

type InstalledMetSensor struct {
	Install
	Point

	MarkCode string
	Comment  string
}

type InstalledMetSensors []InstalledMetSensor

func (m InstalledMetSensors) Len() int      { return len(m) }
func (m InstalledMetSensors) Swap(i, j int) { m[i], m[j] = m[j], m[i] }
func (m InstalledMetSensors) Less(i, j int) { m[i].Install.less(m[j].Install) }

/*
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
*/
