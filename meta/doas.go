package meta

import (
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

var installedDoasHeaders Header = map[string]int{
	"Make":       installedDoasMake,
	"Model":      installedDoasModel,
	"Serial":     installedDoasSerial,
	"Mount":      installedDoasMount,
	"View":       installedDoasView,
	"Dip":        installedDoasDip,
	"Azimuth":    installedDoasAzimuth,
	"Height":     installedDoasHeight,
	"North":      installedDoasNorth,
	"East":       installedDoasEast,
	"Start Date": installedDoasStart,
	"End Date":   installedDoasEnd,
}

type InstalledDoas struct {
	Install
	Orientation
	Offset

	Mount string
	View  string
}

type InstalledDoasList []InstalledDoas

func (id InstalledDoasList) Len() int           { return len(id) }
func (id InstalledDoasList) Swap(i, j int)      { id[i], id[j] = id[j], id[i] }
func (id InstalledDoasList) Less(i, j int) bool { return id[i].Install.Less(id[j].Install) }

func (id InstalledDoasList) encode() [][]string {
	var data [][]string

	data = append(data, installedDoasHeaders.Columns())

	for _, row := range id {
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
		})
	}

	return data
}

func (id *InstalledDoasList) decode(data [][]string) error {
	if !(len(data) > 1) {
		return nil
	}

	var doases []InstalledDoas

	fields := installedDoasHeaders.Fields(data[0])
	for _, row := range data[1:] {
		d := fields.Remap(row)

		dip, err := strconv.ParseFloat(d[installedDoasDip], 64)
		if err != nil {
			return err
		}
		azimuth, err := strconv.ParseFloat(d[installedDoasAzimuth], 64)
		if err != nil {
			return err
		}

		height, err := strconv.ParseFloat(d[installedDoasHeight], 64)
		if err != nil {
			return err
		}
		north, err := strconv.ParseFloat(d[installedDoasNorth], 64)
		if err != nil {
			return err
		}
		east, err := strconv.ParseFloat(d[installedDoasEast], 64)
		if err != nil {
			return err
		}

		start, err := time.Parse(DateTimeFormat, d[installedDoasStart])
		if err != nil {
			return err
		}
		end, err := time.Parse(DateTimeFormat, d[installedDoasEnd])
		if err != nil {
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

	*id = InstalledDoasList(doases)

	return nil
}

func LoadInstalledDoass(path string) ([]InstalledDoas, error) {
	var id []InstalledDoas

	if err := LoadList(path, (*InstalledDoasList)(&id)); err != nil {
		return nil, err
	}

	sort.Sort(InstalledDoasList(id))

	return id, nil
}
