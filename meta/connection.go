package meta

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

const (
	connectionStation int = iota
	connectionLocation
	connectionPlace
	connectionRole
	connectionStart
	connectionEnd
	connectionLast
)

type Connection struct {
	Span

	Station  string
	Location string
	Place    string
	Role     string
}

func (c Connection) less(con Connection) bool {
	switch {
	case c.Station < con.Station:
		return true
	case c.Station > con.Station:
		return false
	case c.Location < con.Location:
		return true
	case c.Location > con.Location:
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
		return c.Start.Before(con.Start)
	}
}

type ConnectionList []Connection

func (c ConnectionList) Len() int           { return len(c) }
func (c ConnectionList) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c ConnectionList) Less(i, j int) bool { return c[i].less(c[j]) }

func (c ConnectionList) encode() [][]string {
	data := [][]string{{
		"Station Code",
		"Location Code",
		"Datalogger Place",
		"Datalogger Role",
		"Start Date",
		"End Date",
	}}
	for _, v := range c {
		data = append(data, []string{
			strings.TrimSpace(v.Station),
			strings.TrimSpace(v.Location),
			strings.TrimSpace(v.Place),
			strings.TrimSpace(v.Role),
			v.Start.Format(DateTimeFormat),
			v.End.Format(DateTimeFormat),
		})
	}
	return data
}

func (c *ConnectionList) decode(data [][]string) error {
	var connections []Connection
	if len(data) > 1 {
		for _, v := range data[1:] {
			if len(v) != connectionLast {
				return fmt.Errorf("incorrect number of installed connection fields")
			}
			var err error

			var start, end time.Time
			if start, err = time.Parse(DateTimeFormat, v[connectionStart]); err != nil {
				return err
			}
			if end, err = time.Parse(DateTimeFormat, v[connectionEnd]); err != nil {
				return err
			}

			connections = append(connections, Connection{
				Station:  strings.TrimSpace(v[connectionStation]),
				Location: strings.TrimSpace(v[connectionLocation]),
				Place:    strings.TrimSpace(v[connectionPlace]),
				Role:     strings.TrimSpace(v[connectionRole]),
				Span: Span{
					Start: start,
					End:   end,
				},
			})
		}

		*c = ConnectionList(connections)
	}
	return nil
}

func LoadConnections(path string) ([]Connection, error) {
	var c []Connection

	if err := LoadList(path, (*ConnectionList)(&c)); err != nil {
		return nil, err
	}

	sort.Sort(ConnectionList(c))

	return c, nil
}
