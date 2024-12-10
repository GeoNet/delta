package meta

import (
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
	calibrationScaleAbsolute
	calibrationFrequency
	calibrationStart
	calibrationEnd
	calibrationLast
)

var calibrationHeaders Header = map[string]int{
	"Make":           calibrationMake,
	"Model":          calibrationModel,
	"Serial":         calibrationSerial,
	"Number":         calibrationNumber,
	"Scale Factor":   calibrationScaleFactor,
	"Scale Bias":     calibrationScaleBias,
	"Scale Absolute": calibrationScaleAbsolute,
	"Frequency":      calibrationFrequency,
	"Start Date":     calibrationStart,
	"End Date":       calibrationEnd,
}

var CalibrationTable Table = Table{
	name:    "Calibration",
	headers: calibrationHeaders,
	primary: []string{"Make", "Model", "Serial", "Number", "Start Date"},
	native:  []string{"Scale Factor", "Scale Bias", "Scale Absolute", "Frequency"},
	foreign: map[string][]string{
		"Asset": {"Make", "Model", "Serial"},
	},
	remap: map[string]string{
		"Scale Factor":   "ScaleFactor",
		"Scale Bias":     "ScaleBias",
		"Scale Absolute": "ScaleAbsolute",
		"Start Date":     "Start",
		"End Date":       "End",
	},
	start: "Start Date",
	end:   "End Date",
}

// Calibration defines times where sensor scaling or offsets are needed, these will be overwrite the
// existing values, i.e. A + BX => A' + B' X, where A' and B' are the given bias and scaling factors.
type Calibration struct {
	Install

	ScaleFactor   float64 `json:"scale-factor"`
	ScaleBias     float64 `json:"scale-bias"`
	ScaleAbsolute float64 `json:"scale-absolute"`
	Frequency     float64 `json:"scale-frequency"`
	Number        int     `json:"number"`

	factor    string
	bias      string
	absolute  string
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

	var data [][]string

	data = append(data, calibrationHeaders.Columns())

	for _, row := range c {
		data = append(data, []string{
			strings.TrimSpace(row.Make),
			strings.TrimSpace(row.Model),
			strings.TrimSpace(row.Serial),
			strconv.Itoa(row.Number),
			strings.TrimSpace(row.factor),
			strings.TrimSpace(row.bias),
			strings.TrimSpace(row.absolute),
			strings.TrimSpace(row.frequency),
			row.Start.Format(DateTimeFormat),
			row.End.Format(DateTimeFormat),
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
	if !(len(data) > 1) {
		return nil
	}

	var calibrations []Calibration

	fields := calibrationHeaders.Fields(data[0])
	for _, row := range data[1:] {
		d := fields.Remap(row)

		factor, err := c.toFloat64(d[calibrationScaleFactor], 1.0)
		if err != nil {
			return err
		}

		bias, err := c.toFloat64(d[calibrationScaleBias], 1.0)
		if err != nil {
			return err
		}

		absolute, err := c.toFloat64(d[calibrationScaleAbsolute], 0.0)
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

			ScaleFactor:   factor,
			ScaleBias:     bias,
			ScaleAbsolute: absolute,
			Frequency:     freq,

			number:    strings.TrimSpace(d[calibrationNumber]),
			factor:    strings.TrimSpace(d[calibrationScaleFactor]),
			bias:      strings.TrimSpace(d[calibrationScaleBias]),
			absolute:  strings.TrimSpace(d[calibrationScaleAbsolute]),
			frequency: strings.TrimSpace(d[calibrationFrequency]),
		})
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
