package meta

import (
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	installedCameraMake int = iota
	installedCameraModel
	installedCameraSerial
	installedCameraMount
	installedCameraView
	installedCameraDip
	installedCameraAzimuth
	installedCameraHeight
	installedCameraNorth
	installedCameraEast
	installedCameraStart
	installedCameraEnd
	installedCameraNotes
	installedCameraLast
)

var installedCameraHeaders Header = map[string]int{
	"Make":       installedCameraMake,
	"Model":      installedCameraModel,
	"Serial":     installedCameraSerial,
	"Mount":      installedCameraMount,
	"View":       installedCameraView,
	"Dip":        installedCameraDip,
	"Azimuth":    installedCameraAzimuth,
	"Height":     installedCameraHeight,
	"North":      installedCameraNorth,
	"East":       installedCameraEast,
	"Start Date": installedCameraStart,
	"End Date":   installedCameraEnd,
	"Notes":      installedCameraNotes,
}

var InstalledCameraTable Table = Table{
	name:    "Camera",
	headers: installedCameraHeaders,
	primary: []string{"Make", "Model", "Serial", "Start Date"},
	native:  []string{"Azimuth", "Dip", "Height", "North", "East"},
	foreign: map[string][]string{
		"Asset": {"Make", "Model", "Serial"},
		"Mount": {"Mount"},
	},
	remap: map[string]string{
		"Start Date": "Start",
		"End Date":   "End",
	},
	start: "Start Date",
	end:   "End Date",
}

type InstalledCamera struct {
	Install
	Orientation
	Offset

	Mount string
	View  string
	Notes string
}

type InstalledCameraList []InstalledCamera

func (ic InstalledCameraList) Len() int           { return len(ic) }
func (ic InstalledCameraList) Swap(i, j int)      { ic[i], ic[j] = ic[j], ic[i] }
func (ic InstalledCameraList) Less(i, j int) bool { return ic[i].Install.Less(ic[j].Install) }

func (ic InstalledCameraList) encode() [][]string {
	var data [][]string

	data = append(data, installedCameraHeaders.Columns())

	for _, row := range ic {
		data = append(data, []string{
			strings.TrimSpace(row.Make),
			strings.TrimSpace(row.Model),
			strings.TrimSpace(row.Serial),
			strings.TrimSpace(row.Mount),
			strings.TrimSpace(row.View),
			strings.TrimSpace(row.dip),
			strings.TrimSpace(row.azimuth),
			strings.TrimSpace(row.vertical),
			strings.TrimSpace(row.north),
			strings.TrimSpace(row.east),
			row.Start.Format(DateTimeFormat),
			row.End.Format(DateTimeFormat),
			strings.TrimSpace(row.Notes),
		})
	}

	return data
}

func (ic *InstalledCameraList) decode(data [][]string) error {

	if !(len(data) > 1) {
		return nil
	}

	var cameras []InstalledCamera

	fields := installedCameraHeaders.Fields(data[0])
	for _, row := range data[1:] {
		d := fields.Remap(row)

		dip, err := strconv.ParseFloat(d[installedCameraDip], 64)
		if err != nil {
			return err
		}
		azimuth, err := strconv.ParseFloat(d[installedCameraAzimuth], 64)
		if err != nil {
			return err
		}

		height, err := strconv.ParseFloat(d[installedCameraHeight], 64)
		if err != nil {
			return err
		}
		north, err := strconv.ParseFloat(d[installedCameraNorth], 64)
		if err != nil {
			return err
		}
		east, err := strconv.ParseFloat(d[installedCameraEast], 64)
		if err != nil {
			return err
		}

		start, err := time.Parse(DateTimeFormat, d[installedCameraStart])
		if err != nil {
			return err
		}
		end, err := time.Parse(DateTimeFormat, d[installedCameraEnd])
		if err != nil {
			return err
		}

		cameras = append(cameras, InstalledCamera{
			Install: Install{
				Equipment: Equipment{
					Make:   strings.TrimSpace(d[installedCameraMake]),
					Model:  strings.TrimSpace(d[installedCameraModel]),
					Serial: strings.TrimSpace(d[installedCameraSerial]),
				},
				Span: Span{
					Start: start,
					End:   end,
				},
			},
			Orientation: Orientation{
				Dip:     dip,
				Azimuth: azimuth,

				dip:     strings.TrimSpace(d[installedCameraDip]),
				azimuth: strings.TrimSpace(d[installedCameraAzimuth]),
			},
			Offset: Offset{
				Vertical: height,
				North:    north,
				East:     east,

				vertical: strings.TrimSpace(d[installedCameraHeight]),
				north:    strings.TrimSpace(d[installedCameraNorth]),
				east:     strings.TrimSpace(d[installedCameraEast]),
			},
			Mount: strings.TrimSpace(d[installedCameraMount]),
			View:  strings.TrimSpace(d[installedCameraView]),
			Notes: strings.TrimSpace(d[installedCameraNotes]),
		})
	}

	*ic = InstalledCameraList(cameras)

	return nil
}

func LoadInstalledCameras(path string) ([]InstalledCamera, error) {
	var ic []InstalledCamera

	if err := LoadList(path, (*InstalledCameraList)(&ic)); err != nil {
		return nil, err
	}

	sort.Sort(InstalledCameraList(ic))

	return ic, nil
}
