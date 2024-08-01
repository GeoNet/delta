package main

import (
	stationxml "github.com/GeoNet/delta/internal/stationxml/v1.1"
)

func A2D(filter string, stage ResponseStage, count int, freq, rate float64) stationxml.ResponseStageType {

	coefs := stationxml.CoefficientsType{
		BaseFilterType: stationxml.BaseFilterType{
			ResourceId:  "Coefficients#" + filter,
			InputUnits:  stationxml.UnitsType{Name: stage.InputUnits},
			OutputUnits: stationxml.UnitsType{Name: stage.OutputUnits},
		},
		CfTransferFunctionType: stationxml.DigitalCfTransferFunction,
	}

	return stationxml.ResponseStageType{
		Number:       stationxml.CounterType(count),
		Coefficients: &coefs,
		StageGain: &stationxml.GainType{
			Value: func() float64 {
				if stage.Gain != 0.0 {
					return stage.Gain
				}
				return 1.0
			}(),
			Frequency: freq,
		},
		Decimation: &stationxml.DecimationType{
			InputSampleRate: stationxml.FrequencyType{FloatType: stationxml.FloatType{Value: func() float64 {
				if rate < 0.0 {
					return -1.0 / rate
				}
				return rate
			}(),
			}},
			Factor: func() int {
				if stage.Decimate != 0 {
					return int(stage.Decimate)
				}
				return 1
			}(),
			Delay:      stationxml.FloatType{Value: stage.Delay},
			Correction: stationxml.FloatType{Value: stage.Correction},
		},
	}
}
