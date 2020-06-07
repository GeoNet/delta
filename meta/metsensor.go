package meta

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	installedMetsensorMake = iota
	installedMetsensorModel
	installedMetsensorSerial
	installedMetsensorMark
	installedMetsensorIMSComment
	installedMetsensorHumidityAccuracy
	installedMetsensorPressureAccuracy
	installedMetsensorTemperatureAccuracy
	installedMetsensorLatitude
	installedMetsensorLongitude
	installedMetsensorElevation
	installedMetsensorDatum
	installedMetsensorStart
	installedMetsensorStop
	installedMetsensorLast
)

type MetSensorAccuracy struct {
	Humidity    float64
	Pressure    float64
	Temperature float64

	humidity    string // shadow variable to maintain formatting
	pressure    string // shadow variable to maintain formatting
	temperature string // shadow variable to maintain formatting
}

type InstalledMetSensor struct {
	Install
	Point

	Mark       string
	IMSComment string
	Accuracy   MetSensorAccuracy
}

type InstalledMetSensorList []InstalledMetSensor

func (m InstalledMetSensorList) Len() int           { return len(m) }
func (m InstalledMetSensorList) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }
func (m InstalledMetSensorList) Less(i, j int) bool { return m[i].Install.Less(m[j].Install) }

func (m InstalledMetSensorList) encode() [][]string {
	data := [][]string{{
		"Make",
		"Model",
		"Serial",
		"Mark",
		"IMS Comment",
		"Humidity",
		"Pressure",
		"Temperature",
		"Latitude",
		"Longitude",
		"Elevation",
		"Datum",
		"Start Date",
		"End Date",
	}}
	for _, v := range m {
		data = append(data, []string{
			strings.TrimSpace(v.Make),
			strings.TrimSpace(v.Model),
			strings.TrimSpace(v.Serial),
			strings.TrimSpace(v.Mark),
			strings.TrimSpace(v.IMSComment),
			strings.TrimSpace(v.Accuracy.humidity),
			strings.TrimSpace(v.Accuracy.pressure),
			strings.TrimSpace(v.Accuracy.temperature),
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

func (m *InstalledMetSensorList) decode(data [][]string) error {
	var installedMetsensors []InstalledMetSensor
	if len(data) > 1 {
		for _, d := range data[1:] {
			if len(d) != installedMetsensorLast {
				return fmt.Errorf("incorrect number of installed metsensor fields")
			}
			var err error

			var h, p, t float64
			if h, err = strconv.ParseFloat(d[installedMetsensorHumidityAccuracy], 64); err != nil {
				return err
			}
			if p, err = strconv.ParseFloat(d[installedMetsensorPressureAccuracy], 64); err != nil {
				return err
			}
			if t, err = strconv.ParseFloat(d[installedMetsensorTemperatureAccuracy], 64); err != nil {
				return err
			}

			var lat, lon, elev float64
			if lat, err = strconv.ParseFloat(d[installedMetsensorLatitude], 64); err != nil {
				return err
			}
			if lon, err = strconv.ParseFloat(d[installedMetsensorLongitude], 64); err != nil {
				return err
			}
			if elev, err = strconv.ParseFloat(d[installedMetsensorElevation], 64); err != nil {
				return err
			}

			var start, end time.Time
			if start, err = time.Parse(DateTimeFormat, d[installedMetsensorStart]); err != nil {
				return err
			}
			if end, err = time.Parse(DateTimeFormat, d[installedMetsensorStop]); err != nil {
				return err
			}

			installedMetsensors = append(installedMetsensors, InstalledMetSensor{
				Install: Install{
					Equipment: Equipment{
						Make:   strings.TrimSpace(d[installedMetsensorMake]),
						Model:  strings.TrimSpace(d[installedMetsensorModel]),
						Serial: strings.TrimSpace(d[installedMetsensorSerial]),
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
					Datum:     strings.TrimSpace(d[installedMetsensorDatum]),

					latitude:  strings.TrimSpace(d[installedMetsensorLatitude]),
					longitude: strings.TrimSpace(d[installedMetsensorLongitude]),
					elevation: strings.TrimSpace(d[installedMetsensorElevation]),
				},
				Mark:       strings.TrimSpace(d[installedMetsensorMark]),
				IMSComment: strings.TrimSpace(d[installedMetsensorIMSComment]),
				Accuracy: MetSensorAccuracy{
					Humidity:    h,
					Pressure:    p,
					Temperature: t,

					humidity:    strings.TrimSpace(d[installedMetsensorHumidityAccuracy]),
					pressure:    strings.TrimSpace(d[installedMetsensorPressureAccuracy]),
					temperature: strings.TrimSpace(d[installedMetsensorTemperatureAccuracy]),
				},
			})
		}

		*m = InstalledMetSensorList(installedMetsensors)
	}
	return nil
}

func LoadInstalledMetSensors(path string) ([]InstalledMetSensor, error) {
	var m []InstalledMetSensor

	if err := LoadList(path, (*InstalledMetSensorList)(&m)); err != nil {
		return nil, err
	}

	sort.Sort(InstalledMetSensorList(m))

	return m, nil
}
