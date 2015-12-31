package meta

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type InstalledRadome struct {
	Install

	MarkCode string
}

type InstalledRadomes []InstalledRadome

func (r InstalledRadomes) Len() int           { return len(r) }
func (r InstalledRadomes) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r InstalledRadomes) Less(i, j int) bool { return r[i].Install.Less(r[j].Install) }

func (r InstalledRadomes) encode() [][]string {
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

func (r *InstalledRadomes) decode(data [][]string) error {
	var radomes []InstalledRadome
	if len(data) > 1 {
		for _, d := range data[1:] {
			if len(d) != 6 {
				return fmt.Errorf("incorrect number of installed radome fields")
			}
			var err error

			var start, end time.Time
			if start, err = time.Parse(DateTimeFormat, d[4]); err != nil {
				return err
			}
			if end, err = time.Parse(DateTimeFormat, d[5]); err != nil {
				return err
			}

			radomes = append(radomes, InstalledRadome{
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
				MarkCode: strings.TrimSpace(d[3]),
			})
		}

		*r = InstalledRadomes(radomes)
	}
	return nil
}

func LoadInstalledRadomes(path string) ([]InstalledRadome, error) {
	var r []InstalledRadome

	if err := LoadList(path, (*InstalledRadomes)(&r)); err != nil {
		return nil, err
	}

	sort.Sort(InstalledRadomes(r))

	return r, nil
}
