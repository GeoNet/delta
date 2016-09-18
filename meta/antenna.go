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
	antennaSerial
	antennaMark
	antennaHeight
	antennaNorth
	antennaEast
	antennaAzimuth
	antennaStart
	antennaEnd
	antennaLast
)

type InstalledAntenna struct {
	Install
	Offset

	Mark    string
	Azimuth float64
}

type InstalledAntennaList []InstalledAntenna

func (a InstalledAntennaList) Len() int           { return len(a) }
func (a InstalledAntennaList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a InstalledAntennaList) Less(i, j int) bool { return a[i].Install.less(a[j].Install) }

func (a InstalledAntennaList) encode() [][]string {
	data := [][]string{{
		"Make",
		"Model",
		"Serial",
		"Mark",
		"Height",
		"North",
		"East",
		"Azimuth",
		"Start Date",
		"End Date",
	}}
	for _, v := range a {
		data = append(data, []string{
			strings.TrimSpace(v.Make),
			strings.TrimSpace(v.Model),
			strings.TrimSpace(v.Serial),
			strings.TrimSpace(v.Mark),
			strconv.FormatFloat(v.Vertical, 'g', -1, 64),
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
			if north, err = strconv.ParseFloat(d[antennaNorth], 64); err != nil {
				return err
			}
			if east, err = strconv.ParseFloat(d[antennaEast], 64); err != nil {
				return err
			}

			var azimuth float64
			if azimuth, err = strconv.ParseFloat(d[antennaAzimuth], 64); err != nil {
				return err
			}

			var start, end time.Time
			if start, err = time.Parse(DateTimeFormat, d[antennaStart]); err != nil {
				return err
			}
			if end, err = time.Parse(DateTimeFormat, d[antennaEnd]); err != nil {
				return err
			}

			antennas = append(antennas, InstalledAntenna{
				Install: Install{
					Equipment: Equipment{
						Make:   strings.TrimSpace(d[antennaMake]),
						Model:  strings.TrimSpace(d[antennaModel]),
						Serial: strings.TrimSpace(d[antennaSerial]),
					},
					Span: Span{
						Start: start,
						End:   end,
					},
				},
				Offset: Offset{
					Vertical: height,
					North:    north,
					East:     east,
				},
				Mark:    strings.TrimSpace(d[antennaMark]),
				Azimuth: azimuth,
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
