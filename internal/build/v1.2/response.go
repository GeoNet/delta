package build

import (
	"strings"

	"github.com/GeoNet/delta/internal/stationxml/v1.2"
)

type Response struct {
	prefix string
	id     string
	freq   float64

	sensor     *stationxml.ResponseType
	datalogger *stationxml.ResponseType

	stages []stationxml.ResponseStageType
}

func NewResponse(prefix, id string, freq float64) *Response {
	return &Response{
		prefix: prefix,
		id:     id,
		freq:   freq,
	}
}

/*
func (r *Response) Stages() []stationxml.ResponseStageType {
	var stages []stationxml.ResponseStageType
	stages = append(stages, r.sensor.Stage...)
	stages = append(stages, r.datalogger.Stage...)
	return stages
}
*/

func (r *Response) ResponseType() (*stationxml.ResponseType, error) {

	if err := r.Normalise(); err != nil {
		return nil, err
	}

	resp := stationxml.ResponseType{
		InstrumentSensitivity: func() *stationxml.SensitivityType {
			if r.sensor.InstrumentSensitivity != nil {
				return &stationxml.SensitivityType{
					InputUnits:  r.sensor.InstrumentSensitivity.InputUnits,
					OutputUnits: r.datalogger.InstrumentSensitivity.OutputUnits,
					GainType: stationxml.GainType{
						Frequency: r.freq,
						Value:     r.Scale(),
					},
				}
			}
			return nil
		}(),
		InstrumentPolynomial: func() *stationxml.PolynomialType {
			if poly := r.Poly(); poly != nil {
				return &stationxml.PolynomialType{
					BaseFilterType: stationxml.BaseFilterType{
						ResourceId:  "Instrument" + poly.ResourceId + ":" + r.id,
						Name:        strings.TrimRight(r.prefix, "."),
						InputUnits:  poly.InputUnits,
						OutputUnits: r.datalogger.InstrumentSensitivity.OutputUnits,
					},
					ApproximationType:       poly.ApproximationType,
					FrequencyLowerBound:     poly.FrequencyLowerBound,
					FrequencyUpperBound:     poly.FrequencyUpperBound,
					ApproximationLowerBound: poly.ApproximationLowerBound,
					ApproximationUpperBound: poly.ApproximationUpperBound,
					MaximumError:            poly.MaximumError,
					Coefficient:             r.Coeffs(),
				}
			}
			return nil
		}(),
		Stage: r.stages,
	}

	return &resp, nil
}
