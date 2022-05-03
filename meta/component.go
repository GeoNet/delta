package meta

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

const (
	componentMake = iota
	componentModel
	componentType
	componentPin
	componentDip
	componentAzimuth
	componentResponse
	componentLast
)

type Component struct {
	Make     string
	Model    string
	Type     string
	Response string

	Pin     int
	Azimuth float64
	Dip     float64

	pin     string
	azimuth string
	dip     string
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
	case c.Pin < comp.Pin:
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
		"Pin",
		"Dip",
		"Azimuth",
		"Response",
	}}

	for _, v := range s {
		data = append(data, []string{
			strings.TrimSpace(v.Make),
			strings.TrimSpace(v.Model),
			strings.TrimSpace(v.Type),
			strings.TrimSpace(v.pin),
			strings.TrimSpace(v.dip),
			strings.TrimSpace(v.azimuth),
			strings.TrimSpace(v.Response),
		})
	}
	return data
}
func (s *ComponentList) decode(data [][]string) error {
	var components []Component
	if len(data) > 1 {
		for _, d := range data[1:] {
			if len(d) != componentLast {
				return fmt.Errorf("incorrect pin of installed component fields")
			}

			pin, err := strconv.Atoi(d[componentPin])
			if err != nil {
				return err
			}

			azimuth, err := strconv.ParseFloat(d[componentAzimuth], 64)
			if err != nil {
				return err
			}
			dip, err := strconv.ParseFloat(d[componentDip], 64)
			if err != nil {
				return err
			}

			components = append(components, Component{
				Make:     strings.TrimSpace(d[componentMake]),
				Model:    strings.TrimSpace(d[componentModel]),
				Type:     strings.TrimSpace(d[componentType]),
				Response: strings.TrimSpace(d[componentResponse]),
				Pin:      pin,
				Dip:      dip,
				Azimuth:  azimuth,

				pin:     strings.TrimSpace(d[componentPin]),
				dip:     strings.TrimSpace(d[componentDip]),
				azimuth: strings.TrimSpace(d[componentAzimuth]),
			})
		}

		*s = ComponentList(components)
	}
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
