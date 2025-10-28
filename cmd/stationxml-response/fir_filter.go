package main

import stationxml "github.com/GeoNet/delta/internal/stationxml/v1.2"

// FirFilter holds the details of a datalogger FIR filter.
type FirFilter struct {
	Name         string    `json:"name"`
	Symmetry     string    `json:"symmetry"`
	Length       int       `json:"length"`
	Coefficients []float64 `json:"coefficients"`
}

func GetFirFilter(stage stationxml.ResponseStageType) (FirFilter, bool) {
	if stage.FIR == nil {
		return FirFilter{}, false
	}

	var coeffs []float64
	for _, v := range stage.FIR.NumeratorCoefficient {
		coeffs = append(coeffs, v.Value)
	}

	fir := FirFilter{
		Name:         stage.FIR.Name,
		Symmetry:     stage.FIR.Symmetry.String(),
		Length:       len(coeffs),
		Coefficients: coeffs,
	}

	return fir, true
}
