package meta

/*
import (
	"sort"
	"time"
)
*/

type InstalledSensor struct {
	Install
	Orientation
	Offset

	StationCode  string `csv:"Station Code"`
	LocationCode string `csv:"Location Code"`
}

type InstalledSensors []InstalledSensor

func (s InstalledSensors) Len() int           { return len(s) }
func (s InstalledSensors) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s InstalledSensors) Less(i, j int) bool { return s[i].Install.less(s[j].Install) }

/*
func (is InstalledSensors) Less(i, j int) bool {
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

func (is InstalledSensors) List()      {}
func (is InstalledSensors) Sort() List { sort.Sort(is); return is }
*/
