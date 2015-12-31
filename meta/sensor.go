package meta

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

type InstalledSensor struct {
	Install
	Orientation
	Offset

	StationCode  string `csv:"Station Code"`
	LocationCode string `csv:"Location Code"`
}

type InstalledSensors []InstalledSensor

func (s InstalledSensors) Len() int           { return len(s) }
func (s InstalledSensors) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s InstalledSensors) Less(i, j int) bool { return s[i].Install.less(s[j].Install) }

func (s InstalledSensors) encode() [][]string {
	data := [][]string{{
		"Sensor Make",
		"Sensor Model",
		"Serial Number",
		"Station Code",
		"Location Code",
		"Installation Azimuth",
		"Installation Dip",
		"Installation Depth",
		"Installation Date",
		"Removal Date",
	}}

	for _, v := range s {
		data = append(data, []string{
			strings.TrimSpace(v.Make),
			strings.TrimSpace(v.Model),
			strings.TrimSpace(v.Serial),
			strings.TrimSpace(v.StationCode),
			strings.TrimSpace(v.LocationCode),
			strconv.FormatFloat(v.Azimuth, 'g', -1, 64),
			strconv.FormatFloat(v.Dip, 'g', -1, 64),
			func() string {
				if v.Height == 0.0 {
					return strconv.FormatFloat(0.0, 'g', -1, 64)
				} else {
					return strconv.FormatFloat(-v.Height, 'g', -1, 64)
				}
			}(),
			v.Start.Format(DateTimeFormat),
			v.End.Format(DateTimeFormat),
		})
	}
	return data
}
func (s *InstalledSensors) decode(data [][]string) error {
	var sensors []InstalledSensor
	if len(data) > 1 {
		for _, d := range data[1:] {
			if len(d) != 10 {
				return fmt.Errorf("incorrect number of installed sensor fields")
			}
			var err error

			var azimuth, dip, depth float64
			if azimuth, err = strconv.ParseFloat(d[5], 64); err != nil {
				return err
			}
			if dip, err = strconv.ParseFloat(d[6], 64); err != nil {
				return err
			}
			if depth, err = strconv.ParseFloat(d[7], 64); err != nil {
				return err
			}

			var start, end time.Time
			if start, err = time.Parse(DateTimeFormat, d[8]); err != nil {
				return err
			}
			if end, err = time.Parse(DateTimeFormat, d[9]); err != nil {
				return err
			}

			sensors = append(sensors, InstalledSensor{
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
				Orientation: Orientation{
					Azimuth: azimuth,
					Dip:     dip,
				},
				Offset: Offset{
					Height: -depth,
				},
				StationCode:  strings.TrimSpace(d[3]),
				LocationCode: strings.TrimSpace(d[4]),
			})
		}

		*s = InstalledSensors(sensors)
	}
	return nil
}

func LoadInstalledSensors(path string) ([]InstalledSensor, error) {
	var s []InstalledSensor

	if err := LoadList(path, (*InstalledSensors)(&s)); err != nil {
		return nil, err
	}

	sort.Sort(InstalledSensors(s))

	return s, nil
}
