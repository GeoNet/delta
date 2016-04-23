package meta

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	recorderMake int = iota
	recorderSensorModel
	recorderDataloggerModel
	recorderSerialNumber
	recorderStationCode
	recorderLocationCode
	recorderInstallationAzimuth
	recorderInstallationDip
	recorderInstallationDepth
	recorderInstallationDate
	recorderRemovalDate
	recorderLast
)

type InstalledRecorder struct {
	InstalledSensor

	DataloggerModel string
}

type InstalledRecorderList []InstalledRecorder

func (r InstalledRecorderList) Len() int           { return len(r) }
func (r InstalledRecorderList) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r InstalledRecorderList) Less(i, j int) bool { return r[i].Install.less(r[j].Install) }

func (r InstalledRecorderList) encode() [][]string {
	data := [][]string{{
		"Recorder Make",
		"Sensor Model",
		"Datalogger Model",
		"Serial Number",
		"Station Code",
		"Location Code",
		"Installation Azimuth",
		"Installation Dip",
		"Installation Depth",
		"Installation Date",
		"Removal Date",
	}}

	for _, v := range r {
		data = append(data, []string{
			strings.TrimSpace(v.Make),
			strings.TrimSpace(v.Model),
			strings.TrimSpace(v.DataloggerModel),
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
func (r *InstalledRecorderList) decode(data [][]string) error {
	var recorders []InstalledRecorder
	if len(data) > 1 {
		for _, d := range data[1:] {
			if len(d) != recorderLast {
				return fmt.Errorf("incorrect number of installed recorder fields")
			}
			var err error

			var azimuth, dip, depth float64
			if azimuth, err = strconv.ParseFloat(d[recorderInstallationAzimuth], 64); err != nil {
				return err
			}
			if dip, err = strconv.ParseFloat(d[recorderInstallationDip], 64); err != nil {
				return err
			}
			if depth, err = strconv.ParseFloat(d[recorderInstallationDepth], 64); err != nil {
				return err
			}

			var start, end time.Time
			if start, err = time.Parse(DateTimeFormat, d[recorderInstallationDate]); err != nil {
				return err
			}
			if end, err = time.Parse(DateTimeFormat, d[recorderRemovalDate]); err != nil {
				return err
			}

			recorders = append(recorders, InstalledRecorder{
				InstalledSensor: InstalledSensor{
					Install: Install{
						Equipment: Equipment{
							Make:   strings.TrimSpace(d[recorderMake]),
							Model:  strings.TrimSpace(d[recorderSensorModel]),
							Serial: strings.TrimSpace(d[recorderSerialNumber]),
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
					StationCode:  strings.TrimSpace(d[recorderStationCode]),
					LocationCode: strings.TrimSpace(d[recorderLocationCode]),
				},
				DataloggerModel: strings.TrimSpace(d[recorderDataloggerModel]),
			})
		}

		*r = InstalledRecorderList(recorders)
	}
	return nil
}

func LoadInstalledRecorders(path string) ([]InstalledRecorder, error) {
	var r []InstalledRecorder

	if err := LoadList(path, (*InstalledRecorderList)(&r)); err != nil {
		return nil, err
	}

	sort.Sort(InstalledRecorderList(r))

	return r, nil
}
