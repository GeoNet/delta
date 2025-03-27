package meta

import (
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

var connectionHeaders Header = map[string]int{
	"Station":    connectionStation,
	"Location":   connectionLocation,
	"Place":      connectionPlace,
	"Role":       connectionRole,
	"Number":     connectionNumber,
	"Start Date": connectionStart,
	"End Date":   connectionEnd,
}

var ConnectionTable Table = Table{
	name:    "Connection",
	headers: connectionHeaders,
	primary: []string{"Station", "Location", "Place", "Number", "Start Date"},
	foreign: map[string][]string{
		"Site": {"Station", "Location"},
	},
	remap: map[string]string{
		"Start Date": "Start",
		"End Date":   "End",
	},
	start: "Start Date",
	end:   "End Date",
}

type Connection struct {
	Span

	Station  string `json:"station"`
	Location string `json:"location"`
	Place    string `json:"place"`
	Role     string `json:"role,omitempty"`
	Number   int    `json:"number,omitempty"`

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
	var data [][]string

	data = append(data, connectionHeaders.Columns())

	for _, row := range c {
		data = append(data, []string{
			strings.TrimSpace(row.Station),
			strings.TrimSpace(row.Location),
			strings.TrimSpace(row.Place),
			strings.TrimSpace(row.Role),
			strings.TrimSpace(row.number),
			row.Start.Format(DateTimeFormat),
			row.End.Format(DateTimeFormat),
		})
	}

	return data
}

func (c *ConnectionList) decode(data [][]string) error {
	if !(len(data) > 1) {
		return nil
	}

	var connections []Connection

	fields := connectionHeaders.Fields(data[0])
	for _, row := range data[1:] {
		d := fields.Remap(row)

		number, err := ParseInt(strings.TrimSpace(d[connectionNumber]))
		if err != nil {
			return err
		}

		start, err := time.Parse(DateTimeFormat, strings.TrimSpace(d[connectionStart]))
		if err != nil {
			return err
		}

		end, err := time.Parse(DateTimeFormat, strings.TrimSpace(d[connectionEnd]))
		if err != nil {
			return err
		}

		connections = append(connections, Connection{
			Station:  strings.TrimSpace(d[connectionStation]),
			Location: strings.TrimSpace(d[connectionLocation]),
			Place:    strings.TrimSpace(d[connectionPlace]),
			Role:     strings.TrimSpace(d[connectionRole]),
			Number:   number,
			Span: Span{
				Start: start,
				End:   end,
			},

			number: strings.TrimSpace(d[connectionNumber]),
		})
	}

	*c = ConnectionList(connections)

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
