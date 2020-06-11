package meta

import (
	"fmt"
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

type InstalledCamera struct {
	Install
	Orientation
	Offset

	Mount string
	Notes string
}

type InstalledCameraList []InstalledCamera

func (a InstalledCameraList) Len() int           { return len(a) }
func (a InstalledCameraList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a InstalledCameraList) Less(i, j int) bool { return a[i].Install.Less(a[j].Install) }

func (a InstalledCameraList) encode() [][]string {
	data := [][]string{{
		"Make",
		"Model",
		"Serial",
		"Mount",
		"Dip",
		"Azimuth",
		"Height",
		"North",
		"East",
		"Start Date",
		"End Date",
		"Notes",
	}}
	for _, v := range a {
		data = append(data, []string{
			strings.TrimSpace(v.Make),
			strings.TrimSpace(v.Model),
			strings.TrimSpace(v.Serial),
			strings.TrimSpace(v.Mount),
			strings.TrimSpace(v.dip),
			strings.TrimSpace(v.azimuth),
			strings.TrimSpace(v.vertical),
			strings.TrimSpace(v.north),
			strings.TrimSpace(v.east),
			v.Start.Format(DateTimeFormat),
			v.End.Format(DateTimeFormat),
			strings.TrimSpace(v.Notes),
		})
	}
	return data
}

func (a *InstalledCameraList) decode(data [][]string) error {
	var cameras []InstalledCamera
	if len(data) > 1 {
		for _, d := range data[1:] {
			if len(d) != installedCameraLast {
				return fmt.Errorf("incorrect number of installed camera fields")
			}
			var err error

			var dip, azimuth float64
			if dip, err = strconv.ParseFloat(d[installedCameraDip], 64); err != nil {
				return err
			}
			if azimuth, err = strconv.ParseFloat(d[installedCameraAzimuth], 64); err != nil {
				return err
			}

			var height, north, east float64
			if height, err = strconv.ParseFloat(d[installedCameraHeight], 64); err != nil {
				return err
			}
			if north, err = strconv.ParseFloat(d[installedCameraNorth], 64); err != nil {
				return err
			}
			if east, err = strconv.ParseFloat(d[installedCameraEast], 64); err != nil {
				return err
			}

			var start, end time.Time
			if start, err = time.Parse(DateTimeFormat, d[installedCameraStart]); err != nil {
				return err
			}
			if end, err = time.Parse(DateTimeFormat, d[installedCameraEnd]); err != nil {
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
				Notes: strings.TrimSpace(d[installedCameraNotes]),
			})
		}

		*a = InstalledCameraList(cameras)
	}
	return nil
}

func LoadInstalledCameras(path string) ([]InstalledCamera, error) {
	var a []InstalledCamera

	if err := LoadList(path, (*InstalledCameraList)(&a)); err != nil {
		return nil, err
	}

	sort.Sort(InstalledCameraList(a))

	return a, nil
}
