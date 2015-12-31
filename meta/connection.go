package meta

/*
import (
	"sort"
	"time"
)
*/

type Connection struct {
	Span

	StationCode  string
	LocationCode string
	Place        string
	Role         string
}

type Connections []Connection

func (c Connection) less(con Connection) bool {
	switch {
	case c.StationCode < con.StationCode:
		return true
	case c.StationCode > con.StationCode:
		return false
	case c.LocationCode < con.LocationCode:
		return true
	case c.LocationCode > con.LocationCode:
		return false
	case c.Place < con.Place:
		return true
	case c.Place > con.Place:
		return false
	case c.Role < con.Role:
		return true
	case c.Role > con.Role:
		return false
	default:
		return c.Span.before(con.Span)
	}
}

/*
func (c Connections) List()      {}
func (c Connections) Sort() List { sort.Sort(c); return c }
*/

func (c Connections) Len() int           { return len(c) }
func (c Connections) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c Connections) Less(i, j int) bool { return c[i].less(c[j]) }
