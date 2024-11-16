package meta

import (
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	viewMount = iota
	viewCode
	viewLabel
	viewAzimuth
	viewMethod
	viewDip
	viewDescription
	viewStart
	viewEnd
	viewLast
)

var viewHeaders Header = map[string]int{
	"Mount":       viewMount,
	"View":        viewCode,
	"Label":       viewLabel,
	"Azimuth":     viewAzimuth,
	"Method":      viewMethod,
	"Dip":         viewDip,
	"Description": viewDescription,
	"Start Date":  viewStart,
	"End Date":    viewEnd,
}

var ViewTable Table = Table{
	name:    "View",
	headers: viewHeaders,
	primary: []string{"Mount", "View", "Start Date"},
	native:  []string{"Azimuth", "Dip"},
	foreign: map[string]map[string]string{
		"Mount": {"Mount": "Mount"},
	},
	nullable: []string{"Label", "Method"},
	remap: map[string]string{
		"Start Date": "Start",
		"End Date":   "End",
	},
	start: "Start Date",
	end:   "End Date",
}

type View struct {
	Orientation
	Span

	Mount       string `json:"mount"`
	Code        string `json:"code"`
	Label       string `json:"label"`
	Description string `json:"description"`
}

type ViewList []View

func (v ViewList) Len() int      { return len(v) }
func (v ViewList) Swap(i, j int) { v[i], v[j] = v[j], v[i] }
func (v ViewList) Less(i, j int) bool {
	switch {
	case v[i].Mount < v[j].Mount:
		return true
	case v[i].Mount > v[j].Mount:
		return false
	case v[i].Code < v[j].Code:
		return true
	default:
		return false
	}
}

func (v ViewList) encode() [][]string {
	var data [][]string

	data = append(data, viewHeaders.Columns())

	for _, row := range v {
		data = append(data, []string{
			strings.TrimSpace(row.Mount),
			strings.TrimSpace(row.Code),
			strings.TrimSpace(row.Label),
			strings.TrimSpace(row.azimuth),
			strings.TrimSpace(row.Method),
			strings.TrimSpace(row.dip),
			strings.TrimSpace(row.Description),
			row.Start.Format(DateTimeFormat),
			row.End.Format(DateTimeFormat),
		})
	}

	return data
}

func (v *ViewList) decode(data [][]string) error {
	if !(len(data) > 1) {
		return nil
	}

	var views []View

	fields := viewHeaders.Fields(data[0])
	for _, row := range data[1:] {
		d := fields.Remap(row)

		dip, err := strconv.ParseFloat(d[viewDip], 64)
		if err != nil {
			return err
		}
		azimuth, err := strconv.ParseFloat(d[viewAzimuth], 64)
		if err != nil {
			return err
		}

		start, err := time.Parse(DateTimeFormat, d[viewStart])
		if err != nil {
			return err
		}
		end, err := time.Parse(DateTimeFormat, d[viewEnd])
		if err != nil {
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
				Method:  strings.TrimSpace(d[viewMethod]),

				azimuth: strings.TrimSpace(d[viewAzimuth]),
				dip:     strings.TrimSpace(d[viewDip]),
			},
			Span: Span{
				Start: start,
				End:   end,
			},
		})
	}

	*v = ViewList(views)

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
