package meta

import (
	"sort"
	"time"
)

type Monument struct {
	Mark         string    `csv:"Monument Mark"`
	Code         string    `csv:"Monument Code"`
	Eccentricity string    `csv:"Eccentricity"`
	StartTime    time.Time `csv:"Date Established"`
	EndTime      time.Time `csv:"Date Removed"`
}

type Monuments []Monument

func (m Monuments) Len() int      { return len(m) }
func (m Monuments) Swap(i, j int) { m[i], m[j] = m[j], m[i] }
func (m Monuments) Less(i, j int) bool {
	switch {
	case m[i].Mark < m[j].Mark:
		return true
	case m[i].Mark > m[j].Mark:
		return false
	default:
		return m[i].StartTime.Before(m[j].StartTime)
	}
}

func (m Monuments) List()      {}
func (m Monuments) Sort() List { sort.Sort(m); return m }
