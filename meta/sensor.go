package meta

import (
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	installedSensorMake = iota
	installedSensorModel
	installedSensorSerial
	installedSensorStation
	installedSensorLocation
	installedSensorAzimuth
	installedSensorMethod
	installedSensorDip
	installedSensorDepth
	installedSensorNorth
	installedSensorEast
	installedSensorScaleFactor
	installedSensorScaleBias
	installedSensorStart
	installedSensorEnd
	installedSensorLast
)

var installedSensorHeaders Header = map[string]int{
	"Make":         installedSensorMake,
	"Model":        installedSensorModel,
	"Serial":       installedSensorSerial,
	"Station":      installedSensorStation,
	"Location":     installedSensorLocation,
	"Azimuth":      installedSensorAzimuth,
	"Method":       installedSensorMethod,
	"Dip":          installedSensorDip,
	"Depth":        installedSensorDepth,
	"North":        installedSensorNorth,
	"East":         installedSensorEast,
	"Scale Factor": installedSensorScaleFactor,
	"Scale Bias":   installedSensorScaleBias,
	"Start Date":   installedSensorStart,
	"End Date":     installedSensorEnd,
}

var InstalledSensorTable Table = Table{
	name:    "Sensor",
	headers: installedSensorHeaders,
	primary: []string{"Make", "Model", "Serial", "Station", "Location", "Start Date"},
	native:  []string{"Azimuth", "Dip", "Depth", "North", "East", "Scale Factor", "Scale Bias"},
	foreign: map[string][]string{
		//            "Asset": {"Make", "Model", "Serial"},
		"Site": {"Station", "Location"},
	},
	remap: map[string]string{
		"Scale Factor": "Factor",
		"Scale Bias":   "Bias",
		"Start Date":   "Start",
		"End Date":     "End",
	},
	start: "Start Date",
	end:   "End Date",
}

type InstalledSensor struct {
	Install
	Orientation
	Offset
	Scale

	Station  string `json:"station"`
	Location string `json:"location"`
}

type InstalledSensorList []InstalledSensor

func (is InstalledSensorList) Len() int           { return len(is) }
func (is InstalledSensorList) Swap(i, j int)      { is[i], is[j] = is[j], is[i] }
func (is InstalledSensorList) Less(i, j int) bool { return is[i].Install.Less(is[j].Install) }

func (is InstalledSensorList) encode() [][]string {
	var data [][]string

	data = append(data, installedSensorHeaders.Columns())

	for _, row := range is {
		data = append(data, []string{
			strings.TrimSpace(row.Make),
			strings.TrimSpace(row.Model),
			strings.TrimSpace(row.Serial),
			strings.TrimSpace(row.Station),
			strings.TrimSpace(row.Location),
			strings.TrimSpace(row.azimuth),
			strings.TrimSpace(row.Method),
			strings.TrimSpace(row.dip),
			strings.TrimSpace(row.vertical),
			strings.TrimSpace(row.north),
			strings.TrimSpace(row.east),
			strings.TrimSpace(row.factor),
			strings.TrimSpace(row.bias),
			row.Start.Format(DateTimeFormat),
			row.End.Format(DateTimeFormat),
		})
	}

	return data
}

func (is *InstalledSensorList) decode(data [][]string) error {
	if !(len(data) > 1) {
		return nil
	}

	var installedSensors []InstalledSensor

	fields := installedSensorHeaders.Fields(data[0])
	for _, row := range data[1:] {
		d := fields.Remap(row)

		azimuth, err := strconv.ParseFloat(d[installedSensorAzimuth], 64)
		if err != nil {
			return err
		}
		dip, err := strconv.ParseFloat(d[installedSensorDip], 64)
		if err != nil {
			return err
		}

		depth, err := strconv.ParseFloat(d[installedSensorDepth], 64)
		if err != nil {
			return err
		}
		north, err := strconv.ParseFloat(d[installedSensorNorth], 64)
		if err != nil {
			return err
		}
		east, err := strconv.ParseFloat(d[installedSensorEast], 64)
		if err != nil {
			return err
		}

		factor, err := strconv.ParseFloat(d[installedSensorScaleFactor], 64)
		if err != nil {
			return err
		}
		bias, err := strconv.ParseFloat(d[installedSensorScaleBias], 64)
		if err != nil {
			return err
		}

		start, err := time.Parse(DateTimeFormat, d[installedSensorStart])
		if err != nil {
			return err
		}
		end, err := time.Parse(DateTimeFormat, d[installedSensorEnd])
		if err != nil {
			return err
		}

		installedSensors = append(installedSensors, InstalledSensor{
			Install: Install{
				Equipment: Equipment{
					Make:   strings.TrimSpace(d[installedSensorMake]),
					Model:  strings.TrimSpace(d[installedSensorModel]),
					Serial: strings.TrimSpace(d[installedSensorSerial]),
				},
				Span: Span{
					Start: start,
					End:   end,
				},
			},
			Orientation: Orientation{
				Azimuth: azimuth,
				Dip:     dip,
				Method:  strings.TrimSpace(d[installedSensorMethod]),

				azimuth: strings.TrimSpace(d[installedSensorAzimuth]),
				dip:     strings.TrimSpace(d[installedSensorDip]),
			},
			Offset: Offset{
				Vertical: -depth,
				North:    north,
				East:     east,

				vertical: strings.TrimSpace(d[installedSensorDepth]),
				north:    strings.TrimSpace(d[installedSensorNorth]),
				east:     strings.TrimSpace(d[installedSensorEast]),
			},
			Scale: Scale{
				Factor: factor,
				Bias:   bias,

				factor: strings.TrimSpace(d[installedSensorScaleFactor]),
				bias:   strings.TrimSpace(d[installedSensorScaleBias]),
			},
			Station:  strings.TrimSpace(d[installedSensorStation]),
			Location: strings.TrimSpace(d[installedSensorLocation]),
		})
	}

	*is = InstalledSensorList(installedSensors)

	return nil
}

func LoadInstalledSensors(path string) ([]InstalledSensor, error) {
	var is []InstalledSensor

	if err := LoadList(path, (*InstalledSensorList)(&is)); err != nil {
		return nil, err
	}

	sort.Sort(InstalledSensorList(is))

	return is, nil
}
