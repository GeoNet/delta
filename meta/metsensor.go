package meta

import (
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	installedMetSensorMake = iota
	installedMetSensorModel
	installedMetSensorSerial
	installedMetSensorMark
	installedMetSensorIMSComment
	installedMetSensorHumidityAccuracy
	installedMetSensorPressureAccuracy
	installedMetSensorTemperatureAccuracy
	installedMetSensorLatitude
	installedMetSensorLongitude
	installedMetSensorElevation
	installedMetSensorDatum
	installedMetSensorStart
	installedMetSensorStop
	installedMetSensorLast
)

var installedMetSensorHeaders Header = map[string]int{
	"Make":        installedMetSensorMake,
	"Model":       installedMetSensorModel,
	"Serial":      installedMetSensorSerial,
	"Mark":        installedMetSensorMark,
	"IMS Comment": installedMetSensorIMSComment,
	"Humidity":    installedMetSensorHumidityAccuracy,
	"Pressure":    installedMetSensorPressureAccuracy,
	"Temperature": installedMetSensorTemperatureAccuracy,
	"Latitude":    installedMetSensorLatitude,
	"Longitude":   installedMetSensorLongitude,
	"Elevation":   installedMetSensorElevation,
	"Datum":       installedMetSensorDatum,
	"Start Date":  installedMetSensorStart,
	"End Date":    installedMetSensorStop,
}

var InstalledMetSensorTable Table = Table{
	name:    "MetSensor",
	headers: installedMetSensorHeaders,
	primary: []string{"Make", "Model", "Serial", "Start Date"},
	native:  []string{"Latitude", "Longitude", "Elevation", "Humidity", "Pressure", "Temperature"},
	foreign: map[string]map[string]string{
		"Asset": {"Make": "Make", "Model": "Model", "Serial": "Serial"},
		"Mark":  {"Mark": "Mark"},
	},
	nullable: []string{"IMS Comment"},
	remap: map[string]string{
		"IMS Comment": "IMSComment",
		"Start Date":  "Start",
		"End Date":    "End",
	},
	start: "Start Date",
	end:   "End Date",
}

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
	Position

	Mark       string
	IMSComment string
	Accuracy   MetSensorAccuracy
}

type InstalledMetSensorList []InstalledMetSensor

func (ims InstalledMetSensorList) Len() int           { return len(ims) }
func (ims InstalledMetSensorList) Swap(i, j int)      { ims[i], ims[j] = ims[j], ims[i] }
func (ims InstalledMetSensorList) Less(i, j int) bool { return ims[i].Install.Less(ims[j].Install) }

func (ims InstalledMetSensorList) encode() [][]string {
	var data [][]string

	data = append(data, installedMetSensorHeaders.Columns())

	for _, row := range ims {
		data = append(data, []string{
			strings.TrimSpace(row.Make),
			strings.TrimSpace(row.Model),
			strings.TrimSpace(row.Serial),
			strings.TrimSpace(row.Mark),
			strings.TrimSpace(row.IMSComment),
			strings.TrimSpace(row.Accuracy.humidity),
			strings.TrimSpace(row.Accuracy.pressure),
			strings.TrimSpace(row.Accuracy.temperature),
			strings.TrimSpace(row.latitude),
			strings.TrimSpace(row.longitude),
			strings.TrimSpace(row.elevation),
			strings.TrimSpace(row.Datum),
			row.Start.Format(DateTimeFormat),
			row.End.Format(DateTimeFormat),
		})
	}

	return data
}

func (ims *InstalledMetSensorList) decode(data [][]string) error {
	if !(len(data) > 1) {
		return nil
	}

	var installedMetSensors []InstalledMetSensor

	fields := installedMetSensorHeaders.Fields(data[0])
	for _, row := range data[1:] {
		d := fields.Remap(row)

		h, err := strconv.ParseFloat(d[installedMetSensorHumidityAccuracy], 64)
		if err != nil {
			return err
		}
		p, err := strconv.ParseFloat(d[installedMetSensorPressureAccuracy], 64)
		if err != nil {
			return err
		}
		t, err := strconv.ParseFloat(d[installedMetSensorTemperatureAccuracy], 64)
		if err != nil {
			return err
		}

		lat, err := strconv.ParseFloat(d[installedMetSensorLatitude], 64)
		if err != nil {
			return err
		}
		lon, err := strconv.ParseFloat(d[installedMetSensorLongitude], 64)
		if err != nil {
			return err
		}
		elev, err := strconv.ParseFloat(d[installedMetSensorElevation], 64)
		if err != nil {
			return err
		}

		start, err := time.Parse(DateTimeFormat, d[installedMetSensorStart])
		if err != nil {
			return err
		}
		end, err := time.Parse(DateTimeFormat, d[installedMetSensorStop])
		if err != nil {
			return err
		}

		installedMetSensors = append(installedMetSensors, InstalledMetSensor{
			Install: Install{
				Equipment: Equipment{
					Make:   strings.TrimSpace(d[installedMetSensorMake]),
					Model:  strings.TrimSpace(d[installedMetSensorModel]),
					Serial: strings.TrimSpace(d[installedMetSensorSerial]),
				},
				Span: Span{
					Start: start,
					End:   end,
				},
			},
			Position: Position{
				Latitude:  lat,
				Longitude: lon,
				Elevation: elev,
				Datum:     strings.TrimSpace(d[installedMetSensorDatum]),

				latitude:  strings.TrimSpace(d[installedMetSensorLatitude]),
				longitude: strings.TrimSpace(d[installedMetSensorLongitude]),
				elevation: strings.TrimSpace(d[installedMetSensorElevation]),
			},
			Mark:       strings.TrimSpace(d[installedMetSensorMark]),
			IMSComment: strings.TrimSpace(d[installedMetSensorIMSComment]),
			Accuracy: MetSensorAccuracy{
				Humidity:    h,
				Pressure:    p,
				Temperature: t,

				humidity:    strings.TrimSpace(d[installedMetSensorHumidityAccuracy]),
				pressure:    strings.TrimSpace(d[installedMetSensorPressureAccuracy]),
				temperature: strings.TrimSpace(d[installedMetSensorTemperatureAccuracy]),
			},
		})
	}

	*ims = InstalledMetSensorList(installedMetSensors)

	return nil
}

func LoadInstalledMetSensors(path string) ([]InstalledMetSensor, error) {
	var ims []InstalledMetSensor

	if err := LoadList(path, (*InstalledMetSensorList)(&ims)); err != nil {
		return nil, err
	}

	sort.Sort(InstalledMetSensorList(ims))

	return ims, nil
}
