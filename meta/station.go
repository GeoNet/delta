package meta

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Station struct {
	Reference
	Point
	Span

	Notes string
}

type Stations []Station

func (s Stations) Len() int           { return len(s) }
func (s Stations) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s Stations) Less(i, j int) bool { return s[i].Reference.Less(s[j].Reference) }

func (s Stations) encode() [][]string {
	data := [][]string{{
		"Station Code",
		"Network Code",
		"Station Name",
		"Latitude",
		"Longitude",
		"Height",
		"Datum",
		"Start Time",
		"End Time",
		"Notes",
	}}
	for _, v := range s {
		data = append(data, []string{
			strings.TrimSpace(v.Code),
			strings.TrimSpace(v.Network),
			strings.TrimSpace(v.Name),
			strconv.FormatFloat(v.Latitude, 'g', -1, 64),
			strconv.FormatFloat(v.Longitude, 'g', -1, 64),
			strconv.FormatFloat(v.Elevation, 'g', -1, 64),
			strings.TrimSpace(v.Datum),
			v.Start.Format(DateTimeFormat),
			v.End.Format(DateTimeFormat),
			strings.TrimSpace(v.Notes),
		})
	}
	return data
}

func (s *Stations) decode(data [][]string) error {
	var stations []Station
	if len(data) > 1 {
		for _, d := range data[1:] {
			if len(d) != 10 {
				return fmt.Errorf("incorrect number of installed station fields")
			}
			var err error

			var lat, lon, elev float64
			if lat, err = strconv.ParseFloat(d[3], 64); err != nil {
				return err
			}
			if lon, err = strconv.ParseFloat(d[4], 64); err != nil {
				return err
			}
			if elev, err = strconv.ParseFloat(d[5], 64); err != nil {
				return err
			}

			var start, end time.Time
			if start, err = time.Parse(DateTimeFormat, d[7]); err != nil {
				return err
			}
			if end, err = time.Parse(DateTimeFormat, d[8]); err != nil {
				return err
			}

			stations = append(stations, Station{
				Reference: Reference{
					Code:    strings.TrimSpace(d[0]),
					Network: strings.TrimSpace(d[1]),
					Name:    strings.TrimSpace(d[2]),
				},
				Span: Span{
					Start: start,
					End:   end,
				},
				Point: Point{
					Latitude:  lat,
					Longitude: lon,
					Elevation: elev,
					Datum:     strings.TrimSpace(d[6]),
				},
				Notes: strings.TrimSpace(d[9]),
			})
		}

		*s = Stations(stations)
	}
	return nil
}

func LoadStations(path string) ([]Station, error) {
	var s []Station

	if err := LoadList(path, (*Stations)(&s)); err != nil {
		return nil, err
	}

	sort.Sort(Stations(s))

	return s, nil
}
