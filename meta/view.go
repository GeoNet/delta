package meta

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	viewMount = iota
	viewCode
	viewLabel
	viewDip
	viewAzimuth
	viewDescription
	viewStart
	viewEnd
	viewLast
)

type View struct {
	Orientation
	Span

	Mount       string
	Code        string
	Label       string
	Description string
}

type ViewList []View

func (l ViewList) Len() int      { return len(l) }
func (l ViewList) Swap(i, j int) { l[i], l[j] = l[j], l[i] }
func (l ViewList) Less(i, j int) bool {
	switch {
	case l[i].Mount < l[j].Mount:
		return true
	case l[i].Mount > l[j].Mount:
		return false
	case l[i].Code < l[j].Code:
		return true
	default:
		return false
	}
}

func (m ViewList) encode() [][]string {
	data := [][]string{{
		"Mount",
		"View",
		"Label",
		"Dip",
		"Azimuth",
		"Description",
		"Start Date",
		"End Date",
	}}
	for _, l := range m {
		data = append(data, []string{
			strings.TrimSpace(l.Mount),
			strings.TrimSpace(l.Code),
			strings.TrimSpace(l.Label),
			strings.TrimSpace(l.dip),
			strings.TrimSpace(l.azimuth),
			strings.TrimSpace(l.Description),
			l.Start.Format(DateTimeFormat),
			l.End.Format(DateTimeFormat),
		})
	}
	return data
}

func (l *ViewList) decode(data [][]string) error {
	var views []View
	if len(data) > 1 {
		for _, d := range data[1:] {
			if len(d) != viewLast {
				return fmt.Errorf("incorrect number of installed view fields")
			}
			var err error

			var dip, azimuth float64
			if dip, err = strconv.ParseFloat(d[viewDip], 64); err != nil {
				return err
			}
			if azimuth, err = strconv.ParseFloat(d[viewAzimuth], 64); err != nil {
				return err
			}

			var start, end time.Time
			if start, err = time.Parse(DateTimeFormat, d[viewStart]); err != nil {
				return err
			}
			if end, err = time.Parse(DateTimeFormat, d[viewEnd]); err != nil {
				return err
			}

			views = append(views, View{
				Mount:       strings.TrimSpace(d[viewMount]),
				Code:        strings.TrimSpace(d[viewCode]),
				Label:       strings.TrimSpace(d[viewLabel]),
				Description: strings.TrimSpace(d[viewDescription]),
				Orientation: Orientation{
					Dip:     dip,
					Azimuth: azimuth,

					dip:     strings.TrimSpace(d[viewDip]),
					azimuth: strings.TrimSpace(d[viewAzimuth]),
				},
				Span: Span{
					Start: start,
					End:   end,
				},
			})
		}

		*l = ViewList(views)
	}
	return nil
}

func LoadViews(path string) ([]View, error) {
	var v []View

	if err := LoadList(path, (*ViewList)(&v)); err != nil {
		return nil, err
	}

	sort.Sort(ViewList(v))

	return v, nil
}
