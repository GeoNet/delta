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
	componentSource
	componentSubsource
	componentDip
	componentAzimuth
	componentTypes
	componentSamplingRate
	componentResponse
	componentLast
)

var componentHeaders Header = map[string]int{
	"Make":          componentMake,
	"Model":         componentModel,
	"Type":          componentType,
	"Number":        componentNumber,
	"Source":        componentSource,
	"Subsource":     componentSubsource,
	"Dip":           componentDip,
	"Azimuth":       componentAzimuth,
	"Types":         componentTypes,
	"Sampling Rate": componentSamplingRate,
	"Response":      componentResponse,
}

var ComponentTable Table = Table{
	name:     "Component",
	headers:  componentHeaders,
	primary:  []string{"Make", "Model", "Number", "Source", "Subsource", "Sampling Rate"},
	native:   []string{"Number", "Dip", "Azimuth", "Sampling Rate"},
	foreign:  map[string]map[string]string{},
	nullable: []string{"Sampling Rate", "Source", "Type"},
	remap: map[string]string{
		"Sampling Rate": "SamplingRate",
	},
}

type Component struct {
	Make         string
	Model        string
	Type         string
	Number       int
	Source       string
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

// Description returns a short label for the channel model family.
func (c Component) Description() string {
	return fmt.Sprintf("%s %s %s", c.Make, strings.Split(strings.Fields(c.Model)[0], "/")[0], c.Type)
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

func (c ComponentList) Len() int           { return len(c) }
func (c ComponentList) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c ComponentList) Less(i, j int) bool { return c[i].Less(c[j]) }

func (c ComponentList) encode() [][]string {
	var data [][]string

	data = append(data, componentHeaders.Columns())

	for _, row := range c {
		data = append(data, []string{
			strings.TrimSpace(row.Make),
			strings.TrimSpace(row.Model),
			strings.TrimSpace(row.Type),
			strings.TrimSpace(row.number),
			strings.TrimSpace(row.Source),
			strings.TrimSpace(row.Subsource),
			strings.TrimSpace(row.dip),
			strings.TrimSpace(row.azimuth),
			strings.TrimSpace(row.Types),
			strings.TrimSpace(row.samplingRate),
			strings.TrimSpace(row.Response),
		})
	}

	return data
}

func (c *ComponentList) decode(data [][]string) error {
	if !(len(data) > 1) {
		return nil
	}

	var components []Component

	fields := componentHeaders.Fields(data[0])
	for _, row := range data[1:] {
		d := fields.Remap(row)

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
			Source:       strings.TrimSpace(d[componentSource]),
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

	*c = ComponentList(components)

	return nil
}

func LoadComponents(path string) ([]Component, error) {
	var c []Component

	if err := LoadList(path, (*ComponentList)(&c)); err != nil {
		return nil, err
	}

	sort.Sort(ComponentList(c))

	return c, nil
}
