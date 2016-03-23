package meta

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Mount struct {
	Reference
	Point
	Span

	Description string
}

type MountList []Mount

func (m MountList) Len() int           { return len(m) }
func (m MountList) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }
func (m MountList) Less(i, j int) bool { return m[i].Code < m[j].Code }

func (m MountList) encode() [][]string {
	data := [][]string{{
		"Mount Code",
		"Mount Network",
		"Mount Name",
		"Latitude",
		"Longitude",
		"Height",
		"Datum",
		"Description",
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
			strings.TrimSpace(v.Description),
			v.Start.Format(DateTimeFormat),
			v.End.Format(DateTimeFormat),
		})
	}
	return data
}

func (m *MountList) decode(data [][]string) error {
	var mounts []Mount
	if len(data) > 1 {
		for _, d := range data[1:] {
			if len(d) != 10 {
				return fmt.Errorf("incorrect number of installed mount fields")
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
			if start, err = time.Parse(DateTimeFormat, d[8]); err != nil {
				return err
			}
			if end, err = time.Parse(DateTimeFormat, d[9]); err != nil {
				return err
			}

			mounts = append(mounts, Mount{
				Reference: Reference{
					Code:    strings.TrimSpace(d[0]),
					Network: strings.TrimSpace(d[1]),
					Name:    strings.TrimSpace(d[2]),
				},
				Point: Point{
					Latitude:  lat,
					Longitude: lon,
					Elevation: elev,
					Datum:     strings.TrimSpace(d[6]),
				},
				Span: Span{
					Start: start,
					End:   end,
				},
				Description: strings.TrimSpace(d[7]),
			})
		}

		*m = MountList(mounts)
	}
	return nil
}

func LoadMounts(path string) ([]Mount, error) {
	var m []Mount

	if err := LoadList(path, (*MountList)(&m)); err != nil {
		return nil, err
	}

	sort.Sort(MountList(m))

	return m, nil
}
