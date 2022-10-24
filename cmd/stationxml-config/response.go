package main

import (
	"bytes"
	"encoding/gob"
	"encoding/xml"
	"fmt"
	"math"
	"math/cmplx"
	"strings"

	"github.com/GeoNet/delta/internal/stationxml/v1.2"
)

func clone(a, b *stationxml.ResponseStageType) error {

	var buff bytes.Buffer

	if err := gob.NewEncoder(&buff).Encode(a); err != nil {
		return err
	}
	if err := gob.NewDecoder(&buff).Decode(b); err != nil {
		return err
	}

	return nil
}

func gain(pz *stationxml.PolesZerosType, freq float64) float64 {

	var w complex128
	switch pz.PzTransferFunctionType {
	case stationxml.LaplaceRadiansSecondPzTransferFunction:
		w = complex(0.0, 2.0*math.Pi*freq)
	default:
		w = complex(0.0, freq)
	}

	h := complex(float64(1.0), float64(0.0))

	for _, zero := range pz.Zero {
		h *= (w - complex(zero.Real.Value, zero.Imaginary.Value))
	}

	for _, pole := range pz.Pole {
		h /= (w - complex(pole.Real.Value, pole.Imaginary.Value))
	}

	return cmplx.Abs(h)
}

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

func (r *Response) Poly() *stationxml.PolynomialType {
	for _, s := range r.stages {
		if s.Polynomial == nil {
			continue
		}
		return s.Polynomial
	}
	return nil
}

func (r *Response) Coeffs() []stationxml.Coefficient {
	poly := r.Poly()
	if poly == nil {
		return nil
	}

	var coeffs []stationxml.Coefficient

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

	return coeffs
}

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

func (r *Response) Derived(data []byte) (*stationxml.ResponseType, error) {

	var derived stationxml.ResponseType
	if err := xml.Unmarshal(data, &derived); err != nil {
		return nil, err
	}

	// must have at least an instrument sensitivity or polynomial
	if derived.InstrumentSensitivity == nil && derived.InstrumentPolynomial == nil {
		return nil, nil
	}

	var stages []stationxml.ResponseStageType
	for n, s := range derived.Stage {
		var stage stationxml.ResponseStageType
		if err := clone(&s, &stage); err != nil {
			return nil, err
		}

		if stage.PolesZeros != nil {
			stage.PolesZeros.Name = fmt.Sprintf("%sstage_%d", r.prefix, n+1)
		}
		if stage.Coefficients != nil {
			stage.Coefficients.Name = fmt.Sprintf("%sstage_%d", r.prefix, n+1)
		}
		if stage.Polynomial != nil {
			stage.Polynomial.Name = fmt.Sprintf("%sstage_%d", r.prefix, n+1)
		}
		stage.Number = stationxml.CounterType(n + 1)

		stages = append(stages, stage)
	}

	derived.Stage = stages

	return &derived, nil
}

func (r *Response) Normalise() error {

	var stages []stationxml.ResponseStageType
	for n, s := range append(r.sensor.Stage, r.datalogger.Stage...) {
		var stage stationxml.ResponseStageType
		if err := clone(&s, &stage); err != nil {
			return err
		}

		if stage.PolesZeros != nil {
			stage.PolesZeros.Name = fmt.Sprintf("%sstage_%d", r.prefix, n+1)
		}
		if stage.Coefficients != nil {
			stage.Coefficients.Name = fmt.Sprintf("%sstage_%d", r.prefix, n+1)
		}
		if stage.Polynomial != nil {
			stage.Polynomial.Name = fmt.Sprintf("%sstage_%d", r.prefix, n+1)
		}

		stage.Number = stationxml.CounterType(n + 1)

		if stage.StageGain != nil {
			scale := stage.StageGain.Value
			if stage.PolesZeros != nil {
				g, z := gain(stage.PolesZeros, r.freq), gain(stage.PolesZeros, stage.PolesZeros.NormalizationFrequency.Value)
				stage.PolesZeros.NormalizationFactor = 1.0 / g
				stage.PolesZeros.NormalizationFrequency = stationxml.FrequencyType{stationxml.FloatType{Value: r.freq}}
				scale /= (z / g)
			}
			stage.StageGain = &stationxml.GainType{
				Value:     scale,
				Frequency: r.freq,
			}
		}
		stages = append(stages, stage)
	}

	r.stages = stages

	return nil
}

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
