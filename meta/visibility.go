package meta

import (
	"sort"
	"strings"
	"time"
)

const (
	visibilityCode = iota
	visibilitySkyVisibility
	visibilityStartTime
	visibilityEndTime
	visibilityLast
)

var visibilityHeaders Header = map[string]int{
	"Code":           visibilityCode,
	"Sky Visibility": visibilitySkyVisibility,
	"Start Date":     visibilityStartTime,
	"End Date":       visibilityEndTime,
}

type Visibility struct {
	Span
	Code          string
	SkyVisibility string
}

type VisibilityList []Visibility

func (v VisibilityList) Len() int      { return len(v) }
func (v VisibilityList) Swap(i, j int) { v[i], v[j] = v[j], v[i] }
func (v VisibilityList) Less(i, j int) bool {
	switch {
	case v[i].Code < v[j].Code:
		return true
	case v[i].Code > v[j].Code:
		return false
	default:
		return v[i].Start.Before(v[j].Start)
	}
}

func (v VisibilityList) encode() [][]string {
	var data [][]string

	data = append(data, visibilityHeaders.Columns())

	for _, row := range v {
		data = append(data, []string{
			strings.TrimSpace(row.Code),
			strings.TrimSpace(row.SkyVisibility),
			row.Start.Format(DateTimeFormat),
			row.End.Format(DateTimeFormat),
		})
	}

	return data
}

func (v *VisibilityList) decode(data [][]string) error {
	if !(len(data) > 1) {
		return nil
	}

	var visibilities []Visibility

	fields := visibilityHeaders.Fields(data[0])
	for _, row := range data[1:] {
		d := fields.Remap(row)

		start, err := time.Parse(DateTimeFormat, d[visibilityStartTime])
		if err != nil {
			return err
		}
		end, err := time.Parse(DateTimeFormat, d[visibilityEndTime])
		if err != nil {
			return err
		}
		visibilities = append(visibilities, Visibility{
			Code:          strings.TrimSpace(d[visibilityCode]),
			SkyVisibility: strings.TrimSpace(d[visibilitySkyVisibility]),
			Span: Span{
				Start: start,
				End:   end,
			},
		})
	}

	*v = VisibilityList(visibilities)

	return nil
}

func LoadVisibilities(path string) ([]Visibility, error) {
	var v []Visibility

	if err := LoadList(path, (*VisibilityList)(&v)); err != nil {
		return nil, err
	}

	sort.Sort(VisibilityList(v))

	return v, nil
}
