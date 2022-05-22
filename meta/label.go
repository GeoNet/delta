package meta

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

const (
	labelType = iota
	labelSamplingRate
	labelAzimuth
	labelDip
	labelCode
	labelFlags
	labelLast
)

// Label is used to describe a channel stream.
type Label struct {
	Type  string
	Code  string
	Flags string

	SamplingRate float64
	Azimuth      float64
	Dip          float64

	azimuth      string
	dip          string
	samplingRate string
}

// Less compares Label structs suitable for sorting.
func (c Label) Less(label Label) bool {

	switch {
	case strings.ToLower(c.Type) < strings.ToLower(label.Type):
		return true
	case strings.ToLower(c.Type) > strings.ToLower(label.Type):
		return false
	case c.SamplingRate < label.SamplingRate:
		return true
	case c.SamplingRate > label.SamplingRate:
		return false
	case c.Dip < label.Dip:
		return true
	case c.Dip > label.Dip:
		return false
	case c.Azimuth < label.Azimuth:
		return true
	default:
		return false
	}
}

type LabelList []Label

func (s LabelList) Len() int           { return len(s) }
func (s LabelList) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s LabelList) Less(i, j int) bool { return s[i].Less(s[j]) }

func (s LabelList) encode() [][]string {
	data := [][]string{{
		"Type",
		"Sampling Rate",
		"Azimuth",
		"Dip",
		"Code",
		"Flags",
	}}

	for _, v := range s {
		data = append(data, []string{
			strings.TrimSpace(v.Type),
			strings.TrimSpace(v.samplingRate),
			strings.TrimSpace(v.azimuth),
			strings.TrimSpace(v.dip),
			strings.TrimSpace(v.Code),
			strings.TrimSpace(string(v.Flags)),
		})
	}
	return data
}
func (s *LabelList) decode(data [][]string) error {
	var labels []Label
	if len(data) > 1 {
		for _, d := range data[1:] {
			if len(d) != labelLast {
				return fmt.Errorf("incorrect number of label fields")
			}

			rate, err := strconv.ParseFloat(d[labelSamplingRate], 64)
			if err != nil {
				return err
			}
			if rate < 0 {
				rate = -1.0 / rate
			}

			azimuth, err := strconv.ParseFloat(d[labelAzimuth], 64)
			if err != nil {
				return err
			}
			dip, err := strconv.ParseFloat(d[labelDip], 64)
			if err != nil {
				return err
			}

			labels = append(labels, Label{
				Type:  strings.TrimSpace(d[labelType]),
				Code:  strings.TrimSpace(d[labelCode]),
				Flags: strings.TrimSpace(d[labelFlags]),

				SamplingRate: rate,
				Azimuth:      azimuth,
				Dip:          dip,

				samplingRate: strings.TrimSpace(d[labelSamplingRate]),
				azimuth:      strings.TrimSpace(d[labelAzimuth]),
				dip:          strings.TrimSpace(d[labelDip]),
			})
		}

		*s = LabelList(labels)
	}
	return nil
}

func LoadLabels(path string) ([]Label, error) {
	var s []Label

	if err := LoadList(path, (*LabelList)(&s)); err != nil {
		return nil, err
	}

	sort.Sort(LabelList(s))

	return s, nil
}
