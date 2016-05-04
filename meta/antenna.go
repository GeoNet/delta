package meta

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	antennaMake int = iota
	antennaModel
	antennaSerialNumber
	antennaMarkCode
	antennaHeight
	antennaOffsetNorth
	antennaOffsetEast
	antennaAzimuth
	antennaInstallationDate
	antennaRemovalDate
	antennaLast
)

type InstalledAntenna struct {
	Install
	Offset

	MarkCode string
	Azimuth  float64
}

type InstalledAntennaList []InstalledAntenna

func (a InstalledAntennaList) Len() int           { return len(a) }
func (a InstalledAntennaList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a InstalledAntennaList) Less(i, j int) bool { return a[i].Install.less(a[j].Install) }

func (a InstalledAntennaList) encode() [][]string {
	data := [][]string{{
		"Antenna Make",
		"Antenna Model",
		"Serial Number",
		"Mark Code",
		"Antenna Height",
		"Offset North",
		"Offset East",
		"Azimuth",
		"Installation Date",
		"Removal Date",
	}}
	for _, v := range a {
		data = append(data, []string{
			strings.TrimSpace(v.Make),
			strings.TrimSpace(v.Model),
			strings.TrimSpace(v.Serial),
			strings.TrimSpace(v.MarkCode),
			strconv.FormatFloat(v.Height, 'g', -1, 64),
			strconv.FormatFloat(v.North, 'g', -1, 64),
			strconv.FormatFloat(v.East, 'g', -1, 64),
			strconv.FormatFloat(v.Azimuth, 'g', -1, 64),
			v.Start.Format(DateTimeFormat),
			v.End.Format(DateTimeFormat),
		})
	}
	return data
}

func (a *InstalledAntennaList) decode(data [][]string) error {
	var antennas []InstalledAntenna
	if len(data) > 1 {
		for _, d := range data[1:] {
			if len(d) != antennaLast {
				return fmt.Errorf("incorrect number of installed antenna fields")
			}
			var err error

			var height, north, east float64
			if height, err = strconv.ParseFloat(d[antennaHeight], 64); err != nil {
				return err
			}
			if north, err = strconv.ParseFloat(d[antennaOffsetNorth], 64); err != nil {
				return err
			}
			if east, err = strconv.ParseFloat(d[antennaOffsetEast], 64); err != nil {
				return err
			}

			var azimuth float64
			if azimuth, err = strconv.ParseFloat(d[antennaAzimuth], 64); err != nil {
				return err
			}

			var start, end time.Time
			if start, err = time.Parse(DateTimeFormat, d[antennaInstallationDate]); err != nil {
				return err
			}
			if end, err = time.Parse(DateTimeFormat, d[antennaRemovalDate]); err != nil {
				return err
			}

			antennas = append(antennas, InstalledAntenna{
				Install: Install{
					Equipment: Equipment{
						Make:   strings.TrimSpace(d[antennaMake]),
						Model:  strings.TrimSpace(d[antennaModel]),
						Serial: strings.TrimSpace(d[antennaSerialNumber]),
					},
					Span: Span{
						Start: start,
						End:   end,
					},
				},
				Offset: Offset{
					Height: height,
					North:  north,
					East:   east,
				},
				MarkCode: strings.TrimSpace(d[antennaMarkCode]),
				Azimuth:  azimuth,
			})
		}

		*a = InstalledAntennaList(antennas)
	}
	return nil
}

func LoadInstalledAntennas(path string) ([]InstalledAntenna, error) {
	var a []InstalledAntenna

	if err := LoadList(path, (*InstalledAntennaList)(&a)); err != nil {
		return nil, err
	}

	sort.Sort(InstalledAntennaList(a))

	return a, nil
}
