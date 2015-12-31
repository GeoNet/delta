package meta

/*
import (
	"sort"
	"time"
)
*/

type Install struct {
	Equipment
	Span
}

func (i Install) less(in Install) bool {
	switch {
	case i.Equipment.less(in.Equipment):
		return true
	case in.Equipment.less(i.Equipment):
		return false
	default:
		return i.Span.before(in.Span)
	}
}

/*
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
*/
