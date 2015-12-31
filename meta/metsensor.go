package meta

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

type InstalledMetSensor struct {
	Install
	Point

	MarkCode string
	Comment  string
}

type InstalledMetSensors []InstalledMetSensor

func (m InstalledMetSensors) Len() int           { return len(m) }
func (m InstalledMetSensors) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }
func (m InstalledMetSensors) Less(i, j int) bool { return m[i].Install.less(m[j].Install) }

func (m InstalledMetSensors) encode() [][]string {
	data := [][]string{{
		"Met Sensor Make",
		"Met Sensor Model",
		"Serial Number",
		"Mark",
		"IMS Comment",
		"Latitude",
		"Longitude",
		"Height",
		"Datum",
		"Installation Date",
		"Removal Date",
	}}
	for _, v := range m {
		data = append(data, []string{
			strings.TrimSpace(v.Make),
			strings.TrimSpace(v.Model),
			strings.TrimSpace(v.Serial),
			strings.TrimSpace(v.MarkCode),
			strings.TrimSpace(v.Comment),
			strconv.FormatFloat(v.Latitude, 'g', -1, 64),
			strconv.FormatFloat(v.Longitude, 'g', -1, 64),
			strconv.FormatFloat(v.Elevation, 'g', -1, 64),
			strings.TrimSpace(v.Datum),
			v.Start.Format(DateTimeFormat),
			v.End.Format(DateTimeFormat),
		})
	}
	return data
}

func (m *InstalledMetSensors) decode(data [][]string) error {
	var metsensors []InstalledMetSensor
	if len(data) > 1 {
		for _, d := range data[1:] {
			if len(d) != 11 {
				return fmt.Errorf("incorrect number of installed metsensor fields")
			}
			var err error

			var lat, lon, elev float64
			if lat, err = strconv.ParseFloat(d[5], 64); err != nil {
				return err
			}
			if lon, err = strconv.ParseFloat(d[6], 64); err != nil {
				return err
			}
			if elev, err = strconv.ParseFloat(d[7], 64); err != nil {
				return err
			}

			var start, end time.Time
			if start, err = time.Parse(DateTimeFormat, d[9]); err != nil {
				return err
			}
			if end, err = time.Parse(DateTimeFormat, d[10]); err != nil {
				return err
			}

			metsensors = append(metsensors, InstalledMetSensor{
				Install: Install{
					Equipment: Equipment{
						Make:   strings.TrimSpace(d[0]),
						Model:  strings.TrimSpace(d[1]),
						Serial: strings.TrimSpace(d[2]),
					},
					Span: Span{
						Start: start,
						End:   end,
					},
				},
				Point: Point{
					Latitude:  lat,
					Longitude: lon,
					Elevation: elev,
					Datum:     strings.TrimSpace(d[8]),
				},
				MarkCode: strings.TrimSpace(d[3]),
				Comment:  strings.TrimSpace(d[4]),
			})
		}

		*m = InstalledMetSensors(metsensors)
	}
	return nil
}

func LoadInstalledMetSensors(path string) ([]InstalledMetSensor, error) {
	var m []InstalledMetSensor

	if err := LoadList(path, (*InstalledMetSensors)(&m)); err != nil {
		return nil, err
	}

	sort.Sort(InstalledMetSensors(m))

	return m, nil
}
