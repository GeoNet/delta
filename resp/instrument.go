package resp

import (
	"encoding/xml"
	"errors"
	"fmt"
	"strings"
)

var ErrInvalidResponse = errors.New("attempt to correct the wrong type of response, biases require a polynomial response")

type InstrumentResponseOpt func(*InstrumentResponse)

// InstrumentResponse is used for building an instrument response based on sensor and datalogger pairs. It makes no assumption about
// the StationXML version, ideally it should encompass all required elements. The conversion from a bas Response to a
// particular version is done via encoding interfaces.
type InstrumentResponse struct {
	Prefix        string
	Serial        string
	Frequency     float64
	ScaleFactor   float64
	ScaleBias     float64
	ScaleAbsolute float64
	GainFactor    float64
	GainBias      float64
	GainAbsolute  float64
	Telemetry     float64
	Preamp        float64

	sensor     *ResponseType
	datalogger *ResponseType

	stages []ResponseStageType
}

// Prefix sets the label used to prefix Response element names.
func Prefix(prefix string) InstrumentResponseOpt {
	return func(r *InstrumentResponse) {
		r.Prefix = prefix
	}
}

// Serial sets the label used to prefix Response equipment labels.
func Serial(serial string) InstrumentResponseOpt {
	return func(r *InstrumentResponse) {
		r.Serial = serial
	}
}

// Frequency is used to set the overall reference frequency for the Response.
func Frequency(frequency float64) InstrumentResponseOpt {
	return func(r *InstrumentResponse) {
		r.Frequency = frequency
	}
}

// Calibration is used to set a initial sensor reference gain, this overrides the default values.
func Calibration(factor, bias, absolute float64) InstrumentResponseOpt {
	return func(r *InstrumentResponse) {
		r.ScaleFactor = factor
		r.ScaleBias = bias
		r.ScaleAbsolute = absolute
	}
}

// Gain is used to adjusts the installed sensor gains, this is in addition to the default values.
func Gain(factor, bias, absolute float64) InstrumentResponseOpt {
	return func(r *InstrumentResponse) {
		r.GainFactor = factor
		r.GainBias = bias
		r.GainAbsolute = absolute
	}
}

// Telemetry is used to adjusts the sensor and datalogger connection gain, this is in addition to the default values.
func Telemetry(gain float64) InstrumentResponseOpt {
	return func(r *InstrumentResponse) {
		r.Telemetry = gain
	}
}

// Preamp is used to adjusts the datalogger gains, this is in addition to the default values.
func Preamp(preamp float64) InstrumentResponseOpt {
	return func(r *InstrumentResponse) {
		r.Preamp = preamp
	}
}

// NewResponse builds a Response with the given options.
func NewInstrumentResponse(opts ...InstrumentResponseOpt) *InstrumentResponse {
	r := InstrumentResponse{
		ScaleFactor: 1.0,
		GainFactor:  1.0,
	}
	r.Config(opts...)
	return &r
}

// Config can be used to set extra Response options.
func (r *InstrumentResponse) Config(opts ...InstrumentResponseOpt) {
	for _, opt := range opts {
		opt(r)
	}
}

// SetPrefix sets the label used to prefix Response element names.
func (r *InstrumentResponse) SetPrefix(prefix string) {
	Prefix(prefix)(r)
}

// SetSerial sets the label used to prefix Response equipment labels.
func (r *InstrumentResponse) SetSerial(serial string) {
	Serial(serial)(r)
}

// SetFrequency is used to set the overall reference frequency for the Response.
func (r *InstrumentResponse) SetFrequency(frequency float64) {
	Frequency(frequency)(r)
}

// SetCalibration is used to set a initial sensor reference gain, this overrides the default values.
func (r *InstrumentResponse) SetCalibration(scale, bias, absolute float64) {
	Calibration(scale, bias, absolute)(r)
}

// SetGain is used to adjusts the installed sensor gains, this is in addition to the default values.
func (r *InstrumentResponse) SetGain(scale, bias, absolute float64) {
	Gain(scale, bias, absolute)(r)
}

// SetPreamp is used to adjusts the datalogger gains, this is in addition to the default values.
func (r *InstrumentResponse) SetPreamp(preamp float64) {
	Preamp(preamp)(r)
}

// SetTelemetry is used to adjusts the datalogger gains, this is in addition to the default values.
func (r *InstrumentResponse) SetTelemetry(gain float64) {
	Telemetry(gain)(r)
}

// Polynomial finds the PolynomialType in the Response if one is present.
func (r *InstrumentResponse) Polynomial() *PolynomialType {
	for _, s := range r.stages {
		if s.Polynomial == nil {
			continue
		}
		return s.Polynomial
	}
	return nil
}

// SetSensor takes an XML encoded ResponseType that represents a Sensor and adds it to the Response.
func (r *InstrumentResponse) SetSensor(data []byte) error {

	var sensor ResponseType
	if err := xml.Unmarshal(data, &sensor); err != nil {
		return err
	}

	switch {
	case sensor.InstrumentSensitivity != nil:
		// check biases, as these imply a polynomial response
		if r.ScaleBias != 0.0 || r.GainBias != 0.0 || r.GainAbsolute != 0.0 {
			return ErrInvalidResponse
		}
		// a simple scaling of the overall sensitivity
		sensor.InstrumentSensitivity.Value *= (r.ScaleFactor * r.GainFactor)
		// look for the first stage with a gain
		for i := range sensor.Stages {
			stage := sensor.Stages[i]
			if stage.StageGain == nil {
				continue
			}
			// update the first one found, and ignore the rest
			stage.StageGain.Value *= (r.ScaleFactor * r.GainFactor)
			break
		}
	case sensor.InstrumentPolynomial != nil:

		// First adjust for any calibrations, these are simply replacing the first two coefficients
		if r.ScaleFactor != 1.0 || r.ScaleAbsolute != 0.0 {
			for i := range sensor.Stages {
				stage := sensor.Stages[i]

				// there can only be one
				if stage.Polynomial == nil {
					continue
				}

				// ignore changes in units, or changes in bounds
				switch c := len(stage.Polynomial.Coefficients); {
				case c > 1:
					stage.Polynomial.Coefficients[1] = PolynomialCoefficient{
						Number: stage.Polynomial.Coefficients[1].Number,
						Value:  r.ScaleFactor,
					}
					stage.Polynomial.Coefficients[0] = PolynomialCoefficient{
						Number: stage.Polynomial.Coefficients[0].Number,
						Value:  r.ScaleAbsolute,
					}
				case c > 0:
					stage.Polynomial.Coefficients[0] = PolynomialCoefficient{
						Number: stage.Polynomial.Coefficients[0].Number,
						Value:  r.ScaleAbsolute,
					}
				}

				// update the overall instrument polynomial
				sensor.InstrumentPolynomial.Coefficients = append([]PolynomialCoefficient{}, stage.Polynomial.Coefficients...)
			}
		}

		// Second adjust for any gains, these will update the first two coefficents
		if r.GainFactor != 1.0 || r.GainBias != 0.0 || r.GainAbsolute != 0.0 {

			for i := range sensor.Stages {
				stage := sensor.Stages[i]

				// there can only be one
				if stage.Polynomial == nil {
					continue
				}

				// ignore changes in units, or changes in bounds
				switch c := len(stage.Polynomial.Coefficients); {
				case c > 1:
					stage.Polynomial.Coefficients[0] = PolynomialCoefficient{
						Number: stage.Polynomial.Coefficients[0].Number,
						Value: stage.Polynomial.Coefficients[0].Value*r.GainFactor +
							stage.Polynomial.Coefficients[1].Value*r.GainBias*r.ScaleBias + r.GainAbsolute,
					}
					stage.Polynomial.Coefficients[1] = PolynomialCoefficient{
						Number: stage.Polynomial.Coefficients[1].Number,
						Value:  stage.Polynomial.Coefficients[1].Value * r.GainFactor,
					}
				case c > 0:
					stage.Polynomial.Coefficients[0] = PolynomialCoefficient{
						Number: stage.Polynomial.Coefficients[0].Number,
						Value:  stage.Polynomial.Coefficients[0].Value*r.GainFactor + r.GainAbsolute,
					}
				}

				// update the overall instrument polynomial
				sensor.InstrumentPolynomial.Coefficients = append([]PolynomialCoefficient{}, stage.Polynomial.Coefficients...)
			}
		}
	default:
		return nil
	}

	r.sensor = &sensor

	return nil
}

// SetDatalogger takes an XML encoded ResponseType that represents a Datalogger and adds it to the Response.
func (r *InstrumentResponse) SetDatalogger(data []byte) error {

	var datalogger ResponseType
	if err := xml.Unmarshal(data, &datalogger); err != nil {
		return err
	}

	if datalogger.InstrumentSensitivity == nil {
		return nil
	}

	// a telemetry gain has been given, prepend an appropriate stage
	if r.Telemetry != 0.0 {
		datalogger.Stages = append([]ResponseStageType{{
			//TODO: technically the poles and zeros are not required, but kept to allow acceptance checks
			PolesZeros: &PolesZerosType{
				ResourceId:             "PolesZeros#Telemetry",
				InputUnits:             datalogger.InstrumentSensitivity.InputUnits,
				OutputUnits:            datalogger.InstrumentSensitivity.InputUnits,
				PzTransferFunctionType: LaplaceRadiansSecondPzTransferFunction,
				NormalizationFactor:    1.0,
				NormalizationFrequency: r.Frequency,
			},
			StageGain: &StageGain{
				Value:     r.Telemetry,
				Frequency: r.Frequency,
			},
		}}, datalogger.Stages...)
	}

	// a preamp has been given, prepend an appropriate stage
	if r.Preamp != 0.0 {
		datalogger.Stages = append([]ResponseStageType{{
			//TODO: technically the poles and zeros are not required, but kept to allow acceptance checks
			PolesZeros: &PolesZerosType{
				ResourceId:             "PolesZeros#Preamp",
				InputUnits:             datalogger.InstrumentSensitivity.InputUnits,
				OutputUnits:            datalogger.InstrumentSensitivity.InputUnits,
				PzTransferFunctionType: LaplaceRadiansSecondPzTransferFunction,
				NormalizationFactor:    1.0,
				NormalizationFrequency: r.Frequency,
			},
			StageGain: &StageGain{
				Value:     r.Preamp,
				Frequency: r.Frequency,
			},
		}}, datalogger.Stages...)
	}

	r.datalogger = &datalogger

	return nil
}

// Derived returns a ResponseType when there is only a single set of derived Response stages.
func (r *InstrumentResponse) Derived(data []byte) (*ResponseType, error) {

	var derived ResponseType
	if err := xml.Unmarshal(data, &derived); err != nil {
		return nil, err
	}

	// must have at least an instrument sensitivity or polynomial
	if derived.InstrumentSensitivity == nil && derived.InstrumentPolynomial == nil {
		return nil, nil
	}

	var stages []ResponseStageType
	for n, s := range derived.Stages {
		stage, err := s.Clone()
		if err != nil {
			return nil, err
		}

		if stage.PolesZeros != nil {
			stage.PolesZeros.Name = fmt.Sprintf("%sstage_%d", r.Prefix, n+1)
		}
		if stage.Coefficients != nil {
			stage.Coefficients.Name = fmt.Sprintf("%sstage_%d", r.Prefix, n+1)
		}
		if stage.Polynomial != nil {
			stage.Polynomial.Name = fmt.Sprintf("%sstage_%d", r.Prefix, n+1)
		}
		stage.Number = n + 1

		stages = append(stages, stage)
	}

	derived.Stages = stages

	return &derived, nil
}

// Coeffs returns a slice of PolynomialCoeffiencent values present in the Response.
func (r *InstrumentResponse) Coeffs() []PolynomialCoefficient {
	var coeffs []PolynomialCoefficient

	if p := r.Polynomial(); p != nil && len(p.Coefficients) > 0 {
		coeffs = append(coeffs, PolynomialCoefficient{
			Number: len(coeffs) + 1,
			Value:  p.Coefficients[0].Value,
		})
	}

	if p := r.Polynomial(); p != nil && len(p.Coefficients) > 1 {
		if scale := r.Scale(); scale != 0.0 {
			coeffs = append(coeffs, PolynomialCoefficient{
				Number: len(coeffs) + 1,
				Value:  p.Coefficients[1].Value / scale,
			})
		}
	}

	return coeffs
}

// Scale calculates the overall response scale factor.
func (r *InstrumentResponse) Scale() float64 {
	scale := 1.0
	for _, s := range r.stages {
		if s.StageGain == nil || s.StageGain.Value == 0.0 {
			continue
		}
		scale *= s.StageGain.Value
	}
	return scale
}

// Normalise adjusts the labels and stage gains of a Response.
func (r *InstrumentResponse) Normalise() error {
	var stages []ResponseStageType
	for n, s := range append(r.sensor.Stages, r.datalogger.Stages...) {

		stage, err := s.Clone()
		if err != nil {
			return err
		}

		if stage.PolesZeros != nil {
			stage.PolesZeros.Name = fmt.Sprintf("%sstage_%d", r.Prefix, n+1)
			for i := range stage.PolesZeros.Zeros {
				stage.PolesZeros.Zeros[i].Number = i + 1
			}
			for i := range stage.PolesZeros.Poles {
				stage.PolesZeros.Poles[i].Number = i + 1
			}
		}
		if stage.Coefficients != nil {
			stage.Coefficients.Name = fmt.Sprintf("%sstage_%d", r.Prefix, n+1)
			for i := range stage.Coefficients.Numerators {
				stage.Coefficients.Numerators[i].Number = i + 1
			}
			for i := range stage.Coefficients.Denominators {
				stage.Coefficients.Denominators[i].Number = i + 1
			}
		}
		if stage.Polynomial != nil {
			stage.Polynomial.Name = fmt.Sprintf("%sstage_%d", r.Prefix, n+1)
			for i := range stage.Polynomial.Coefficients {
				stage.Polynomial.Coefficients[i].Number = i + 1
			}
		}

		stage.Number = n + 1

		if stage.StageGain != nil {
			scale := stage.StageGain.Value
			if stage.PolesZeros != nil {
				g, z := stage.PolesZeros.Gain(r.Frequency), stage.PolesZeros.Gain(stage.PolesZeros.NormalizationFrequency)
				stage.PolesZeros.NormalizationFactor = 1.0 / g
				stage.PolesZeros.NormalizationFrequency = r.Frequency
				scale /= (z / g)
			}
			stage.StageGain = &StageGain{
				Value:     scale,
				Frequency: r.Frequency,
			}
		}
		stages = append(stages, stage)
	}

	r.stages = stages

	return nil
}

// ResponseType builds a combined ResponseType from a Response.
func (r *InstrumentResponse) ResponseType() (*ResponseType, error) {

	if err := r.Normalise(); err != nil {
		return nil, err
	}

	resp := ResponseType{
		InstrumentSensitivity: func() *InstrumentSensitivity {
			if r.sensor.InstrumentSensitivity != nil {
				return &InstrumentSensitivity{
					InputUnits:  r.sensor.InstrumentSensitivity.InputUnits,
					OutputUnits: r.datalogger.InstrumentSensitivity.OutputUnits,
					Frequency:   r.Frequency,
					Value:       r.Scale(),
				}
			}
			return nil
		}(),
		InstrumentPolynomial: func() *InstrumentPolynomial {
			if poly := r.Polynomial(); poly != nil {
				return &InstrumentPolynomial{
					ResourceId:              "Instrument" + poly.ResourceId + ":" + r.Serial,
					Name:                    strings.TrimRight(r.Prefix, "."),
					InputUnits:              poly.InputUnits,
					OutputUnits:             r.datalogger.InstrumentSensitivity.OutputUnits,
					ApproximationType:       poly.ApproximationType,
					FrequencyLowerBound:     poly.FrequencyLowerBound,
					FrequencyUpperBound:     poly.FrequencyUpperBound,
					ApproximationLowerBound: poly.ApproximationLowerBound,
					ApproximationUpperBound: poly.ApproximationUpperBound,
					MaximumError:            poly.MaximumError,
					Coefficients:            r.Coeffs(),
				}
			}
			return nil
		}(),
		Stages: r.stages,

		// used for polynomial calculations
		frequency: r.Frequency,
	}

	return &resp, nil
}

// Marshal generates an XML encoded version of the Response as a ResponseType.
func (r *InstrumentResponse) Marshal() ([]byte, error) {
	resp, err := r.ResponseType()
	if err != nil {
		return nil, err
	}

	data, err := xml.Marshal(resp)
	if err != nil {
		return nil, err
	}

	return data, nil
}
