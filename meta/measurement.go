package meta

import (
	"fmt"
	"sort"
	"strings"
)

const (
	measurementLocation = iota
	measurementName
	measurementSensor
	measurementType
	measurementUnit
	measurementDescription
	measurementLast
)

type Measurement struct {
	Location    string
	Name        string
	Sensor      string
	Type        string
	Unit        string
	Description string
}

type MeasurementList []Measurement

func (s MeasurementList) Len() int      { return len(s) }
func (s MeasurementList) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s MeasurementList) Less(i, j int) bool {
	switch {
	case s[i].Location < s[j].Location:
		return true
	case s[i].Location > s[j].Location:
		return false
	case s[i].Name < s[j].Name:
		return true
	default:
		return false
	}
}

func (s MeasurementList) encode() [][]string {
	data := [][]string{{
		"Location",
		"Name",
		"Sensor",
		"Type",
		"Unit",
		"Description",
	}}

	for _, v := range s {
		data = append(data, []string{
			strings.TrimSpace(v.Location),
			strings.TrimSpace(v.Name),
			strings.TrimSpace(v.Sensor),
			strings.TrimSpace(v.Type),
			strings.TrimSpace(v.Unit),
			strings.TrimSpace(v.Description),
		})
	}
	return data
}
func (s *MeasurementList) decode(data [][]string) error {
	var measurements []Measurement
	if len(data) > 1 {
		for _, d := range data[1:] {
			if len(d) != measurementLast {
				return fmt.Errorf("incorrect number of measurement fields")
			}

			measurements = append(measurements, Measurement{
				Location:    strings.TrimSpace(d[measurementLocation]),
				Name:        strings.TrimSpace(d[measurementName]),
				Sensor:      strings.TrimSpace(d[measurementSensor]),
				Type:        strings.TrimSpace(d[measurementType]),
				Unit:        strings.TrimSpace(d[measurementUnit]),
				Description: strings.TrimSpace(d[measurementDescription]),
			})
		}

		*s = MeasurementList(measurements)
	}
	return nil
}

func LoadMeasurements(path string) ([]Measurement, error) {
	var s []Measurement

	if err := LoadList(path, (*MeasurementList)(&s)); err != nil {
		return nil, err
	}

	sort.Sort(MeasurementList(s))

	return s, nil
}
