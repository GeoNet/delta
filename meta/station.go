package meta

import (
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	stationCode = iota
	stationNetwork
	stationName
	stationLatitude
	stationLongitude
	stationElevation
	stationDepth
	stationDatum
	stationStart
	stationEnd
	stationLast
)

var stationHeaders Header = map[string]int{
	"Station":    stationCode,
	"Network":    stationNetwork,
	"Name":       stationName,
	"Latitude":   stationLatitude,
	"Longitude":  stationLongitude,
	"Elevation":  stationElevation,
	"Depth":      stationDepth,
	"Datum":      stationDatum,
	"Start Date": stationStart,
	"End Date":   stationEnd,
}

type Station struct {
	Reference
	Position
	Span
}

var StationTable = Table{
	name:    "Station",
	headers: stationHeaders,
	primary: []string{"Station", "Start Date"},
	native:  []string{"Latitude", "Longitude", "Elevation", "Depth"},
	foreign: map[string]map[string]string{
		"Network": {"Network": "Network"},
	},
	nullable: []string{"Depth"},
	remap: map[string]string{
		"Station":    "Code",
		"Start Date": "Start",
		"End Date":   "End",
	},
	start: "Start Date",
	end:   "End Date",
}

type StationList []Station

func (s StationList) Len() int      { return len(s) }
func (s StationList) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s StationList) Less(i, j int) bool {
	switch {
	case s[i].Code < s[j].Code:
		return true
	case s[i].Code > s[j].Code:
		return false
	case s[i].Network < s[j].Network:
		return true
	default:
		return false
	}
}

func (s StationList) encode() [][]string {
	var data [][]string

	data = append(data, stationHeaders.Columns())
	for _, row := range s {
		data = append(data, []string{
			strings.TrimSpace(row.Code),
			strings.TrimSpace(row.Network),
			strings.TrimSpace(row.Name),
			strings.TrimSpace(row.latitude),
			strings.TrimSpace(row.longitude),
			strings.TrimSpace(row.elevation),
			strings.TrimSpace(row.depth),
			strings.TrimSpace(row.Datum),
			row.Start.Format(DateTimeFormat),
			row.End.Format(DateTimeFormat),
		})
	}

	return data
}

func (s *StationList) decode(data [][]string) error {
	if !(len(data) > 1) {
		return nil
	}

	var stations []Station

	fields := stationHeaders.Fields(data[0])
	for _, row := range data[1:] {
		d := fields.Remap(row)

		lat, err := strconv.ParseFloat(d[stationLatitude], 64)
		if err != nil {
			return err
		}

		lon, err := strconv.ParseFloat(d[stationLongitude], 64)
		if err != nil {
			return err
		}

		elev, err := ParseFloat64(d[stationElevation])
		if err != nil {
			return err
		}

		depth, err := ParseFloat64(d[stationDepth])
		if err != nil {
			return err
		}

		start, err := time.Parse(DateTimeFormat, d[stationStart])
		if err != nil {
			return err
		}
		end, err := time.Parse(DateTimeFormat, d[stationEnd])
		if err != nil {
			return err
		}

		stations = append(stations, Station{
			Reference: Reference{
				Code:    strings.TrimSpace(d[stationCode]),
				Network: strings.TrimSpace(d[stationNetwork]),
				Name:    strings.TrimSpace(d[stationName]),
			},
			Span: Span{
				Start: start,
				End:   end,
			},
			Position: Position{
				Latitude:  lat,
				Longitude: lon,
				Elevation: elev,
				Datum:     strings.TrimSpace(d[stationDatum]),
				Depth:     depth,

				latitude:  strings.TrimSpace(d[stationLatitude]),
				longitude: strings.TrimSpace(d[stationLongitude]),
				elevation: strings.TrimSpace(d[stationElevation]),
				depth:     strings.TrimSpace(d[stationDepth]),
			},
		})
	}

	*s = StationList(stations)

	return nil
}

func LoadStations(path string) ([]Station, error) {
	var s []Station

	if err := LoadList(path, (*StationList)(&s)); err != nil {
		return nil, err
	}

	sort.Sort(StationList(s))

	return s, nil
}
