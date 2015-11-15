package meta

import (
	"sort"
	"time"
)

type Mark struct {
	Code        string    `csv:"Mark Code"`
	Network     string    `csv:"Network Code"`
	Name        string    `csv:"Mark Name"`
	Type        string    `csv:"Mark Type"`
	Latitude    float64   `csv:"Latitude"`
	Longitude   float64   `csv:"Longitude"`
	Height      float64   `csv:"Height"`
	Datum       string    `csv:"Datum"`
	Offset      float64   `csv:"Vertical Offset"`
	StartTime   time.Time `csv:"Start Time"`
	EndTime     time.Time `csv:"End Time"`
	Dome        string    `csv:"Dome Number"`
	Plan        string    `csv:"Plan Reference"`
	Protection  string    `csv:"Protection"`
	Sky         string    `csv:"Sky View"`
	Established string    `csv:"Established By"`
}

type Marks []Mark

func (m Marks) Len() int           { return len(m) }
func (m Marks) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }
func (m Marks) Less(i, j int) bool { return m[i].Code < m[j].Code }

func (m Marks) List()      {}
func (m Marks) Sort() List { sort.Sort(m); return m }
