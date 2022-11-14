package meta

import (
	"fmt"
	"sort"
	"strings"
)

const (
	componentMake = iota
	componentModel
	componentType
	componentNumber
	componentSubsource
	componentDip
	componentAzimuth
	componentTypes
	componentSamplingRate
	componentResponse
	componentLast
)

type Component struct {
	Make         string
	Model        string
	Type         string
	Number       int
	Subsource    string
	Dip          float64
	Azimuth      float64
	Types        string
	SamplingRate float64
	Response     string

	number       string
	dip          string
	azimuth      string
	samplingRate string
}

// Less compares Component structs suitable for sorting.
func (c Component) Less(comp Component) bool {

	switch {
	case strings.ToLower(c.Make) < strings.ToLower(comp.Make):
		return true
	case strings.ToLower(c.Make) > strings.ToLower(comp.Make):
		return false
	case strings.ToLower(c.Model) < strings.ToLower(comp.Model):
		return true
	case strings.ToLower(c.Model) > strings.ToLower(comp.Model):
		return false
	case c.Number < comp.Number:
		return true
	case c.Number > comp.Number:
		return false
	case c.SamplingRate < comp.SamplingRate:
		return true
	default:
		return false
	}
}

type ComponentList []Component

func (s ComponentList) Len() int           { return len(s) }
func (s ComponentList) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s ComponentList) Less(i, j int) bool { return s[i].Less(s[j]) }

func (s ComponentList) encode() [][]string {
	data := [][]string{{
		"Make",
		"Model",
		"Type",
		"Number",
		"Subsource",
		"Dip",
		"Azimuth",
		"Types",
		"Sampling Rate",
		"Response",
	}}

	for _, v := range s {
		data = append(data, []string{
			strings.TrimSpace(v.Make),
			strings.TrimSpace(v.Model),
			strings.TrimSpace(v.Type),
			strings.TrimSpace(v.number),
			strings.TrimSpace(v.Subsource),
			strings.TrimSpace(v.dip),
			strings.TrimSpace(v.azimuth),
			strings.TrimSpace(v.Types),
			strings.TrimSpace(v.samplingRate),
			strings.TrimSpace(v.Response),
		})
	}
	return data
}
func (s *ComponentList) decode(data [][]string) error {
	var components []Component

	if !(len(data) > 1) {
		return nil
	}

	for _, d := range data[1:] {
		if len(d) != componentLast {
			return fmt.Errorf("incorrect pin of installed component fields")
		}

		number, err := ParseInt(d[componentNumber])
		if err != nil {
			return err
		}

		dip, err := ParseFloat64(d[componentDip])
		if err != nil {
			return err
		}

		azimuth, err := ParseFloat64(d[componentAzimuth])
		if err != nil {
			return err
		}

		samplingRate, err := ParseFloat64(d[componentSamplingRate])
		if err != nil {
			return err
		}

		if samplingRate < 0.0 {
			samplingRate = -1.0 / samplingRate
		}

		components = append(components, Component{
			Make:         strings.TrimSpace(d[componentMake]),
			Model:        strings.TrimSpace(d[componentModel]),
			Type:         strings.TrimSpace(d[componentType]),
			Number:       number,
			Subsource:    strings.TrimSpace(d[componentSubsource]),
			Dip:          dip,
			Azimuth:      azimuth,
			Types:        strings.TrimSpace(d[componentTypes]),
			SamplingRate: samplingRate,
			Response:     strings.TrimSpace(d[componentResponse]),

			number:       strings.TrimSpace(d[componentNumber]),
			dip:          strings.TrimSpace(d[componentDip]),
			azimuth:      strings.TrimSpace(d[componentAzimuth]),
			samplingRate: strings.TrimSpace(d[componentSamplingRate]),
		})
	}

	*s = ComponentList(components)

	return nil
}

func LoadComponents(path string) ([]Component, error) {
	var s []Component

	if err := LoadList(path, (*ComponentList)(&s)); err != nil {
		return nil, err
	}

	sort.Sort(ComponentList(s))

	return s, nil
}
