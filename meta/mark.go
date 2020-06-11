package meta

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	markCode = iota
	markNetwork
	markIgs
	markName
	markLatitude
	markLongitude
	markElevation
	markDatum
	markStartTime
	markEndTime
	markLast
)

type Mark struct {
	Reference
	Point
	Span

	Igs bool
}

type MarkList []Mark

func (m MarkList) Len() int           { return len(m) }
func (m MarkList) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }
func (m MarkList) Less(i, j int) bool { return m[i].Code < m[j].Code }

func (m MarkList) encode() [][]string {
	data := [][]string{{
		"Mark",
		"Network",
		"Igs",
		"Name",
		"Latitude",
		"Longitude",
		"Elevation",
		"Datum",
		"Start Date",
		"End Date",
	}}
	for _, v := range m {
		data = append(data, []string{
			strings.TrimSpace(v.Code),
			strings.TrimSpace(v.Network),
			func() string {
				if v.Igs {
					return "yes"
				}
				return "no"
			}(),
			strings.TrimSpace(v.Name),
			strings.TrimSpace(v.latitude),
			strings.TrimSpace(v.longitude),
			strings.TrimSpace(v.elevation),
			strings.TrimSpace(v.Datum),
			v.Start.Format(DateTimeFormat),
			v.End.Format(DateTimeFormat),
		})
	}
	return data
}

func (m *MarkList) decode(data [][]string) error {
	var marks []Mark
	if len(data) > 1 {
		for _, d := range data[1:] {
			if len(d) != markLast {
				return fmt.Errorf("incorrect number of installed mark fields")
			}
			var err error

			var igs bool
			if igs, err = strconv.ParseBool(d[markIgs]); err != nil {
				switch d[markIgs] {
				case "y", "Y", "yes", "YES":
					igs = true
				case "n", "N", "no", "NO":
					igs = false
				default:
					return err
				}
			}

			var lat, lon, elev float64
			if lat, err = strconv.ParseFloat(d[markLatitude], 64); err != nil {
				return err
			}
			if lon, err = strconv.ParseFloat(d[markLongitude], 64); err != nil {
				return err
			}
			if elev, err = strconv.ParseFloat(d[markElevation], 64); err != nil {
				return err
			}

			var start, end time.Time
			if start, err = time.Parse(DateTimeFormat, d[markStartTime]); err != nil {
				return err
			}
			if end, err = time.Parse(DateTimeFormat, d[markEndTime]); err != nil {
				return err
			}
			marks = append(marks, Mark{
				Reference: Reference{
					Code:    strings.TrimSpace(d[markCode]),
					Network: strings.TrimSpace(d[markNetwork]),
					Name:    strings.TrimSpace(d[markName]),
				},
				Span: Span{
					Start: start,
					End:   end,
				},
				Point: Point{
					Latitude:  lat,
					Longitude: lon,
					Elevation: elev,
					Datum:     strings.TrimSpace(d[markDatum]),

					latitude:  strings.TrimSpace(d[markLatitude]),
					longitude: strings.TrimSpace(d[markLongitude]),
					elevation: strings.TrimSpace(d[markElevation]),
				},
				Igs: igs,
			})
		}

		*m = MarkList(marks)
	}
	return nil
}

func LoadMarks(path string) ([]Mark, error) {
	var m []Mark

	if err := LoadList(path, (*MarkList)(&m)); err != nil {
		return nil, err
	}

	sort.Sort(MarkList(m))

	return m, nil
}
