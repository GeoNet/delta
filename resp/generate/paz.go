package main

import (
	"math"
	"math/cmplx"
	"sort"

	stationxml "github.com/GeoNet/delta/internal/stationxml/v1.1"
)

type PAZ struct {
	Code  string       `yaml:"code"`
	Type  string       `yaml:"type"`
	Notes string       `yaml:"notes"`
	Poles []Complex128 `yaml:"poles"`
	Zeros []Complex128 `yaml:"zeros"`
}

func (p PAZ) PzTransferFunction() string {
	switch p.Code {
	case "A":
		return "PZFunctionLaplaceRadiansPerSecond"
	case "B":
		return "PZFunctionLaplaceHertz"
	case "D":
		return "PZFunctionLaplaceZTransform"
	default:
		return "PZFunctionUnknown"
	}
}

func (p PAZ) PzTransferType() stationxml.PzTransferFunctionType {
	switch p.Code {
	case "A":
		return stationxml.LaplaceRadiansSecondPzTransferFunction
	case "B":
		return stationxml.LaplaceHertzPzTransferFunction
	case "D":
		return stationxml.DigitalZTransformPzTransferFunction
	default:
		return stationxml.PzTransferFunctionType(0)
	}
}

func (p PAZ) Gain(freq float64) float64 {
	w := complex(0.0, func() float64 {
		switch p.PzTransferType() {
		case stationxml.LaplaceRadiansSecondPzTransferFunction:
			return 2.0 * math.Pi * freq
		default:
			return freq
		}
	}())
	h := complex(float64(1.0), float64(0.0))

	for _, zero := range p.Zeros {
		h *= (w - complex128(zero))
	}

	for _, pole := range p.Poles {
		h /= (w - complex128(pole))
	}
	return cmplx.Abs(h)
}

type pazMap map[string]PAZ

func (p pazMap) Keys() []string {
	var keys []string
	for k := range p {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (p pazMap) Values() []PAZ {
	var values []PAZ
	for _, k := range p.Keys() {
		values = append(values, p[k])
	}
	return values
}

func (p PAZ) Stage(filter string, stage ResponseStage, count int, freq, gain float64) stationxml.ResponseStageType {

	var poles []stationxml.PoleZeroType
	for j, pz := range p.Poles {
		poles = append(poles, stationxml.PoleZeroType{
			Number:    j,
			Real:      stationxml.FloatNoUnitType{Value: real(pz)},
			Imaginary: stationxml.FloatNoUnitType{Value: imag(pz)},
		})
	}

	var zeros []stationxml.PoleZeroType
	for j, pz := range p.Zeros {
		zeros = append(zeros, stationxml.PoleZeroType{
			Number:    len(p.Poles) + j,
			Real:      stationxml.FloatNoUnitType{Value: real(pz)},
			Imaginary: stationxml.FloatNoUnitType{Value: imag(pz)},
		})
	}

	return stationxml.ResponseStageType{
		Number: stationxml.CounterType(count),
		PolesZeros: &stationxml.PolesZerosType{
			BaseFilterType: stationxml.BaseFilterType{
				ResourceId:  "PolesZeros#" + filter,
				InputUnits:  stationxml.UnitsType{Name: stage.InputUnits},
				OutputUnits: stationxml.UnitsType{Name: stage.OutputUnits},
			},
			PzTransferFunctionType: p.PzTransferType(),
			NormalizationFactor: func() float64 {
				return 1.0 / p.Gain(freq)
			}(),
			NormalizationFrequency: stationxml.FrequencyType{
				FloatType: stationxml.FloatType{Value: freq},
			},
			Zero: zeros,
			Pole: poles,
		},
		StageGain: &stationxml.GainType{
			Value: func() float64 {
				if gain != 0.0 {
					return gain
				}
				return 1.0
			}(),
			Frequency: freq,
		},
	}
}

func (p PAZ) Sensitivity(stage ResponseStage) stationxml.SensitivityType {
	return stationxml.SensitivityType{
		GainType: stationxml.GainType{
			Value: func() float64 {
				if stage.Gain != 0.0 {
					return stage.Gain
				}
				return p.Gain(stage.Frequency)
			}(),
			Frequency: stage.Frequency,
		},
		InputUnits:  stationxml.UnitsType{Name: stage.InputUnits},
		OutputUnits: stationxml.UnitsType{Name: stage.OutputUnits},
	}
}
