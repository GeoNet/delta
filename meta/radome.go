package meta

import (
	"sort"
	"strings"
	"time"
)

const (
	installedRadomeMake int = iota
	installedRadomeModel
	installedRadomeSerial
	installedRadomeMark
	installedRadomeStart
	installedRadomeEnd
	installedRadomeLast
)

var installedRadomeHeaders Header = map[string]int{
	"Make":       installedRadomeMake,
	"Model":      installedRadomeModel,
	"Serial":     installedRadomeSerial,
	"Mark":       installedRadomeMark,
	"Start Date": installedRadomeStart,
	"End Date":   installedRadomeEnd,
}

var InstalledRadomeTable Table = Table{
	name:    "Radome",
	headers: installedRadomeHeaders,
	primary: []string{"Make", "Model", "Serial", "Start Date"},
	native:  []string{},
	foreign: map[string]map[string]string{
		"Mark": {"Mark": "Mark"},
	},
	remap: map[string]string{
		"Start Date": "Start",
		"End Date":   "End",
	},
	start: "Start Date",
	end:   "End Date",
}

type InstalledRadome struct {
	Install

	Mark string
}

type InstalledRadomeList []InstalledRadome

func (ir InstalledRadomeList) Len() int           { return len(ir) }
func (ir InstalledRadomeList) Swap(i, j int)      { ir[i], ir[j] = ir[j], ir[i] }
func (ir InstalledRadomeList) Less(i, j int) bool { return ir[i].Install.Less(ir[j].Install) }

func (ir InstalledRadomeList) encode() [][]string {
	var data [][]string

	data = append(data, installedRadomeHeaders.Columns())

	for _, row := range ir {
		data = append(data, []string{
			strings.TrimSpace(row.Make),
			strings.TrimSpace(row.Model),
			strings.TrimSpace(row.Serial),
			strings.TrimSpace(row.Mark),
			row.Start.Format(DateTimeFormat),
			row.End.Format(DateTimeFormat),
		})
	}

	return data
}

func (ir *InstalledRadomeList) decode(data [][]string) error {
	if !(len(data) > 1) {
		return nil
	}

	var radomes []InstalledRadome

	fields := installedRadomeHeaders.Fields(data[0])
	for _, row := range data[1:] {
		d := fields.Remap(row)

		start, err := time.Parse(DateTimeFormat, d[installedRadomeStart])
		if err != nil {
			return err
		}
		end, err := time.Parse(DateTimeFormat, d[installedRadomeEnd])
		if err != nil {
			return err
		}

		radomes = append(radomes, InstalledRadome{
			Install: Install{
				Equipment: Equipment{
					Make:   strings.TrimSpace(d[installedRadomeMake]),
					Model:  strings.TrimSpace(d[installedRadomeModel]),
					Serial: strings.TrimSpace(d[installedRadomeSerial]),
				},
				Span: Span{
					Start: start,
					End:   end,
				},
			},
			Mark: strings.TrimSpace(d[installedRadomeMark]),
		})
	}

	*ir = InstalledRadomeList(radomes)

	return nil
}

func LoadInstalledRadomes(path string) ([]InstalledRadome, error) {
	var ir []InstalledRadome

	if err := LoadList(path, (*InstalledRadomeList)(&ir)); err != nil {
		return nil, err
	}

	sort.Sort(InstalledRadomeList(ir))

	return ir, nil
}
