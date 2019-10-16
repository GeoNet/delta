package resp

//go:generate bash -c "go run generate/*.go | gofmt -s > auto.go; test -s auto.go || rm auto.go"

import (
	"math"
	"math/cmplx"
	"strings"
)

type Symmetry uint

const (
	SymmetryUnknown Symmetry = iota
	SymmetryNone
	SymmetryEven
	SymmetryOdd
)

type PzTransferFunction uint

const (
	PZFunctionUnknown PzTransferFunction = iota
	PZFunctionLaplaceRadiansPerSecond
	PZFunctionLaplaceHertz
	PZFunctionLaplaceZTransform
)

type ApproximationType uint

const (
	ApproximationTypeUnknown ApproximationType = iota
	ApproximationTypeMaclaurin
)

type Datalogger struct {
	DataloggerList []string
	Type           string
	Label          string
	SampleRate     float64
	Frequency      float64
	StorageFormat  string
	ClockDrift     float64
	FilterList     []string
	Stages         []ResponseStage
	Reversed       bool
}

type DataloggerModel struct {
	Name         string
	Type         string // FDSN StationXML Datalogger Type
	Description  string // FDSN StationXML Datalogger Description
	Manufacturer string // FDSN StationXML Datalogger Manufacturer
	Vendor       string // FDSN StationXML Datalogger Vendor
}

type Sensor struct {
	SensorList []string
	FilterList []string
	Stages     []ResponseStage
	Channels   string
	Reversed   bool
}

func (s Sensor) Labels(axial string) string {
	labels := s.Channels
	switch strings.ToUpper(axial) {
	case "TRUE", "Z12":
		labels = strings.Replace(labels, "N", "1", -1)
		labels = strings.Replace(labels, "E", "2", -1)
		labels = strings.Replace(labels, "Y", "1", -1)
		labels = strings.Replace(labels, "X", "2", -1)
	case "FALSE", "ZNE":
		labels = strings.Replace(labels, "1", "N", -1)
		labels = strings.Replace(labels, "2", "E", -1)
		labels = strings.Replace(labels, "Y", "N", -1)
		labels = strings.Replace(labels, "X", "E", -1)
	case "ZYX", "XYZ":
		labels = strings.Replace(labels, "N", "Y", -1)
		labels = strings.Replace(labels, "E", "X", -1)
		labels = strings.Replace(labels, "1", "Y", -1)
		labels = strings.Replace(labels, "2", "X", -1)
	}
	return labels
}

type SensorComponent struct {
	Azimuth float64
	Dip     float64
}

type SensorModel struct {
	Name         string
	Type         string // FDSN StationXML Sensor Type
	Description  string // FDSN StationXML Sensor Description
	Manufacturer string // FDSN StationXML Vendor Description
	Vendor       string // FDSN StationXML Vendor Description

	Components []SensorComponent
}

type Response struct {
	Name        string
	Sensors     []Sensor
	Dataloggers []Datalogger
}

type Stream struct {
	Datalogger
	Sensor

	Components []SensorComponent
}

func (s Stream) Channels(axial string) []string {
	var channels []string

	labels := s.Sensor.Labels(axial)
	if len(s.Components) < len(labels) && len(labels) > 0 {
		labels = labels[0:len(s.Components)]
	}

	for _, component := range labels {
		channels = append(channels, s.Datalogger.Label+string(component))
	}

	return channels
}

func (s Stream) Gain() float64 {
	var gain float64 = 1.0

	for _, stage := range append(s.Sensor.Stages, s.Datalogger.Stages...) {
		if stage.StageSet == nil {
			continue
		}
		switch stage.StageSet.GetType() {
		case "fir":
			gain *= stage.StageSet.(FIR).Gain
		default:
			gain *= stage.Gain
		}
	}

	return gain
}

type StageSet interface {
	GetType() string
}

type ResponseStage struct {
	Type       string
	Lookup     string
	Filter     string
	StageSet   StageSet
	Frequency  float64
	SampleRate float64
	Decimate   int32
	Gain       float64
	//	Scale       float64
	Correction  float64
	Delay       float64
	InputUnits  string
	OutputUnits string
}

type PAZ struct {
	Name  string
	Code  PzTransferFunction
	Type  string
	Notes string
	Poles []complex128
	Zeros []complex128
}

func (p PAZ) GetType() string {
	return "paz"
}

func (p PAZ) Gain(freq float64) float64 {
	w := complex(0.0, func() float64 {
		switch p.Code {
		case PZFunctionLaplaceRadiansPerSecond:
			return 2.0 * math.Pi * freq
		default:
			return freq
		}
	}())
	h := complex(float64(1.0), float64(0.0))

	for _, zero := range p.Zeros {
		h *= (w - zero)
	}

	for _, pole := range p.Poles {
		h /= (w - pole)
	}
	return cmplx.Abs(h)
}

type FIR struct {
	Name       string
	Causal     bool
	Symmetry   Symmetry
	Decimation float64
	Gain       float64
	Notes      *string
	Factors    []float64
	Reversed   *bool
}

func (f FIR) GetType() string {
	return "fir"
}

type Coefficient struct {
	Value float64
}

type Polynomial struct {
	Name                    string
	Gain                    float64
	ApproximationType       ApproximationType
	FrequencyLowerBound     float64
	FrequencyUpperBound     float64
	ApproximationLowerBound float64
	ApproximationUpperBound float64
	MaximumError            float64
	Notes                   *string

	Coefficients []Coefficient
}

func (p Polynomial) GetType() string {
	return "poly"
}

type A2D struct {
	Name  string
	Code  PzTransferFunction
	Type  string
	Notes string
}

func (a A2D) GetType() string {
	return "a2d"
}
