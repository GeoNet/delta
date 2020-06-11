package meta

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	mountCode = iota
	mountNetwork
	mountName
	mountLatitude
	mountLongitude
	mountElevation
	mountDatum
	mountDescription
	mountStart
	mountEnd
	mountLast
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
		"Mount",
		"Network",
		"Name",
		"Latitude",
		"Longitude",
		"Elevation",
		"Datum",
		"Description",
		"Start Date",
		"End Date",
	}}
	for _, v := range m {
		data = append(data, []string{
			strings.TrimSpace(v.Code),
			strings.TrimSpace(v.Network),
			strings.TrimSpace(v.Name),
			strings.TrimSpace(v.latitude),
			strings.TrimSpace(v.longitude),
			strings.TrimSpace(v.elevation),
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
			if len(d) != mountLast {
				return fmt.Errorf("incorrect number of installed mount fields")
			}
			var err error

			var lat, lon, elev float64
			if lat, err = strconv.ParseFloat(d[mountLatitude], 64); err != nil {
				return err
			}
			if lon, err = strconv.ParseFloat(d[mountLongitude], 64); err != nil {
				return err
			}
			if elev, err = strconv.ParseFloat(d[mountElevation], 64); err != nil {
				return err
			}

			var start, end time.Time
			if start, err = time.Parse(DateTimeFormat, d[mountStart]); err != nil {
				return err
			}
			if end, err = time.Parse(DateTimeFormat, d[mountEnd]); err != nil {
				return err
			}

			mounts = append(mounts, Mount{
				Reference: Reference{
					Code:    strings.TrimSpace(d[mountCode]),
					Network: strings.TrimSpace(d[mountNetwork]),
					Name:    strings.TrimSpace(d[mountName]),
				},
				Point: Point{
					Latitude:  lat,
					Longitude: lon,
					Elevation: elev,
					Datum:     strings.TrimSpace(d[mountDatum]),

					latitude:  strings.TrimSpace(d[mountLatitude]),
					longitude: strings.TrimSpace(d[mountLongitude]),
					elevation: strings.TrimSpace(d[mountElevation]),
				},
				Span: Span{
					Start: start,
					End:   end,
				},
				Description: strings.TrimSpace(d[mountDescription]),
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
