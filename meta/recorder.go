package meta

import (
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	installedRecorderMake int = iota
	installedRecorderSensorModel
	installedRecorderDataloggerModel
	installedRecorderSerial
	installedRecorderStation
	installedRecorderLocation
	installedRecorderAzimuth
	installedRecorderMethod
	installedRecorderDip
	installedRecorderDepth
	installedRecorderStart
	installedRecorderEnd
	installedRecorderLast
)

var installedRecorderHeaders Header = map[string]int{
	"Make":       installedRecorderMake,
	"Sensor":     installedRecorderSensorModel,
	"Datalogger": installedRecorderDataloggerModel,
	"Serial":     installedRecorderSerial,
	"Station":    installedRecorderStation,
	"Location":   installedRecorderLocation,
	"Azimuth":    installedRecorderAzimuth,
	"Method":     installedRecorderMethod,
	"Dip":        installedRecorderDip,
	"Depth":      installedRecorderDepth,
	"Start Date": installedRecorderStart,
	"End Date":   installedRecorderEnd,
}

type InstalledRecorder struct {
	InstalledSensor

	DataloggerModel string
}

type InstalledRecorderList []InstalledRecorder

func (ir InstalledRecorderList) Len() int           { return len(ir) }
func (ir InstalledRecorderList) Swap(i, j int)      { ir[i], ir[j] = ir[j], ir[i] }
func (ir InstalledRecorderList) Less(i, j int) bool { return ir[i].Install.Less(ir[j].Install) }

func (ir InstalledRecorderList) encode() [][]string {
	var data [][]string

	data = append(data, installedRecorderHeaders.Columns())

	for _, row := range ir {
		data = append(data, []string{
			strings.TrimSpace(row.Make),
			strings.TrimSpace(row.Model),
			strings.TrimSpace(row.DataloggerModel),
			strings.TrimSpace(row.Serial),
			strings.TrimSpace(row.Station),
			strings.TrimSpace(row.Location),
			strings.TrimSpace(row.azimuth),
			strings.TrimSpace(row.Method),
			strings.TrimSpace(row.dip),
			strings.TrimSpace(row.vertical),
			row.Start.Format(DateTimeFormat),
			row.End.Format(DateTimeFormat),
		})
	}

	return data
}

func (ir *InstalledRecorderList) decode(data [][]string) error {
	if !(len(data) > 1) {
		return nil
	}

	var installedRecorders []InstalledRecorder

	fields := installedRecorderHeaders.Fields(data[0])
	for _, row := range data[1:] {
		d := fields.Remap(row)

		azimuth, err := strconv.ParseFloat(d[installedRecorderAzimuth], 64)
		if err != nil {
			return err
		}
		dip, err := strconv.ParseFloat(d[installedRecorderDip], 64)
		if err != nil {
			return err
		}
		depth, err := strconv.ParseFloat(d[installedRecorderDepth], 64)
		if err != nil {
			return err
		}

		start, err := time.Parse(DateTimeFormat, d[installedRecorderStart])
		if err != nil {
			return err
		}
		end, err := time.Parse(DateTimeFormat, d[installedRecorderEnd])
		if err != nil {
			return err
		}

		installedRecorders = append(installedRecorders, InstalledRecorder{
			InstalledSensor: InstalledSensor{
				Install: Install{
					Equipment: Equipment{
						Make:   strings.TrimSpace(d[installedRecorderMake]),
						Model:  strings.TrimSpace(d[installedRecorderSensorModel]),
						Serial: strings.TrimSpace(d[installedRecorderSerial]),
					},
					Span: Span{
						Start: start,
						End:   end,
					},
				},
				Orientation: Orientation{
					Azimuth: azimuth,
					Dip:     dip,
					Method:  strings.TrimSpace(d[installedRecorderMethod]),

					azimuth: strings.TrimSpace(d[installedRecorderAzimuth]),
					dip:     strings.TrimSpace(d[installedRecorderDip]),
				},
				Offset: Offset{
					Vertical: -depth,

					vertical: strings.TrimSpace(d[installedRecorderDepth]),
				},
				Station:  strings.TrimSpace(d[installedRecorderStation]),
				Location: strings.TrimSpace(d[installedRecorderLocation]),
			},
			DataloggerModel: strings.TrimSpace(d[installedRecorderDataloggerModel]),
		})
	}

	*ir = InstalledRecorderList(installedRecorders)

	return nil
}

func LoadInstalledRecorders(path string) ([]InstalledRecorder, error) {
	var ir []InstalledRecorder

	if err := LoadList(path, (*InstalledRecorderList)(&ir)); err != nil {
		return nil, err
	}

	sort.Sort(InstalledRecorderList(ir))

	return ir, nil
}
