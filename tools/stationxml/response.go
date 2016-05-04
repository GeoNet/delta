package main

import (
	"math"
	"math/cmplx"
)

type Datalogger struct {
	Dataloggers   []string
	Type          string
	Label         string
	Rate          float64
	Frequency     float64
	StorageFormat string
	ClockDrift    float64
	Filters       []string
	Reversed      bool
	Match         string
	Skip          string
}

type DataloggerModel struct {
	Type         string // FDSN StationXML Datalogger Type
	Description  string // FDSN StationXML Datalogger Description
	Manufacturer string // FDSN StationXML Datalogger Manufacturer
	Vendor       string // FDSN StationXML Datalogger Vendor
}

type Sensor struct {
	Sensors  []string
	Filters  []string
	Channels string
	Reversed bool
	Match    string
	Skip     string
}

type SensorComponent struct {
	Azimuth float64
	Dip     float64
}

type SensorModel struct {
	Type         string // FDSN StationXML Sensor Type
	Description  string // FDSN StationXML Sensor Description
	Manufacturer string // FDSN StationXML Vendor Description
	Vendor       string // FDSN StationXML Vendor Description

	Components []SensorComponent
}

type Response struct {
	Sensors     []Sensor
	Dataloggers []Datalogger
}

type Stream struct {
	Datalogger
	Sensor
}

type ResponseStage struct {
	Type        string
	Lookup      string
	Frequency   float64
	SampleRate  float64
	Decimate    int32
	Gain        float64
	Scale       float64
	Correction  float64
	Delay       float64
	InputUnits  string
	OutputUnits string
}

type Filter struct {
	Stages []ResponseStage
}

type PAZ struct {
	Code  string
	Type  string
	Notes string
	Poles []complex128
	Zeros []complex128
}

func (p PAZ) Gain(freq float64) float64 {
	w := complex(0.0, func() float64 {
		switch p.Code {
		case "A":
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
	Causal     bool
	Symmetry   string
	Decimation float64
	Gain       float64
	Notes      *string
	Factors    []float64
	Reversed   *bool
}

type Coefficient struct {
	Value float64
}

type Polynomial struct {
	Gain                    float64
	ApproximationType       string
	FrequencyLowerBound     float64
	FrequencyUpperBound     float64
	ApproximationLowerBound float64
	ApproximationUpperBound float64
	MaximumError            float64
	Notes                   *string

	Coefficients []Coefficient
}
