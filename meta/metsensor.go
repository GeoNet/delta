package meta

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	metsensorMetSensorMake = iota
	metsensorMetSensorModel
	metsensorSerialNumber
	metsensorMarkCode
	metsensorIMSComment
	metsensorLatitude
	metsensorLongitude
	metsensorElevation
	metsensorDatum
	metsensorInstallationDate
	metsensorRemovalDate
	metsensorLast
)

type InstalledMetSensor struct {
	Install
	Point

	MarkCode   string
	IMSComment string
}

type InstalledMetSensorList []InstalledMetSensor

func (m InstalledMetSensorList) Len() int           { return len(m) }
func (m InstalledMetSensorList) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }
func (m InstalledMetSensorList) Less(i, j int) bool { return m[i].Install.less(m[j].Install) }

func (m InstalledMetSensorList) encode() [][]string {
	data := [][]string{{
		"Met Sensor Make",
		"Met Sensor Model",
		"Serial Number",
		"Mark",
		"IMS Comment",
		"Latitude",
		"Longitude",
		"Elevation",
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
			strings.TrimSpace(v.IMSComment),
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

func (m *InstalledMetSensorList) decode(data [][]string) error {
	var metsensors []InstalledMetSensor
	if len(data) > 1 {
		for _, d := range data[1:] {
			if len(d) != metsensorLast {
				return fmt.Errorf("incorrect number of installed metsensor fields")
			}
			var err error

			var lat, lon, elev float64
			if lat, err = strconv.ParseFloat(d[metsensorLatitude], 64); err != nil {
				return err
			}
			if lon, err = strconv.ParseFloat(d[metsensorLongitude], 64); err != nil {
				return err
			}
			if elev, err = strconv.ParseFloat(d[metsensorElevation], 64); err != nil {
				return err
			}

			var start, end time.Time
			if start, err = time.Parse(DateTimeFormat, d[metsensorInstallationDate]); err != nil {
				return err
			}
			if end, err = time.Parse(DateTimeFormat, d[metsensorRemovalDate]); err != nil {
				return err
			}

			metsensors = append(metsensors, InstalledMetSensor{
				Install: Install{
					Equipment: Equipment{
						Make:   strings.TrimSpace(d[metsensorMetSensorMake]),
						Model:  strings.TrimSpace(d[metsensorMetSensorModel]),
						Serial: strings.TrimSpace(d[metsensorSerialNumber]),
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
					Datum:     strings.TrimSpace(d[metsensorDatum]),
				},
				MarkCode:   strings.TrimSpace(d[metsensorMarkCode]),
				IMSComment: strings.TrimSpace(d[metsensorIMSComment]),
			})
		}

		*m = InstalledMetSensorList(metsensors)
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
