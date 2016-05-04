package meta

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type Connection struct {
	Span

	StationCode  string
	LocationCode string
	Place        string
	Role         string
}

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
			strings.TrimSpace(v.StationCode),
			strings.TrimSpace(v.LocationCode),
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
			if len(v) != 6 {
				return fmt.Errorf("incorrect number of installed connection fields")
			}
			var err error

			var start, end time.Time
			if start, err = time.Parse(DateTimeFormat, v[4]); err != nil {
				return err
			}
			if end, err = time.Parse(DateTimeFormat, v[5]); err != nil {
				return err
			}

			connections = append(connections, Connection{
				StationCode:  strings.TrimSpace(v[0]),
				LocationCode: strings.TrimSpace(v[1]),
				Place:        strings.TrimSpace(v[2]),
				Role:         strings.TrimSpace(v[3]),
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
