package meta

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/GeoNet/delta/internal/expr"
)

const (
	calibrationMake = iota
	calibrationModel
	calibrationSerial
	calibrationNumber
	calibrationScaleFactor
	calibrationScaleBias
	calibrationFrequency
	calibrationStart
	calibrationEnd
	calibrationLast
)

// Calibration defines times where sensor scaling or offsets are needed, these will be overwrite the
// existing values, i.e. A + BX => A' + B' X, where A' and B' are the given bias and scaling factors.
type Calibration struct {
	Install

	ScaleFactor float64
	ScaleBias   float64
	Frequency   float64
	Number      int

	factor    string
	bias      string
	frequency string
	number    string
}

// Id returns a unique string which can be used for sorting or checking.
func (c Calibration) Id() string {
	return strings.Join([]string{c.Make, c.Model, c.Serial, strconv.Itoa(c.Number)}, ":")
}

// Less returns whether one Calibration sorts before another.
func (s Calibration) Less(calibration Calibration) bool {
	switch {
	case s.Install.Less(calibration.Install):
		return true
	case calibration.Install.Less(s.Install):
		return false
	case s.Number < calibration.Number:
		return true
	default:
		return false
	}
}

// CalibrationList is a slice of Calibration types.
type CalibrationList []Calibration

func (c CalibrationList) Len() int           { return len(c) }
func (c CalibrationList) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c CalibrationList) Less(i, j int) bool { return c[i].Less(c[j]) }

func (c CalibrationList) encode() [][]string {
	data := [][]string{{
		"Make",
		"Model",
		"Serial",
		"Number",
		"Scale Factor",
		"Scale Bias",
		"Frequency",
		"Start Date",
		"End Date",
	}}

	for _, v := range c {
		data = append(data, []string{
			strings.TrimSpace(v.Make),
			strings.TrimSpace(v.Model),
			strings.TrimSpace(v.Serial),
			strconv.Itoa(v.Number),
			strings.TrimSpace(v.factor),
			strings.TrimSpace(v.bias),
			strings.TrimSpace(v.frequency),
			v.Start.Format(DateTimeFormat),
			v.End.Format(DateTimeFormat),
		})
	}
	return data
}

func (c *CalibrationList) toFloat64(str string, def float64) (float64, error) {
	switch s := strings.TrimSpace(str); {
	case s != "":
		return expr.ToFloat64(s)
	default:
		return def, nil
	}
}

func (c *CalibrationList) toInt(str string, def int) (int, error) {
	switch s := strings.TrimSpace(str); {
	case s != "":
		return expr.ToInt(s)
	default:
		return def, nil
	}
}

func (c *CalibrationList) decode(data [][]string) error {
	var calibrations []Calibration
	if len(data) > 1 {
		for _, d := range data[1:] {
			if len(d) != calibrationLast {
				return fmt.Errorf("incorrect number of installed calibration fields")
			}

			factor, err := c.toFloat64(d[calibrationScaleFactor], 1.0)
			if err != nil {
				return err
			}

			bias, err := c.toFloat64(d[calibrationScaleBias], 0.0)
			if err != nil {
				return err
			}

			freq, err := c.toFloat64(d[calibrationFrequency], 0.0)
			if err != nil {
				return err
			}

			number, err := c.toInt(d[calibrationNumber], 0)
			if err != nil {
				return err
			}

			start, err := time.Parse(DateTimeFormat, d[calibrationStart])
			if err != nil {
				return err
			}

			end, err := time.Parse(DateTimeFormat, d[calibrationEnd])
			if err != nil {
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
				Number: number,

				ScaleFactor: factor,
				ScaleBias:   bias,
				Frequency:   freq,

				number:    strings.TrimSpace(d[calibrationNumber]),
				factor:    strings.TrimSpace(d[calibrationScaleFactor]),
				bias:      strings.TrimSpace(d[calibrationScaleBias]),
				frequency: strings.TrimSpace(d[calibrationFrequency]),
			})
		}
	}

	*c = CalibrationList(calibrations)

	return nil
}

// LoadCalibrations reads a CSV formatted file and returns a slice of Calibration types.
func LoadCalibrations(path string) ([]Calibration, error) {
	var c []Calibration

	if err := LoadList(path, (*CalibrationList)(&c)); err != nil {
		return nil, err
	}

	sort.Sort(CalibrationList(c))

	return c, nil
}
