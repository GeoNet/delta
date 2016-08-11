package meta

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	sensorSensorMake = iota
	sensorSensorModel
	sensorSerialNumber
	sensorStationCode
	sensorLocationCode
	sensorInstallationAzimuth
	sensorInstallationDip
	sensorInstallationDepth
	sensorScaleFactor
	sensorScaleBias
	sensorInstallationDate
	sensorRemovalDate
	sensorLast
)

type InstalledSensor struct {
	Install
	Orientation
	Offset
	Scale

	StationCode  string
	LocationCode string
}

type InstalledSensorList []InstalledSensor

func (s InstalledSensorList) Len() int           { return len(s) }
func (s InstalledSensorList) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s InstalledSensorList) Less(i, j int) bool { return s[i].Install.less(s[j].Install) }

func (s InstalledSensorList) encode() [][]string {
	data := [][]string{{
		"Sensor Make",
		"Sensor Model",
		"Serial Number",
		"Station Code",
		"Location Code",
		"Installation Azimuth",
		"Installation Dip",
		"Installation Depth",
		"Scale Factor",
		"Scale Bias",
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
			strconv.FormatFloat(v.Factor, 'g', -1, 64),
			strconv.FormatFloat(v.Bias, 'g', -1, 64),
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

			var azimuth, dip, depth float64
			if azimuth, err = strconv.ParseFloat(d[sensorInstallationAzimuth], 64); err != nil {
				return err
			}
			if dip, err = strconv.ParseFloat(d[sensorInstallationDip], 64); err != nil {
				return err
			}
			if depth, err = strconv.ParseFloat(d[sensorInstallationDepth], 64); err != nil {
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
			if start, err = time.Parse(DateTimeFormat, d[sensorInstallationDate]); err != nil {
				return err
			}
			if end, err = time.Parse(DateTimeFormat, d[sensorRemovalDate]); err != nil {
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
				Scale: Scale{
					Factor: factor,
					Bias:   bias,
				},
				StationCode:  strings.TrimSpace(d[3]),
				LocationCode: strings.TrimSpace(d[4]),
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
