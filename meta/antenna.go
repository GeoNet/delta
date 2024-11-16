package meta

import (
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	installedAntennaMake int = iota
	installedAntennaModel
	installedAntennaSerial
	installedAntennaMark
	installedAntennaHeight
	installedAntennaNorth
	installedAntennaEast
	installedAntennaAzimuth
	installedAntennaStart
	installedAntennaEnd
	installedAntennaLast
)

var installedAntennaHeaders Header = map[string]int{
	"Make":       installedAntennaMake,
	"Model":      installedAntennaModel,
	"Serial":     installedAntennaSerial,
	"Mark":       installedAntennaMark,
	"Height":     installedAntennaHeight,
	"North":      installedAntennaNorth,
	"East":       installedAntennaEast,
	"Azimuth":    installedAntennaAzimuth,
	"Start Date": installedAntennaStart,
	"End Date":   installedAntennaEnd,
}

var InstalledAntennaTable Table = Table{
	name:    "Antenna",
	headers: installedAntennaHeaders,
	primary: []string{"Make", "Model", "Serial", "Mark", "Start Date"},
	native:  []string{"Height", "North", "East", "Azimuth"},
	foreign: map[string]map[string]string{
		"Asset": {"Make": "Make", "Model": "Model", "Serial": "Serial"},
		"Mark":  {"Mark": "Mark"},
	},
	remap: map[string]string{
		"Start Date": "Start",
		"End Date":   "End",
	},
	start: "Start Date",
	end:   "End Date",
}

type InstalledAntenna struct {
	Install
	Offset

	Mark    string
	Azimuth float64

	azimuth string // shadow variable to maintain formatting
}

type InstalledAntennaList []InstalledAntenna

func (ia InstalledAntennaList) Len() int           { return len(ia) }
func (ia InstalledAntennaList) Swap(i, j int)      { ia[i], ia[j] = ia[j], ia[i] }
func (ia InstalledAntennaList) Less(i, j int) bool { return ia[i].Install.Less(ia[j].Install) }

func (ia InstalledAntennaList) encode() [][]string {
	var data [][]string

	data = append(data, installedAntennaHeaders.Columns())

	for _, row := range ia {
		data = append(data, []string{
			strings.TrimSpace(row.Make),
			strings.TrimSpace(row.Model),
			strings.TrimSpace(row.Serial),
			strings.TrimSpace(row.Mark),
			strings.TrimSpace(row.vertical),
			strings.TrimSpace(row.north),
			strings.TrimSpace(row.east),
			strings.TrimSpace(row.azimuth),
			row.Start.Format(DateTimeFormat),
			row.End.Format(DateTimeFormat),
		})
	}
	return data
}

func (ia *InstalledAntennaList) decode(data [][]string) error {
	if !(len(data) > 1) {
		return nil
	}

	var installedAntennas []InstalledAntenna

	fields := installedAntennaHeaders.Fields(data[0])
	for _, row := range data[1:] {
		d := fields.Remap(row)

		height, err := strconv.ParseFloat(d[installedAntennaHeight], 64)
		if err != nil {
			return err
		}
		north, err := strconv.ParseFloat(d[installedAntennaNorth], 64)
		if err != nil {
			return err
		}
		east, err := strconv.ParseFloat(d[installedAntennaEast], 64)
		if err != nil {
			return err
		}

		azimuth, err := strconv.ParseFloat(d[installedAntennaAzimuth], 64)
		if err != nil {
			return err
		}

		start, err := time.Parse(DateTimeFormat, d[installedAntennaStart])
		if err != nil {
			return err
		}
		end, err := time.Parse(DateTimeFormat, d[installedAntennaEnd])
		if err != nil {
			return err
		}

		installedAntennas = append(installedAntennas, InstalledAntenna{
			Install: Install{
				Equipment: Equipment{
					Make:   strings.TrimSpace(d[installedAntennaMake]),
					Model:  strings.TrimSpace(d[installedAntennaModel]),
					Serial: strings.TrimSpace(d[installedAntennaSerial]),
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

				vertical: strings.TrimSpace(d[installedAntennaHeight]),
				north:    strings.TrimSpace(d[installedAntennaNorth]),
				east:     strings.TrimSpace(d[installedAntennaEast]),
			},
			Mark:    strings.TrimSpace(d[installedAntennaMark]),
			Azimuth: azimuth,
			azimuth: strings.TrimSpace(d[installedAntennaAzimuth]),
		})
	}

	*ia = InstalledAntennaList(installedAntennas)

	return nil
}

func LoadInstalledAntennas(path string) ([]InstalledAntenna, error) {
	var ia []InstalledAntenna

	if err := LoadList(path, (*InstalledAntennaList)(&ia)); err != nil {
		return nil, err
	}

	sort.Sort(InstalledAntennaList(ia))

	return ia, nil
}
