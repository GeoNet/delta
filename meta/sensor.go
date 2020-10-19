package meta

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	sensorMake = iota
	sensorModel
	sensorSerial
	sensorStation
	sensorLocation
	sensorAzimuth
	sensorMethod
	sensorDip
	sensorDepth
	sensorNorth
	sensorEast
	sensorScaleFactor
	sensorScaleBias
	sensorStart
	sensorEnd
	sensorLast
)

type InstalledSensor struct {
	Install
	Orientation
	Offset
	Scale

	Station  string
	Location string
}

type InstalledSensorList []InstalledSensor

func (s InstalledSensorList) Len() int           { return len(s) }
func (s InstalledSensorList) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s InstalledSensorList) Less(i, j int) bool { return s[i].Install.Less(s[j].Install) }

func (s InstalledSensorList) encode() [][]string {
	data := [][]string{{
		"Make",
		"Model",
		"Serial",
		"Station",
		"Location",
		"Azimuth",
		"Method",
		"Dip",
		"Depth",
		"North",
		"East",
		"Scale Factor",
		"Scale Bias",
		"Start Date",
		"End Date",
	}}

	for _, v := range s {
		data = append(data, []string{
			strings.TrimSpace(v.Make),
			strings.TrimSpace(v.Model),
			strings.TrimSpace(v.Serial),
			strings.TrimSpace(v.Station),
			strings.TrimSpace(v.Location),
			strings.TrimSpace(v.azimuth),
			strings.TrimSpace(v.Method),
			strings.TrimSpace(v.dip),
			strings.TrimSpace(v.vertical),
			strings.TrimSpace(v.north),
			strings.TrimSpace(v.east),
			strings.TrimSpace(v.factor),
			strings.TrimSpace(v.bias),
			v.Start.Format(DateTimeFormat),
			v.End.Format(DateTimeFormat),
		})
	}
	return data
}
func (s *InstalledSensorList) decode(data [][]string) error {
	var sensors []InstalledSensor
	if len(data) > 1 {
		for _, d := range data[1:] {
			if len(d) != sensorLast {
				return fmt.Errorf("incorrect number of installed sensor fields")
			}
			var err error

			var azimuth, dip float64
			if azimuth, err = strconv.ParseFloat(d[sensorAzimuth], 64); err != nil {
				return err
			}
			if dip, err = strconv.ParseFloat(d[sensorDip], 64); err != nil {
				return err
			}

			var depth, north, east float64
			if depth, err = strconv.ParseFloat(d[sensorDepth], 64); err != nil {
				return err
			}
			if north, err = strconv.ParseFloat(d[sensorNorth], 64); err != nil {
				return err
			}
			if east, err = strconv.ParseFloat(d[sensorEast], 64); err != nil {
				return err
			}

			var factor, bias float64
			if factor, err = strconv.ParseFloat(d[sensorScaleFactor], 64); err != nil {
				return err
			}
			if bias, err = strconv.ParseFloat(d[sensorScaleBias], 64); err != nil {
				return err
			}

			var start, end time.Time
			if start, err = time.Parse(DateTimeFormat, d[sensorStart]); err != nil {
				return err
			}
			if end, err = time.Parse(DateTimeFormat, d[sensorEnd]); err != nil {
				return err
			}

			sensors = append(sensors, InstalledSensor{
				Install: Install{
					Equipment: Equipment{
						Make:   strings.TrimSpace(d[sensorMake]),
						Model:  strings.TrimSpace(d[sensorModel]),
						Serial: strings.TrimSpace(d[sensorSerial]),
					},
					Span: Span{
						Start: start,
						End:   end,
					},
				},
				Orientation: Orientation{
					Azimuth: azimuth,
					Dip:     dip,
					Method:  strings.TrimSpace(d[sensorMethod]),

					azimuth: strings.TrimSpace(d[sensorAzimuth]),
					dip:     strings.TrimSpace(d[sensorDip]),
				},
				Offset: Offset{
					Vertical: -depth,
					North:    north,
					East:     east,

					vertical: strings.TrimSpace(d[sensorDepth]),
					north:    strings.TrimSpace(d[sensorNorth]),
					east:     strings.TrimSpace(d[sensorEast]),
				},
				Scale: Scale{
					Factor: factor,
					Bias:   bias,

					factor: strings.TrimSpace(d[sensorScaleFactor]),
					bias:   strings.TrimSpace(d[sensorScaleBias]),
				},
				Station:  strings.TrimSpace(d[sensorStation]),
				Location: strings.TrimSpace(d[sensorLocation]),
			})
		}

		*s = InstalledSensorList(sensors)
	}
	return nil
}

func LoadInstalledSensors(path string) ([]InstalledSensor, error) {
	var s []InstalledSensor

	if err := LoadList(path, (*InstalledSensorList)(&s)); err != nil {
		return nil, err
	}

	sort.Sort(InstalledSensorList(s))

	return s, nil
}
