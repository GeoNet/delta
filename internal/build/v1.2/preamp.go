package build

import (
	"github.com/GeoNet/delta/internal/stationxml/v1.2"
)

// Not strictly needed, but used to match existing pre-amp gains
func PreampLegacy(resp *stationxml.ResponseType, gain, freq float64) stationxml.ResponseStageType {
	var units stationxml.UnitsType
	if resp.InstrumentSensitivity != nil {
		units = resp.InstrumentSensitivity.OutputUnits
	}
	if resp.InstrumentPolynomial != nil {
		units = resp.InstrumentPolynomial.OutputUnits
	}
	return stationxml.ResponseStageType{
		PolesZeros: &stationxml.PolesZerosType{
			BaseFilterType: stationxml.BaseFilterType{
				InputUnits:  units,
				OutputUnits: units,
			},
			PzTransferFunctionType: stationxml.LaplaceRadiansSecondPzTransferFunction,
			NormalizationFactor:    1.0,
			NormalizationFrequency: stationxml.FrequencyType{FloatType: stationxml.FloatType{Value: freq}},
		},
		StageGain: &stationxml.GainType{
			Value:     gain,
			Frequency: freq,
		},
	}
}

func Preamp(gain, freq float64) stationxml.ResponseStageType {
	return stationxml.ResponseStageType{
		StageGain: &stationxml.GainType{
			Value:     gain,
			Frequency: freq,
		},
	}
}
