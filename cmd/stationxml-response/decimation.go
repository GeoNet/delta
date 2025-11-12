package main

import (
	stationxml "github.com/GeoNet/delta/internal/stationxml/v1.2"
)

// DecimationFilter holds the stream decimation details and
// the name of the anti-alias FIR filter used.
type DecimationFilter struct {
	FirFilter       string  `json:"fir_filter"`
	InputSampleRate float64 `json:"input_sample_rate"`
	Factor          int     `json:"factor"`
	Offset          int     `json:"offset,omitzero"`
	Delay           float64 `json:"delay,omitzero"`
	Correction      float64 `json:"correction,omitzero"`
}

func GetDecimationFilter(stage stationxml.ResponseStageType) (DecimationFilter, bool) {
	if stage.Decimation == nil {
		return DecimationFilter{}, false
	}
	if stage.FIR == nil {
		return DecimationFilter{}, false
	}

	filter := DecimationFilter{
		FirFilter:       stage.FIR.Name,
		InputSampleRate: stage.Decimation.InputSampleRate.Value,
		Factor:          stage.Decimation.Factor,
		Offset:          stage.Decimation.Offset,
		Delay:           stage.Decimation.Delay.Value,
		Correction:      stage.Decimation.Correction.Value,
	}

	return filter, true
}
