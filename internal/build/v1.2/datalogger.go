package build

import (
	"encoding/xml"

	"github.com/GeoNet/delta/internal/stationxml/v1.2"
)

func (r *Response) Datalogger(preamp float64, data []byte) error {

	var datalogger stationxml.ResponseType
	if err := xml.Unmarshal(data, &datalogger); err != nil {
		return err
	}

	if datalogger.InstrumentSensitivity == nil {
		return nil
	}

	if preamp != 0.0 {
		datalogger.Stage = append([]stationxml.ResponseStageType{{
			//TODO: technically the poles and zeros are not required, but kept to allow acceptance checks
			PolesZeros: &stationxml.PolesZerosType{
				BaseFilterType: stationxml.BaseFilterType{
					InputUnits:  datalogger.InstrumentSensitivity.InputUnits,
					OutputUnits: datalogger.InstrumentSensitivity.InputUnits,
				},
				PzTransferFunctionType: stationxml.LaplaceRadiansSecondPzTransferFunction,
				NormalizationFactor:    1.0,
				NormalizationFrequency: stationxml.FrequencyType{FloatType: stationxml.FloatType{Value: r.freq}},
			},
			StageGain: &stationxml.GainType{
				Value:     preamp,
				Frequency: r.freq,
			},
		}}, datalogger.Stage...)
	}

	r.datalogger = &datalogger

	return nil
}
