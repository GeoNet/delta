package meta

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Mark struct {
	Reference
	Point
	Span

	MarkType           string
	MonumentType       string
	DomeNumber         string
	GroundRelationship float64
}

type Marks []Mark

func (m Marks) Len() int           { return len(m) }
func (m Marks) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }
func (m Marks) Less(i, j int) bool { return m[i].Reference.less(m[j].Reference) }

func (m Marks) encode() [][]string {
	data := [][]string{{
		"Mark Code",
		"Network Code",
		"Mark Name",
		"Latitude",
		"Longitude",
		"Height",
		"Datum",
		"Ground Relationship",
		"Mark Type",
		"Monument Type",
		"Dome Number",
		"Start Time",
		"End Time",
	}}
	for _, v := range m {
		data = append(data, []string{
			strings.TrimSpace(v.Code),
			strings.TrimSpace(v.Network),
			strings.TrimSpace(v.Name),
			strconv.FormatFloat(v.Latitude, 'g', -1, 64),
			strconv.FormatFloat(v.Longitude, 'g', -1, 64),
			strconv.FormatFloat(v.Elevation, 'g', -1, 64),
			strings.TrimSpace(v.Datum),
			strconv.FormatFloat(v.GroundRelationship, 'g', -1, 64),
			strings.TrimSpace(v.MarkType),
			strings.TrimSpace(v.MonumentType),
			strings.TrimSpace(v.DomeNumber),
			v.Start.Format(DateTimeFormat),
			v.End.Format(DateTimeFormat),
		})
	}
	return data
}

func (m *Marks) decode(data [][]string) error {
	var marks []Mark
	if len(data) > 1 {
		for _, d := range data[1:] {
			if len(d) != 13 {
				return fmt.Errorf("incorrect number of installed mark fields")
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

			var ground float64
			if ground, err = strconv.ParseFloat(d[7], 64); err != nil {
				return err
			}

			var start, end time.Time
			if start, err = time.Parse(DateTimeFormat, d[11]); err != nil {
				return err
			}
			if end, err = time.Parse(DateTimeFormat, d[12]); err != nil {
				return err
			}

			marks = append(marks, Mark{
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
				GroundRelationship: ground,
				MarkType:           strings.TrimSpace(d[8]),
				MonumentType:       strings.TrimSpace(d[9]),
				DomeNumber:         strings.TrimSpace(d[10]),
			})
		}

		*m = Marks(marks)
	}
	return nil
}

func LoadMarks(path string) ([]Mark, error) {
	var m []Mark

	if err := LoadList(path, (*Marks)(&m)); err != nil {
		return nil, err
	}

	sort.Sort(Marks(m))

	return m, nil
}
