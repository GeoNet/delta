package meta

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	installedDoasMake int = iota
	installedDoasModel
	installedDoasSerial
	installedDoasMount
	installedDoasView
	installedDoasDip
	installedDoasAzimuth
	installedDoasHeight
	installedDoasNorth
	installedDoasEast
	installedDoasStart
	installedDoasEnd
	installedDoasLast
)

type InstalledDoas struct {
	Install
	Orientation
	Offset

	Mount string
	View  string
}

type InstalledDoasList []InstalledDoas

func (a InstalledDoasList) Len() int           { return len(a) }
func (a InstalledDoasList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a InstalledDoasList) Less(i, j int) bool { return a[i].Install.Less(a[j].Install) }

func (a InstalledDoasList) encode() [][]string {
	data := [][]string{{
		"Make",
		"Model",
		"Serial",
		"Mount",
		"View",
		"Dip",
		"Azimuth",
		"Height",
		"North",
		"East",
		"Start Date",
		"End Date",
	}}
	for _, v := range a {
		data = append(data, []string{
			strings.TrimSpace(v.Make),
			strings.TrimSpace(v.Model),
			strings.TrimSpace(v.Serial),
			strings.TrimSpace(v.Mount),
			strings.TrimSpace(v.View),
			strings.TrimSpace(v.dip),
			strings.TrimSpace(v.azimuth),
			strings.TrimSpace(v.vertical),
			strings.TrimSpace(v.north),
			strings.TrimSpace(v.east),
			v.Start.Format(DateTimeFormat),
			v.End.Format(DateTimeFormat),
		})
	}
	return data
}

func (a *InstalledDoasList) decode(data [][]string) error {
	var doases []InstalledDoas
	if len(data) > 1 {
		for _, d := range data[1:] {
			if len(d) != installedDoasLast {
				return fmt.Errorf("incorrect number of installed doas fields")
			}
			var err error

			var dip, azimuth float64
			if dip, err = strconv.ParseFloat(d[installedDoasDip], 64); err != nil {
				return err
			}
			if azimuth, err = strconv.ParseFloat(d[installedDoasAzimuth], 64); err != nil {
				return err
			}

			var height, north, east float64
			if height, err = strconv.ParseFloat(d[installedDoasHeight], 64); err != nil {
				return err
			}
			if north, err = strconv.ParseFloat(d[installedDoasNorth], 64); err != nil {
				return err
			}
			if east, err = strconv.ParseFloat(d[installedDoasEast], 64); err != nil {
				return err
			}

			var start, end time.Time
			if start, err = time.Parse(DateTimeFormat, d[installedDoasStart]); err != nil {
				return err
			}
			if end, err = time.Parse(DateTimeFormat, d[installedDoasEnd]); err != nil {
				return err
			}

			doases = append(doases, InstalledDoas{
				Install: Install{
					Equipment: Equipment{
						Make:   strings.TrimSpace(d[installedDoasMake]),
						Model:  strings.TrimSpace(d[installedDoasModel]),
						Serial: strings.TrimSpace(d[installedDoasSerial]),
					},
					Span: Span{
						Start: start,
						End:   end,
					},
				},
				Orientation: Orientation{
					Dip:     dip,
					Azimuth: azimuth,

					dip:     strings.TrimSpace(d[installedDoasDip]),
					azimuth: strings.TrimSpace(d[installedDoasAzimuth]),
				},
				Offset: Offset{
					Vertical: height,
					North:    north,
					East:     east,

					vertical: strings.TrimSpace(d[installedDoasHeight]),
					north:    strings.TrimSpace(d[installedDoasNorth]),
					east:     strings.TrimSpace(d[installedDoasEast]),
				},
				Mount: strings.TrimSpace(d[installedDoasMount]),
				View:  strings.TrimSpace(d[installedDoasView]),
			})
		}

		*a = InstalledDoasList(doases)
	}
	return nil
}

func LoadInstalledDoass(path string) ([]InstalledDoas, error) {
	var a []InstalledDoas

	if err := LoadList(path, (*InstalledDoasList)(&a)); err != nil {
		return nil, err
	}

	sort.Sort(InstalledDoasList(a))

	return a, nil
}
