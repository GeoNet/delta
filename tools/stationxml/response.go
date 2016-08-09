package main

import (
	"math"
	"math/cmplx"
	"strconv"

	"github.com/ozym/fdsn/stationxml"
)

//go:generate bash -c "go run generate/*.go | gofmt > responses_auto.go"

type Datalogger struct {
	Dataloggers   []string
	Type          string
	Label         string
	SampleRate    float64
	Frequency     float64
	StorageFormat string
	ClockDrift    float64
	Filters       []string
	Stages        [][]ResponseStage
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
	Stages   [][]ResponseStage
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

type StageSet interface {
	GetType() string
	ResponseStage(stage Stage) stationxml.ResponseStage
}

type ResponseStage struct {
	Type        string
	Lookup      string
	Filter      string
	StageSet    StageSet
	PAZ         PAZ
	FIR         FIR
	Polynomial  Polynomial
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

type PAZ struct {
	Name  string
	Code  stationxml.PzTransferFunction
	Type  string
	Notes string
	Poles []complex128
	Zeros []complex128
}

func (p PAZ) Gain(freq float64) float64 {
	w := complex(0.0, func() float64 {
		switch p.Code {
		case stationxml.PZFunctionLaplaceRadiansPerSecond:
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

func (p PAZ) GetType() string {
	return "paz"
}

type FIR struct {
	Name       string
	Causal     bool
	Symmetry   stationxml.Symmetry
	Decimation float64
	Gain       float64
	Notes      *string
	Factors    []float64
	Reversed   *bool
}

func (f FIR) GetType() string {
	return "fir"
}

func (f FIR) ResponseStage(stage Stage) stationxml.ResponseStage {
	var coeffs []stationxml.NumeratorCoefficient
	for j, c := range f.Factors {
		coeffs = append(coeffs, stationxml.NumeratorCoefficient{
			Coefficient: int32(j + 1),
			Value:       c,
		})
	}

	fir := stationxml.FIR{
		BaseFilter: stationxml.BaseFilter{
			ResourceId:  "FIR#" + stage.id,
			Name:        f.Name,
			InputUnits:  stationxml.Units{Name: stage.responseStage.InputUnits},
			OutputUnits: stationxml.Units{Name: stage.responseStage.OutputUnits},
		},
		Symmetry:              f.Symmetry,
		NumeratorCoefficients: coeffs,
	}

	return stationxml.ResponseStage{
		Number: stationxml.Counter(uint32(stage.count)),
		FIR:    &fir,
		//TODO: check we may need to adjust gain for different frequency
		StageGain: stationxml.Gain{
			Value: func() float64 {
				if stage.responseStage.Gain != 0.0 {
					return stage.responseStage.Gain
				}
				return 1.0
			}(),
			Frequency: stage.frequency,
		},
		Decimation: &stationxml.Decimation{
			InputSampleRate: stationxml.Frequency{stationxml.Float{Value: f.Decimation * stage.responseStage.SampleRate}},
			Factor: func() int32 {
				if stage.responseStage.Decimate != 0 {
					return stage.responseStage.Decimate
				}
				return int32(f.Decimation)
			}(),
			Delay:      stationxml.Float{Value: stage.responseStage.Delay},
			Correction: stationxml.Float{Value: stage.responseStage.Correction},
		},
	}
}

type Coefficient struct {
	Value float64
}

type Polynomial struct {
	Name                    string
	Gain                    float64
	ApproximationType       stationxml.ApproximationType
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

func (p Polynomial) ResponseStage(stage Stage) stationxml.ResponseStage {
	var coeffs []stationxml.Coefficient
	for n, c := range p.Coefficients {
		coeffs = append(coeffs, stationxml.Coefficient{
			Number: uint32(n) + 1,
			Value:  c.Value,
		})
	}

	poly := stationxml.Polynomial{
		BaseFilter: stationxml.BaseFilter{
			ResourceId:  "Polynomial#" + stage.id,
			Name:        stage.name,
			InputUnits:  stationxml.Units{Name: stage.responseStage.InputUnits},
			OutputUnits: stationxml.Units{Name: stage.responseStage.OutputUnits},
		},
		ApproximationType:       p.ApproximationType,
		FrequencyLowerBound:     stationxml.Frequency{stationxml.Float{Value: p.FrequencyLowerBound}},
		FrequencyUpperBound:     stationxml.Frequency{stationxml.Float{Value: p.FrequencyUpperBound}},
		ApproximationLowerBound: strconv.FormatFloat(p.ApproximationLowerBound, 'g', -1, 64),
		ApproximationUpperBound: strconv.FormatFloat(p.ApproximationUpperBound, 'g', -1, 64),
		MaximumError:            p.MaximumError,
		Coefficients:            coeffs,
	}

	return stationxml.ResponseStage{
		Number:     stationxml.Counter(uint32(stage.count)),
		Polynomial: &poly,
		//TODO: check we may need to adjust gain for different frequency
		StageGain: stationxml.Gain{
			Value: func() float64 {
				if p.Gain != 0.0 {
					return p.Gain
				}
				return 1.0
			}(),
			//Frequency: stage.Frequency,
			Frequency: stage.frequency,
		},
	}
}

func (p PAZ) ResponseStage(stage Stage) stationxml.ResponseStage {

	var poles []stationxml.PoleZero
	for j, pz := range p.Poles {
		poles = append(poles, stationxml.PoleZero{
			Number:    uint32(j),
			Real:      stationxml.FloatNoUnit{Value: real(pz)},
			Imaginary: stationxml.FloatNoUnit{Value: imag(pz)},
		})
	}
	var zeros []stationxml.PoleZero
	for j, pz := range p.Zeros {
		zeros = append(zeros, stationxml.PoleZero{
			Number:    uint32(len(p.Poles) + j),
			Real:      stationxml.FloatNoUnit{Value: real(pz)},
			Imaginary: stationxml.FloatNoUnit{Value: imag(pz)},
		})
	}

	paz := stationxml.PolesZeros{
		BaseFilter: stationxml.BaseFilter{
			ResourceId:  "PolesZeros#" + stage.id,
			Name:        stage.name,
			InputUnits:  stationxml.Units{Name: stage.responseStage.InputUnits},
			OutputUnits: stationxml.Units{Name: stage.responseStage.OutputUnits},
		},
		PzTransferFunction: p.Code,
		NormalizationFactor: func() float64 {
			return 1.0 / p.Gain(stage.frequency)
		}(),
		NormalizationFrequency: stationxml.Frequency{
			stationxml.Float{Value: stage.frequency},
		},
		Zeros: zeros,
		Poles: poles,
	}

	return stationxml.ResponseStage{
		Number:     stationxml.Counter(uint32(stage.count)),
		PolesZeros: &paz,
		StageGain: stationxml.Gain{
			Value: func() float64 {
				if stage.responseStage.Gain != 0.0 {
					return p.Gain(stage.frequency) * stage.responseStage.Gain / p.Gain(stage.responseStage.Frequency)
				}
				return 1.0
			}(),
			Frequency: stage.frequency,
		},
	}

}

type A2D struct {
	Name  string
	Code  stationxml.PzTransferFunction
	Type  string
	Notes string
}

func (a A2D) GetType() string {
	return "a2d"
}

func (a A2D) ResponseStage(stage Stage) stationxml.ResponseStage {

	coefs := stationxml.Coefficients{
		BaseFilter: stationxml.BaseFilter{
			ResourceId:  "Coefficients#" + stage.id,
			Name:        stage.name,
			InputUnits:  stationxml.Units{Name: stage.responseStage.InputUnits},
			OutputUnits: stationxml.Units{Name: stage.responseStage.OutputUnits},
		},
		CfTransferFunctionType: stationxml.CfFunctionDigital,
	}

	return stationxml.ResponseStage{
		Number:       stationxml.Counter(uint32(stage.count)),
		Coefficients: &coefs,
		StageGain: stationxml.Gain{
			Value: func() float64 {
				if stage.responseStage.Gain != 0.0 {
					return stage.responseStage.Gain
				}
				return 1.0
			}(),
			//Frequency: s.Frequency,
			Frequency: stage.frequency,
		},
		Decimation: &stationxml.Decimation{
			InputSampleRate: stationxml.Frequency{stationxml.Float{Value: stage.responseStage.SampleRate}},
			Factor: func() int32 {
				if stage.responseStage.Decimate != 0 {
					return stage.responseStage.Decimate
				}
				return 1
			}(),
			Delay:      stationxml.Float{Value: stage.responseStage.Delay},
			Correction: stationxml.Float{Value: stage.responseStage.Correction},
		},
	}
}
