package meta

import (
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

var mountHeaders Header = map[string]int{
	"Mount":       mountCode,
	"Network":     mountNetwork,
	"Name":        mountName,
	"Latitude":    mountLatitude,
	"Longitude":   mountLongitude,
	"Elevation":   mountElevation,
	"Datum":       mountDatum,
	"Description": mountDescription,
	"Start Date":  mountStart,
	"End Date":    mountEnd,
}

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

	var data [][]string

	data = append(data, mountHeaders.Columns())

	for _, row := range m {
		data = append(data, []string{
			strings.TrimSpace(row.Code),
			strings.TrimSpace(row.Network),
			strings.TrimSpace(row.Name),
			strings.TrimSpace(row.latitude),
			strings.TrimSpace(row.longitude),
			strings.TrimSpace(row.elevation),
			strings.TrimSpace(row.Datum),
			strings.TrimSpace(row.Description),
			row.Start.Format(DateTimeFormat),
			row.End.Format(DateTimeFormat),
		})
	}

	return data
}

func (m *MountList) decode(data [][]string) error {
	if !(len(data) > 1) {
		return nil
	}

	var mounts []Mount

	fields := mountHeaders.Fields(data[0])
	for _, row := range data[1:] {
		d := fields.Remap(row)

		lat, err := strconv.ParseFloat(d[mountLatitude], 64)
		if err != nil {
			return err
		}
		lon, err := strconv.ParseFloat(d[mountLongitude], 64)
		if err != nil {
			return err
		}
		elev, err := strconv.ParseFloat(d[mountElevation], 64)
		if err != nil {
			return err
		}

		start, err := time.Parse(DateTimeFormat, d[mountStart])
		if err != nil {
			return err
		}
		end, err := time.Parse(DateTimeFormat, d[mountEnd])
		if err != nil {
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
