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
	connectionNumber
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
	Number   int

	number string
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
	case c.Number < con.Number:
		return true
	case c.Number > con.Number:
		return false
	case c.Start.Before(con.Start):
		return true
	default:
		return false
	}
}

type ConnectionList []Connection

func (c ConnectionList) Len() int           { return len(c) }
func (c ConnectionList) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c ConnectionList) Less(i, j int) bool { return c[i].less(c[j]) }

func (c ConnectionList) encode() [][]string {
	data := [][]string{{
		"Station",
		"Location",
		"Place",
		"Role",
		"Number",
		"Start Date",
		"End Date",
	}}
	for _, v := range c {
		data = append(data, []string{
			strings.TrimSpace(v.Station),
			strings.TrimSpace(v.Location),
			strings.TrimSpace(v.Place),
			strings.TrimSpace(v.Role),
			strings.TrimSpace(v.number),
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

			number, err := ParseInt(strings.TrimSpace(v[connectionNumber]))
			if err != nil {
				return err
			}

			start, err := time.Parse(DateTimeFormat, strings.TrimSpace(v[connectionStart]))
			if err != nil {
				return err
			}

			end, err := time.Parse(DateTimeFormat, strings.TrimSpace(v[connectionEnd]))
			if err != nil {
				return err
			}

			connections = append(connections, Connection{
				Station:  strings.TrimSpace(v[connectionStation]),
				Location: strings.TrimSpace(v[connectionLocation]),
				Place:    strings.TrimSpace(v[connectionPlace]),
				Role:     strings.TrimSpace(v[connectionRole]),
				Number:   number,
				Span: Span{
					Start: start,
					End:   end,
				},

				number: strings.TrimSpace(v[connectionNumber]),
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
