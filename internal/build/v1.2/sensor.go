package build

import (
	"encoding/xml"

	"github.com/GeoNet/delta/internal/stationxml/v1.2"
)

//TODO: handle calibration etc
/*
func Sensor(gain, bias, freq float64, data []byte) (*stationxml.ResponseType, error) {

	var sensor stationxml.ResponseType
	if err := xml.Unmarshal(data, &sensor); err != nil {
		return nil, err
	}

	switch {
	case sensor.InstrumentSensitivity != nil:
		switch {
		case bias != 0.0:
			//units := s.InstrumentSensitivity.OutputUnits
		case gain != 0.0:
			//units := s.InstrumentSensitivity.OutputUnits
		}
	case sensor.InstrumentPolynomial != nil:
		if gain != 1.0 || bias != 0.0 {
			//units = s.InstrumentPolynomial.OutputUnits
		}
	default:
		return nil, nil
	}

	return &sensor, nil
}
*/

//TODO: handle calibration etc
func (r *Response) Sensor(gain, bias float64, data []byte) error {

	var sensor stationxml.ResponseType
	if err := xml.Unmarshal(data, &sensor); err != nil {
		return err
	}

	switch {
	case sensor.InstrumentSensitivity != nil:
		switch {
		case bias != 0.0:
			//units := s.InstrumentSensitivity.OutputUnits
		case gain != 0.0:
			//units := s.InstrumentSensitivity.OutputUnits
		}
	case sensor.InstrumentPolynomial != nil:
		if gain != 1.0 || bias != 0.0 {
			//units = s.InstrumentPolynomial.OutputUnits
		}
	default:
		return nil
	}

	r.sensor = &sensor

	return nil
}

/*
func Datalogger(preamp, freq float64, data []byte) (*stationxml.ResponseType, error) {

	var datalogger stationxml.ResponseType
	if err := xml.Unmarshal(data, &datalogger); err != nil {
		return nil, err
	}

	if datalogger.InstrumentSensitivity == nil {
		return nil, nil
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
				NormalizationFrequency: stationxml.FrequencyType{FloatType: stationxml.FloatType{Value: freq}},
			},
			StageGain: &stationxml.GainType{
				Value:     preamp,
				Frequency: freq,
			},
		}}, datalogger.Stage...)
	}

	return &datalogger, nil
}
*/

/*
func InstrumentScale(freq float64, stages ...stationxml.ResponseStageType) float64 {
	scale := 1.0
	for _, s := range stages {
		if s.StageGain == nil || s.StageGain.Value == 0.0 {
			continue
		}
		scale *= s.StageGain.Value
	}
	return scale
}
*/

func (r *Response) Scale() float64 {
	scale := 1.0
	for _, s := range r.stages {
		if s.StageGain == nil || s.StageGain.Value == 0.0 {
			continue
		}
		scale *= s.StageGain.Value
	}
	return scale
}

/*
func InstrumentPoly(stages ...stationxml.ResponseStageType) *stationxml.PolynomialType {
	for _, s := range stages {
		if s.Polynomial == nil {
			continue
		}
		return s.Polynomial
	}
	return nil
}
*/

func (r *Response) Poly() *stationxml.PolynomialType {
	for _, s := range r.stages {
		if s.Polynomial == nil {
			continue
		}
		return s.Polynomial
	}
	return nil
}

/*
func InstrumentCoeffs(freq float64, stages ...stationxml.ResponseStageType) []stationxml.Coefficient {
	var coeffs []stationxml.Coefficient
	if poly := InstrumentPoly(stages...); poly != nil {
		scale := InstrumentScale(freq, stages...)

		if len(poly.Coefficient) > 0 {
			coeffs = append(coeffs, poly.Coefficient[0])
		}

		if len(poly.Coefficient) > 1 {
			coeffs = append(coeffs, stationxml.Coefficient{
				Number: 1,
				FloatNoUnitType: stationxml.FloatNoUnitType{
					Value: poly.Coefficient[1].Value / scale,
				},
			})
		}
	}
	return coeffs
}
*/

func (r *Response) Coeffs() []stationxml.Coefficient {
	var coeffs []stationxml.Coefficient
	if poly := r.Poly(); poly != nil {
		if len(poly.Coefficient) > 0 {
			coeffs = append(coeffs, poly.Coefficient[0])
		}

		if len(poly.Coefficient) > 1 {
			if scale := r.Scale(); scale != 0.0 {
				coeffs = append(coeffs, stationxml.Coefficient{
					Number: 1,
					FloatNoUnitType: stationxml.FloatNoUnitType{
						Value: poly.Coefficient[1].Value / scale,
					},
				})
			}
		}
	}
	return coeffs
}

/*
func Resp(prefix, id string, freq, gain, bias, preamp float64, sensor, datalogger []byte) (*stationxml.ResponseType, error) {

        if len(sensor) == 0 {
                return nil, nil
        }

        var s stationxml.ResponseType
        if err := xml.Unmarshal(sensor, &s); err != nil {
                return nil, err
        }

        if s.InstrumentSensitivity == nil && s.InstrumentPolynomial == nil {
                return nil, nil
        }

        if len(datalogger) == 0 {
                return nil, nil
        }

        var d stationxml.ResponseType
        if err := xml.Unmarshal(datalogger, &d); err != nil {
                return nil, err
        }

        if d.InstrumentSensitivity == nil && d.InstrumentPolynomial == nil {
                return nil, nil
        }

        var stages []stationxml.ResponseStageType

        stages = append(stages, s.Stage...)

        // a simple gain stage may come after the sensor
        if gain != 0.0 && gain != 1.0 && bias == 0.0 {
                var units stationxml.UnitsType
                if s.InstrumentSensitivity != nil {
                        units = s.InstrumentSensitivity.OutputUnits
                }
                if s.InstrumentPolynomial != nil {
                        units = s.InstrumentPolynomial.OutputUnits
                }
                stages = append(stages, PreampLegacy(gain, freq, units))
        }
*/
