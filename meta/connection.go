package meta

import (
	"sort"
	"time"
)

type Connection struct {
	Station   string    `csv:"Station Code"`
	Location  string    `csv:"Location Code"`
	Place     string    `csv:"Datalogger Place"`
	Role      string    `csv:"Datalogger Role"`
	StartTime time.Time `csv:"Start Date"`
	EndTime   time.Time `csv:"End Date"`
}

type Connections []Connection

func (c Connections) Len() int      { return len(c) }
func (c Connections) Swap(i, j int) { c[i], c[j] = c[j], c[i] }
func (c Connections) Less(i, j int) bool {
	switch {
	case c[i].Station < c[j].Station:
		return true
	case c[i].Station > c[j].Station:
		return false
	case c[i].Location < c[j].Location:
		return true
	case c[i].Location > c[j].Location:
		return false
	case c[i].Place < c[j].Place:
		return true
	case c[i].Place > c[j].Place:
		return false
	case c[i].Role < c[j].Role:
		return true
	case c[i].Role > c[j].Role:
		return false
	case c[i].StartTime.Before(c[j].StartTime):
		return true
	case c[i].StartTime.After(c[j].StartTime):
		return false
		/*
			case c[i].Offset < c[j].Offset:
				return true
			case c[i].Offset > c[j].Offset:
				return false
		*/
	default:
		return false
	}
}

func (c Connections) List()      {}
func (c Connections) Sort() List { sort.Sort(c); return c }
