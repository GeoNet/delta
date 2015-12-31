package meta

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

type InstalledAntenna struct {
	Install
	Offset

	MarkCode string
}

type InstalledAntennas []InstalledAntenna

func (a InstalledAntennas) Len() int           { return len(a) }
func (a InstalledAntennas) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a InstalledAntennas) Less(i, j int) bool { return a[i].Install.Less(a[j].Install) }

func (a InstalledAntennas) encode() [][]string {
	data := [][]string{{
		"Antenna Make",
		"Antenna Model",
		"Serial Number",
		"Mark Code",
		"Antenna Height",
		"Offset North",
		"Offset East",
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
			v.Start.Format(DateTimeFormat),
			v.End.Format(DateTimeFormat),
		})
	}
	return data
}

func (a *InstalledAntennas) decode(data [][]string) error {
	var antennas []InstalledAntenna
	if len(data) > 1 {
		for _, d := range data[1:] {
			if len(d) != 9 {
				return fmt.Errorf("incorrect number of installed antenna fields")
			}
			var err error

			var height, north, east float64
			if height, err = strconv.ParseFloat(d[4], 64); err != nil {
				return err
			}
			if north, err = strconv.ParseFloat(d[5], 64); err != nil {
				return err
			}
			if east, err = strconv.ParseFloat(d[6], 64); err != nil {
				return err
			}

			var start, end time.Time
			if start, err = time.Parse(DateTimeFormat, d[7]); err != nil {
				return err
			}
			if end, err = time.Parse(DateTimeFormat, d[8]); err != nil {
				return err
			}

			antennas = append(antennas, InstalledAntenna{
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
				Offset: Offset{
					Height: height,
					North:  north,
					East:   east,
				},
				MarkCode: strings.TrimSpace(d[3]),
			})
		}

		*a = InstalledAntennas(antennas)
	}
	return nil
}

func LoadInstalledAntennas(path string) ([]InstalledAntenna, error) {
	var a []InstalledAntenna

	if err := LoadList(path, (*InstalledAntennas)(&a)); err != nil {
		return nil, err
	}

	sort.Sort(InstalledAntennas(a))

	return a, nil
}
