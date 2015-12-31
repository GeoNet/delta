package meta

import (
//	"sort"
//	"time"
)

type Mark struct {
	Reference
	Point
	Span

	Type        string  `csv:"Mark Type"`
	Offset      float64 `csv:"Ground Relationship"`
	Dome        string  `csv:"Dome Number"`
	Plan        string  `csv:"Plan Reference"`
	Protection  string  `csv:"Protection"`
	Sky         string  `csv:"Sky View"`
	Monument    string  `csv:"Monument Type"`
	Established string  `csv:"Established By"`
}

type Marks []Mark

func (m Marks) Len() int           { return len(m) }
func (m Marks) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }
func (m Marks) Less(i, j int) bool { return m[i].Reference.less(m[j].Reference) }

/*
func (m Marks) List()      {}
func (m Marks) Sort() List { sort.Sort(m); return m }
*/
