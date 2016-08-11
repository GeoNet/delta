package meta

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

type InstalledCamera struct {
	Install
	Orientation
	Offset

	MountCode string
	Notes     string
}

type InstalledCameraList []InstalledCamera

func (a InstalledCameraList) Len() int           { return len(a) }
func (a InstalledCameraList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a InstalledCameraList) Less(i, j int) bool { return a[i].Install.less(a[j].Install) }

func (a InstalledCameraList) encode() [][]string {
	data := [][]string{{
		"Camera Make",
		"Camera Model",
		"Serial Number",
		"Mount Code",
		"Dip",
		"Azimuth",
		"Camera Height",
		"Offset North",
		"Offset East",
		"Installation Date",
		"Removal Date",
		"Notes",
	}}
	for _, v := range a {
		data = append(data, []string{
			strings.TrimSpace(v.Make),
			strings.TrimSpace(v.Model),
			strings.TrimSpace(v.Serial),
			strings.TrimSpace(v.MountCode),
			strconv.FormatFloat(v.Dip, 'g', -1, 64),
			strconv.FormatFloat(v.Azimuth, 'g', -1, 64),
			strconv.FormatFloat(v.Vertical, 'g', -1, 64),
			strconv.FormatFloat(v.North, 'g', -1, 64),
			strconv.FormatFloat(v.East, 'g', -1, 64),
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
			if len(d) != 12 {
				return fmt.Errorf("incorrect number of installed camera fields")
			}
			var err error

			var dip, azimuth float64
			if dip, err = strconv.ParseFloat(d[4], 64); err != nil {
				return err
			}
			if azimuth, err = strconv.ParseFloat(d[5], 64); err != nil {
				return err
			}

			var height, north, east float64
			if height, err = strconv.ParseFloat(d[6], 64); err != nil {
				return err
			}
			if north, err = strconv.ParseFloat(d[7], 64); err != nil {
				return err
			}
			if east, err = strconv.ParseFloat(d[8], 64); err != nil {
				return err
			}

			var start, end time.Time
			if start, err = time.Parse(DateTimeFormat, d[9]); err != nil {
				return err
			}
			if end, err = time.Parse(DateTimeFormat, d[10]); err != nil {
				return err
			}

			cameras = append(cameras, InstalledCamera{
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
					Dip:     dip,
					Azimuth: azimuth,
				},
				Offset: Offset{
					Vertical: height,
					North:    north,
					East:     east,
				},
				MountCode: strings.TrimSpace(d[3]),
				Notes:     strings.TrimSpace(d[11]),
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
