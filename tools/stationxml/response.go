package main

import (
	"strconv"

	"github.com/GeoNet/delta/resp"
	"github.com/ozym/fdsn/stationxml"
)

type Stage struct {
	responseStage resp.ResponseStage
	count         int
	id            string
	name          string
	frequency     float64
}

func a2dResponseStage(stage Stage) stationxml.ResponseStage {

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
			InputSampleRate: stationxml.Frequency{Float: stationxml.Float{Value: stage.responseStage.SampleRate}},
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

func firResponseStage(filter resp.FIR, stage Stage) stationxml.ResponseStage {
	var coeffs []stationxml.NumeratorCoefficient
	for j, c := range filter.Factors {
		coeffs = append(coeffs, stationxml.NumeratorCoefficient{
			Coefficient: int32(j + 1),
			Value:       c,
		})
	}

	fir := stationxml.FIR{
		BaseFilter: stationxml.BaseFilter{
			ResourceId:  "FIR#" + stage.id,
			Name:        filter.Name,
			InputUnits:  stationxml.Units{Name: stage.responseStage.InputUnits},
			OutputUnits: stationxml.Units{Name: stage.responseStage.OutputUnits},
		},
		Symmetry:              stationxml.Symmetry(filter.Symmetry),
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
			InputSampleRate: stationxml.Frequency{Float: stationxml.Float{Value: filter.Decimation * stage.responseStage.SampleRate}},
			Factor: func() int32 {
				if stage.responseStage.Decimate != 0 {
					return stage.responseStage.Decimate
				}
				return int32(filter.Decimation)
			}(),
			//Best leave this to the application as there are multiple interpretations.
			Delay:      stationxml.Float{Value: 0.0},
			Correction: stationxml.Float{Value: 0.0},
		},
	}
}

func polyResponseStage(filter resp.Polynomial, stage Stage) stationxml.ResponseStage {
	var coeffs []stationxml.Coefficient
	for n, c := range filter.Coefficients {
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
		ApproximationType:       stationxml.ApproximationType(filter.ApproximationType),
		FrequencyLowerBound:     stationxml.Frequency{Float: stationxml.Float{Value: filter.FrequencyLowerBound}},
		FrequencyUpperBound:     stationxml.Frequency{Float: stationxml.Float{Value: filter.FrequencyUpperBound}},
		ApproximationLowerBound: strconv.FormatFloat(filter.ApproximationLowerBound, 'g', -1, 64),
		ApproximationUpperBound: strconv.FormatFloat(filter.ApproximationUpperBound, 'g', -1, 64),
		MaximumError:            filter.MaximumError,
		Coefficients:            coeffs,
	}

	return stationxml.ResponseStage{
		Number:     stationxml.Counter(uint32(stage.count)),
		Polynomial: &poly,
		//TODO: check we may need to adjust gain for different frequency
		StageGain: stationxml.Gain{
			Value: func() float64 {
				if filter.Gain != 0.0 {
					return filter.Gain
				}
				return 1.0
			}(),
			//Frequency: stage.Frequency,
			Frequency: stage.frequency,
		},
	}
}

func pazResponseStage(filter resp.PAZ, stage Stage) stationxml.ResponseStage {
	var poles []stationxml.PoleZero
	for j, pz := range filter.Poles {
		poles = append(poles, stationxml.PoleZero{
			Number:    uint32(j),
			Real:      stationxml.FloatNoUnit{Value: real(pz)},
			Imaginary: stationxml.FloatNoUnit{Value: imag(pz)},
		})
	}
	var zeros []stationxml.PoleZero
	for j, pz := range filter.Zeros {
		zeros = append(zeros, stationxml.PoleZero{
			Number:    uint32(len(filter.Poles) + j),
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
		PzTransferFunctionType: stationxml.PzTransferFunctionType(filter.Code),
		NormalizationFactor: func() float64 {
			return 1.0 / filter.Gain(stage.frequency)
		}(),
		NormalizationFrequency: stationxml.Frequency{
			Float: stationxml.Float{Value: stage.frequency},
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
					return filter.Gain(stage.frequency) * stage.responseStage.Gain / filter.Gain(stage.responseStage.Frequency)
				}
				return 1.0
			}(),
			Frequency: stage.frequency,
		},
	}

}
