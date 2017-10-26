package meta

import (
	"fmt"
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

type Visibility struct {
	Span
	Code          string
	SkyVisibility string
}

type VisibilityList []Visibility

func (m VisibilityList) Len() int      { return len(m) }
func (m VisibilityList) Swap(i, j int) { m[i], m[j] = m[j], m[i] }
func (m VisibilityList) Less(i, j int) bool {
	switch {
	case m[i].Code < m[j].Code:
		return true
	case m[i].Code > m[j].Code:
		return false
	default:
		return m[i].Start.Before(m[j].Start)
	}
}

func (m VisibilityList) encode() [][]string {
	data := [][]string{{
		"Code",
		"Sky Visibility",
		"Start Date",
		"End Date",
	}}
	for _, v := range m {
		data = append(data, []string{
			strings.TrimSpace(v.Code),
			strings.TrimSpace(v.SkyVisibility),
			v.Start.Format(DateTimeFormat),
			v.End.Format(DateTimeFormat),
		})
	}
	return data
}

func (m *VisibilityList) decode(data [][]string) error {
	var visibilities []Visibility
	if len(data) > 1 {
		for _, d := range data[1:] {
			if len(d) != visibilityLast {
				return fmt.Errorf("incorrect number of installed visibility fields")
			}
			var err error

			var start, end time.Time
			if start, err = time.Parse(DateTimeFormat, d[visibilityStartTime]); err != nil {
				return err
			}
			if end, err = time.Parse(DateTimeFormat, d[visibilityEndTime]); err != nil {
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

		*m = VisibilityList(visibilities)
	}
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
