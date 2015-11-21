package meta

import (
	"sort"
)

type Response struct {
	Datalogger string `csv:"Datalogger Model"`
	Sensor     string `csv:"Sensor Model"`
	Reversed   bool   `csv:"Reversed Connection"`
	Lookup     string `csv:"Response Lookup"`
	Match      string `csv:"Site Match"`
}

type Responses []Response

func (r Responses) Len() int      { return len(r) }
func (r Responses) Swap(i, j int) { r[i], r[j] = r[j], r[i] }
func (r Responses) Less(i, j int) bool {
	switch {
	case r[i].Datalogger < r[j].Datalogger:
		return true
	case r[i].Datalogger > r[j].Datalogger:
		return false
	case r[i].Sensor < r[j].Sensor:
		return true
	case r[i].Sensor > r[j].Sensor:
		return false
	case r[i].Match < r[j].Match:
		return true
	case r[i].Match > r[j].Match:
		return false
	default:
		return false
	}
}

func (r Responses) List()      {}
func (r Responses) Sort() List { sort.Sort(r); return r }
