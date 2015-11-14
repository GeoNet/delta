package meta

import (
	"sort"
	"time"
)

type Station struct {
	Code      string    `csv:"Station Code",`
	Network   string    `csv:"Network Code",`
	Name      string    `csv:"Station Name",`
	Latitude  float64   `csv:"Latitude",`
	Longitude float64   `csv:"Longitude",`
	Depth     float64   `csv:"Depth",`
	StartTime time.Time `csv:"Start Time"`
	EndTime   time.Time `csv:"End Time"`
}

type Stations []Station

func (s Stations) Len() int           { return len(s) }
func (s Stations) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s Stations) Less(i, j int) bool { return s[i].Code < s[j].Code }

func (s Stations) List()      {}
func (s Stations) Sort() List { sort.Sort(s); return s }
