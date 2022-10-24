package build

import (
	"encoding/xml"

	"github.com/GeoNet/delta/internal/stationxml/v1.2"
)

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
