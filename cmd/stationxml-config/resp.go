package main

import (
	//	"log"

	//	"encoding/xml"
	//"fmt"
	//	"math"
	//"math/cmplx"
	//	"strings"

	"github.com/GeoNet/delta/internal/build/v1.2"
	"github.com/GeoNet/delta/internal/stationxml/v1.2"
)

/*
func Derived(prefix string, data []byte) (*stationxml.ResponseType, error) {

	if len(data) == 0 {
		return nil, nil
	}

	var derived stationxml.ResponseType

	if err := xml.Unmarshal(data, &derived); err != nil {
		return nil, err
	}

	if derived.InstrumentSensitivity == nil && derived.InstrumentPolynomial == nil {
		return nil, nil
	}

	for n, x := range derived.Stage {
		if x.PolesZeros != nil {
			derived.Stage[n].PolesZeros.Name = fmt.Sprintf("%sstage_%d", prefix, n+1)
		}
		if x.Coefficients != nil {
			derived.Stage[n].Coefficients.Name = fmt.Sprintf("%sstage_%d", prefix, n+1)
		}
		if x.Polynomial != nil {
			derived.Stage[n].Polynomial.Name = fmt.Sprintf("%sstage_%d", prefix, n+1)
		}
	}

	return &derived, nil
}
*/

func Resp(resp *build.Response, prefix, id string, freq, gain, bias, preamp float64, sensor, datalogger []byte) (*stationxml.ResponseType, error) {

	//resp := build.NewResponse(prefix, id, freq)

	//var stages []stationxml.ResponseStageType

	if len(sensor) == 0 {
		return nil, nil
	}

	/*
		if err := resp.Sensor(gain, bias, sensor); err != nil {
			return nil, err
		}
	*/

	/*
		s, err := build.Sensor(gain, bias, freq, sensor)
		if err != nil || s == nil {
			return nil, err
		}
		stages = append(stages, s.Stage...)
	*/

	if len(datalogger) == 0 {
		return nil, nil
	}
	/*

		if err := resp.Datalogger(preamp, datalogger); err != nil {
			return nil, err
		}
	*/

	/*
		d, err := build.Datalogger(preamp, freq, datalogger)
		if err != nil || d == nil {
			return nil, err
		}
		stages = append(stages, d.Stage...)
	*/

	/*
		res, err := build.Stages(prefix, freq, resp.Stages()...)
		if err != nil {
			return nil, err
		}
	*/

	/*
		if err := resp.Normalise(); err != nil {
			return nil, err
		}
	*/

	//scale := build.InstrumentScale(freq, res...)
	//poly := build.InstrumentPoly(res...)
	//coeffs := build.InstrumentCoeffs(freq, res...)

	//TODO:
	//stages = res

	/*
		r := &stationxml.ResponseType{
			InstrumentSensitivity: func() *stationxml.SensitivityType {
				if s.InstrumentSensitivity != nil {
					return &stationxml.SensitivityType{
						InputUnits:  s.InstrumentSensitivity.InputUnits,
						OutputUnits: d.InstrumentSensitivity.OutputUnits,
						GainType: stationxml.GainType{
							Frequency: freq,
							Value:     scale,
						},
					}
				}
				return nil
			}(),
			InstrumentPolynomial: func() *stationxml.PolynomialType {
				if poly != nil {
					return &stationxml.PolynomialType{
						BaseFilterType: stationxml.BaseFilterType{
							ResourceId:  "Instrument" + poly.ResourceId + ":" + id,
							Name:        strings.TrimRight(prefix, "."),
							InputUnits:  poly.InputUnits,
							OutputUnits: d.InstrumentSensitivity.OutputUnits,
						},
						ApproximationType:       poly.ApproximationType,
						FrequencyLowerBound:     poly.FrequencyLowerBound,
						FrequencyUpperBound:     poly.FrequencyUpperBound,
						ApproximationLowerBound: poly.ApproximationLowerBound,
						ApproximationUpperBound: poly.ApproximationUpperBound,
						MaximumError:            poly.MaximumError,
						Coefficient:             coeffs,
					}
				}
				return nil
			}(),
			Stage: res,
		}

		if false {
			return r, nil
		}
	*/
	return resp.ResponseType(), nil
}

/*
func Gain(pz *stationxml.PolesZerosType, freq float64) float64 {

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
*/
