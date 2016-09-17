package meta

import (
	"fmt"
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

type InstalledRadome struct {
	Install

	MarkCode string
}

type InstalledRadomeList []InstalledRadome

func (r InstalledRadomeList) Len() int           { return len(r) }
func (r InstalledRadomeList) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r InstalledRadomeList) Less(i, j int) bool { return r[i].Install.less(r[j].Install) }

func (r InstalledRadomeList) encode() [][]string {
	data := [][]string{{
		"Radome Make",
		"Radome Model",
		"Serial Number",
		"Mark Code",
		"Installation Date",
		"Removal Date",
	}}
	for _, v := range r {
		data = append(data, []string{
			strings.TrimSpace(v.Make),
			strings.TrimSpace(v.Model),
			strings.TrimSpace(v.Serial),
			strings.TrimSpace(v.MarkCode),
			v.Start.Format(DateTimeFormat),
			v.End.Format(DateTimeFormat),
		})
	}
	return data
}

func (r *InstalledRadomeList) decode(data [][]string) error {
	var radomes []InstalledRadome
	if len(data) > 1 {
		for _, d := range data[1:] {
			if len(d) != installedRadomeLast {
				return fmt.Errorf("incorrect number of installed radome fields")
			}
			var err error

			var start, end time.Time
			if start, err = time.Parse(DateTimeFormat, d[installedRadomeStart]); err != nil {
				return err
			}
			if end, err = time.Parse(DateTimeFormat, d[installedRadomeEnd]); err != nil {
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
				MarkCode: strings.TrimSpace(d[installedRadomeMark]),
			})
		}

		*r = InstalledRadomeList(radomes)
	}
	return nil
}

func LoadInstalledRadomes(path string) ([]InstalledRadome, error) {
	var r []InstalledRadome

	if err := LoadList(path, (*InstalledRadomeList)(&r)); err != nil {
		return nil, err
	}

	sort.Sort(InstalledRadomeList(r))

	return r, nil
}
