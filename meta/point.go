package meta

import (
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	pointSample = iota
	pointLocation
	pointLatitude
	pointLongitude
	pointElevation
	pointDepth
	pointDatum
	pointSurvey
	pointStart
	pointEnd
	pointLast
)

var pointHeaders Header = map[string]int{
	"Sample":     pointSample,
	"Location":   pointLocation,
	"Latitude":   pointLatitude,
	"Longitude":  pointLongitude,
	"Elevation":  pointElevation,
	"Depth":      pointDepth,
	"Datum":      pointDatum,
	"Survey":     pointSurvey,
	"Start Date": pointStart,
	"End Date":   pointEnd,
}

var PointTable Table = Table{
	name:    "Point",
	headers: pointHeaders,
	primary: []string{"Sample", "Location", "Start Date"},
	native:  []string{"Latitude", "Longitude", "Elevation", "Depth"},
	foreign: map[string]map[string]string{
		"Sample": {"Code": "Sample"},
	},
	nullable: []string{"Depth"},
	remap: map[string]string{
		"Start Date": "Start",
		"End Date":   "End",
	},
	start: "Start Date",
	end:   "End Date",
}

type Point struct {
	Position
	Span

	Sample   string `json:"sample"`
	Location string `json:"location"`
	Survey   string `json:"survey,omitempty"`
}

func (s Point) Less(point Point) bool {
	switch {
	case s.Sample < point.Sample:
		return true
	case s.Sample > point.Sample:
		return false
	case s.Location < point.Location:
		return true
	case s.Location > point.Location:
		return false
	default:
		return false
	}
}

type PointList []Point

func (s PointList) Len() int           { return len(s) }
func (s PointList) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s PointList) Less(i, j int) bool { return s[i].Less(s[j]) }

func (s PointList) encode() [][]string {
	var data [][]string

	data = append(data, pointHeaders.Columns())
	for _, row := range s {
		data = append(data, []string{
			strings.TrimSpace(row.Sample),
			strings.TrimSpace(row.Location),
			strings.TrimSpace(row.latitude),
			strings.TrimSpace(row.longitude),
			strings.TrimSpace(row.elevation),
			strings.TrimSpace(row.depth),
			strings.TrimSpace(row.Datum),
			strings.TrimSpace(row.Survey),
			row.Start.Format(DateTimeFormat),
			row.End.Format(DateTimeFormat),
		})
	}

	return data
}

func (s *PointList) decode(data [][]string) error {
	var points []Point

	if !(len(data) > 1) {
		return nil
	}

	fields := pointHeaders.Fields(data[0])
	for _, row := range data[1:] {
		d := fields.Remap(row)

		lat, err := strconv.ParseFloat(d[pointLatitude], 64)
		if err != nil {
			return err
		}
		lon, err := strconv.ParseFloat(d[pointLongitude], 64)
		if err != nil {
			return err
		}
		elev, err := ParseFloat64(d[pointElevation])
		if err != nil {
			return err
		}
		depth, err := ParseFloat64(d[pointDepth])
		if err != nil {
			return err
		}

		start, err := time.Parse(DateTimeFormat, d[pointStart])
		if err != nil {
			return err
		}
		end, err := time.Parse(DateTimeFormat, d[pointEnd])
		if err != nil {
			return err
		}

		points = append(points, Point{
			Position: Position{
				// geographic details
				Latitude:  lat,
				Longitude: lon,
				Elevation: elev,
				Depth:     depth,
				Datum:     strings.TrimSpace(d[pointDatum]),

				// shadow variables
				latitude:  d[pointLatitude],
				longitude: d[pointLongitude],
				elevation: d[pointElevation],
				depth:     d[pointDepth],
			},
			Span: Span{
				Start: start,
				End:   end,
			},
			Sample:   strings.TrimSpace(d[pointSample]),
			Location: strings.TrimSpace(d[pointLocation]),
			Survey:   strings.TrimSpace(d[pointSurvey]),
		})
	}

	*s = PointList(points)

	return nil
}

func LoadPoints(path string) ([]Point, error) {
	var s []Point

	if err := LoadList(path, (*PointList)(&s)); err != nil {
		return nil, err
	}

	sort.Sort(PointList(s))

	return s, nil
}
