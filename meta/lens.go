package meta

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

const (
	lensMake = iota
	lensModel
	lensSerial
	lensMount
	lensView
	lensType
	lensStart
	lensEnd
	lensLast
)

type InstalledLens struct {
	Install

	Mount string
	View  string
	Type  string
}

type InstalledLensList []InstalledLens

func (s InstalledLensList) Len() int           { return len(s) }
func (s InstalledLensList) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s InstalledLensList) Less(i, j int) bool { return s[i].Install.Less(s[j].Install) }

func (s InstalledLensList) encode() [][]string {
	data := [][]string{{
		"Make",
		"Model",
		"Serial",
		"Mount",
		"View",
		"Type",
		"Start Date",
		"End Date",
	}}

	for _, v := range s {
		data = append(data, []string{
			strings.TrimSpace(v.Make),
			strings.TrimSpace(v.Model),
			strings.TrimSpace(v.Serial),
			strings.TrimSpace(v.Mount),
			strings.TrimSpace(v.View),
			strings.TrimSpace(v.Type),
			v.Start.Format(DateTimeFormat),
			v.End.Format(DateTimeFormat),
		})
	}
	return data
}
func (s *InstalledLensList) decode(data [][]string) error {
	var lenses []InstalledLens
	if len(data) > 1 {
		for _, d := range data[1:] {
			if len(d) != lensLast {
				return fmt.Errorf("incorrect number of installed lens fields")
			}
			var err error

			var start, end time.Time
			if start, err = time.Parse(DateTimeFormat, d[lensStart]); err != nil {
				return err
			}
			if end, err = time.Parse(DateTimeFormat, d[lensEnd]); err != nil {
				return err
			}

			lenses = append(lenses, InstalledLens{
				Install: Install{
					Equipment: Equipment{
						Make:   strings.TrimSpace(d[lensMake]),
						Model:  strings.TrimSpace(d[lensModel]),
						Serial: strings.TrimSpace(d[lensSerial]),
					},
					Span: Span{
						Start: start,
						End:   end,
					},
				},
				Mount: strings.TrimSpace(d[lensMount]),
				View:  strings.TrimSpace(d[lensView]),
				Type:  strings.TrimSpace(d[lensType]),
			})
		}

		*s = InstalledLensList(lenses)
	}
	return nil
}

func LoadInstalledLenses(path string) ([]InstalledLens, error) {
	var s []InstalledLens

	if err := LoadList(path, (*InstalledLensList)(&s)); err != nil {
		return nil, err
	}

	sort.Sort(InstalledLensList(s))

	return s, nil
}
