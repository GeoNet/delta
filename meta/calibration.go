package meta

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	calibrationMake = iota
	calibrationModel
	calibrationSerial
	calibrationComponent
	calibrationScaleFactor
	calibrationScaleBias
	calibrationStart
	calibrationEnd
	calibrationLast
)

type Calibration struct {
	Install
	Scale

	Component int
}

type CalibrationList []Calibration

func (s CalibrationList) Len() int           { return len(s) }
func (s CalibrationList) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s CalibrationList) Less(i, j int) bool { return s[i].Install.Less(s[j].Install) }

func (s CalibrationList) encode() [][]string {
	data := [][]string{{
		"Make",
		"Model",
		"Serial",
		"Component",
		"Scale Factor",
		"Scale Bias",
		"Start Date",
		"End Date",
	}}

	for _, v := range s {
		data = append(data, []string{
			strings.TrimSpace(v.Make),
			strings.TrimSpace(v.Model),
			strings.TrimSpace(v.Serial),
			strconv.Itoa(v.Component),
			strings.TrimSpace(v.factor),
			strings.TrimSpace(v.bias),
			v.Start.Format(DateTimeFormat),
			v.End.Format(DateTimeFormat),
		})
	}
	return data
}
func (s *CalibrationList) decode(data [][]string) error {
	var calibrations []Calibration
	if len(data) > 1 {
		for _, d := range data[1:] {
			if len(d) != calibrationLast {
				return fmt.Errorf("incorrect number of installed calibration fields")
			}
			var err error

			var component int
			if component, err = strconv.Atoi(d[calibrationComponent]); err != nil {
				return err
			}

			var factor, bias float64
			if factor, err = strconv.ParseFloat(d[calibrationScaleFactor], 64); err != nil {
				return err
			}
			if bias, err = strconv.ParseFloat(d[calibrationScaleBias], 64); err != nil {
				return err
			}

			var start, end time.Time
			if start, err = time.Parse(DateTimeFormat, d[calibrationStart]); err != nil {
				return err
			}
			if end, err = time.Parse(DateTimeFormat, d[calibrationEnd]); err != nil {
				return err
			}

			calibrations = append(calibrations, Calibration{
				Install: Install{
					Equipment: Equipment{
						Make:   strings.TrimSpace(d[calibrationMake]),
						Model:  strings.TrimSpace(d[calibrationModel]),
						Serial: strings.TrimSpace(d[calibrationSerial]),
					},
					Span: Span{
						Start: start,
						End:   end,
					},
				},
				Scale: Scale{
					Factor: factor,
					Bias:   bias,

					factor: strings.TrimSpace(d[calibrationScaleFactor]),
					bias:   strings.TrimSpace(d[calibrationScaleBias]),
				},
				Component: component,
			})
		}

		*s = CalibrationList(calibrations)
	}
	return nil
}

func LoadCalibrations(path string) ([]Calibration, error) {
	var s []Calibration

	if err := LoadList(path, (*CalibrationList)(&s)); err != nil {
		return nil, err
	}

	sort.Sort(CalibrationList(s))

	return s, nil
}
