package main

import (
	"strconv"
	//	"strings"

	"github.com/ozym/fdsn/stationxml"
)

type Stage struct {
	responseStage ResponseStage
	count         int
	id            string
	name          string
	frequency     float64
}

type Stager interface {
	ResponseStage(stage ResponseStage, count int, id, name string, frequency float64) stationxml.ResponseStage
}

func (s Stage) a2d() stationxml.ResponseStage {

	coefs := stationxml.Coefficients{
		BaseFilter: stationxml.BaseFilter{
			ResourceId:  "Coefficients#" + s.id,
			Name:        s.name,
			InputUnits:  stationxml.Units{Name: s.responseStage.InputUnits},
			OutputUnits: stationxml.Units{Name: s.responseStage.OutputUnits},
		},
		CfTransferFunctionType: stationxml.CfFunctionDigital,
	}

	return stationxml.ResponseStage{
		Number:       stationxml.Counter(uint32(s.count)),
		Coefficients: &coefs,
		StageGain: stationxml.Gain{
			Value: func() float64 {
				if s.responseStage.Gain != 0.0 {
					return s.responseStage.Gain
				}
				return 1.0
			}(),
			//Frequency: s.Frequency,
			Frequency: s.frequency,
		},
		Decimation: &stationxml.Decimation{
			InputSampleRate: stationxml.Frequency{stationxml.Float{Value: s.responseStage.SampleRate}},
			Factor: func() int32 {
				if s.responseStage.Decimate != 0 {
					return s.responseStage.Decimate
				}
				return 1
			}(),
			Delay:      stationxml.Float{Value: s.responseStage.Delay},
			Correction: stationxml.Float{Value: s.responseStage.Correction},
		},
	}
}

func (s Stage) paz(pz PAZ) stationxml.ResponseStage {

	var poles []stationxml.PoleZero
	for j, p := range pz.Poles {
		poles = append(poles, stationxml.PoleZero{
			Number:    uint32(j),
			Real:      stationxml.FloatNoUnit{Value: real(p)},
			Imaginary: stationxml.FloatNoUnit{Value: imag(p)},
		})
	}
	var zeros []stationxml.PoleZero
	for j, z := range pz.Zeros {
		zeros = append(zeros, stationxml.PoleZero{
			Number:    uint32(len(pz.Poles) + j),
			Real:      stationxml.FloatNoUnit{Value: real(z)},
			Imaginary: stationxml.FloatNoUnit{Value: imag(z)},
		})
	}

	paz := stationxml.PolesZeros{
		BaseFilter: stationxml.BaseFilter{
			ResourceId:  "PolesZeros#" + s.id,
			Name:        s.name,
			InputUnits:  stationxml.Units{Name: s.responseStage.InputUnits},
			OutputUnits: stationxml.Units{Name: s.responseStage.OutputUnits},
		},
		PzTransferFunction: pz.Code,
		NormalizationFactor: func() float64 {
			return 1.0 / pz.Gain(s.frequency)
		}(),
		NormalizationFrequency: stationxml.Frequency{
			stationxml.Float{Value: s.frequency},
		},
		Zeros: zeros,
		Poles: poles,
	}

	return stationxml.ResponseStage{
		Number:     stationxml.Counter(uint32(s.count)),
		PolesZeros: &paz,
		StageGain: stationxml.Gain{
			Value: func() float64 {
				if s.responseStage.Gain != 0.0 {
					return pz.Gain(s.frequency) * s.responseStage.Gain / pz.Gain(s.responseStage.Frequency)
				}
				return 1.0
			}(),
			Frequency: s.frequency,
		},
	}

}

func (s Stage) polynomial(p Polynomial) stationxml.ResponseStage {

	var coeffs []stationxml.Coefficient
	for n, c := range p.Coefficients {
		coeffs = append(coeffs, stationxml.Coefficient{
			Number: uint32(n) + 1,
			Value:  c.Value,
		})
	}

	poly := stationxml.Polynomial{
		BaseFilter: stationxml.BaseFilter{
			ResourceId:  "Polynomial#" + s.id,
			Name:        s.name,
			InputUnits:  stationxml.Units{Name: s.responseStage.InputUnits},
			OutputUnits: stationxml.Units{Name: s.responseStage.OutputUnits},
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
		Number:     stationxml.Counter(uint32(s.count)),
		Polynomial: &poly,
		//TODO: check we may need to adjust gain for different frequency
		StageGain: stationxml.Gain{
			Value: func() float64 {
				if p.Gain != 0.0 {
					return p.Gain
				}
				return 1.0
			}(),
			//Frequency: s.Frequency,
			Frequency: s.frequency,
		},
	}
}

func (s Stage) fir(f FIR) stationxml.ResponseStage {

	var coeffs []stationxml.NumeratorCoefficient
	for j, c := range f.Factors {
		coeffs = append(coeffs, stationxml.NumeratorCoefficient{
			Coefficient: int32(j + 1),
			Value:       c,
		})
	}

	fir := stationxml.FIR{
		BaseFilter: stationxml.BaseFilter{
			ResourceId:  "FIR#" + s.id,
			Name:        s.responseStage.Lookup,
			InputUnits:  stationxml.Units{Name: s.responseStage.InputUnits},
			OutputUnits: stationxml.Units{Name: s.responseStage.OutputUnits},
		},
		Symmetry: f.Symmetry,
		/*
			func() stationxml.Symmetry {
				switch strings.ToUpper(f.Symmetry) {
				case "EVEN":
					return stationxml.SymmetryEven
				case "ODD":
					return stationxml.SymmetryOdd
				default:
					return stationxml.SymmetryNone
				}
			}(),
		*/
		NumeratorCoefficients: coeffs,
	}

	return stationxml.ResponseStage{
		Number: stationxml.Counter(uint32(s.count)),
		FIR:    &fir,
		//TODO: check we may need to adjust gain for different frequency
		StageGain: stationxml.Gain{
			Value: func() float64 {
				if s.responseStage.Gain != 0.0 {
					return s.responseStage.Gain
				}
				return 1.0
			}(),
			//Frequency: stage.Frequency,
			Frequency: s.frequency,
		},
		Decimation: &stationxml.Decimation{
			InputSampleRate: stationxml.Frequency{stationxml.Float{Value: f.Decimation * s.responseStage.SampleRate}},
			Factor: func() int32 {
				if s.responseStage.Decimate != 0 {
					return s.responseStage.Decimate
				}
				return int32(f.Decimation)
			}(),
			Delay:      stationxml.Float{Value: s.responseStage.Delay},
			Correction: stationxml.Float{Value: s.responseStage.Correction},
		},
	}
}
