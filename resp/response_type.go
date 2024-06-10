package resp

import (
	"bytes"
	"encoding/gob"
	"encoding/xml"
	"math"
	"math/cmplx"
)

const LaplaceRadiansSecondPzTransferFunction = "LAPLACE (RADIANS/SECOND)"

type Float struct {
	Value float64 `xml:",chardata"`
}

type Units struct {
	Name        string `xml:"Name"`
	Description string `xml:"Description,omitempty"`
}

type ApproximationBound struct {
	Value float64 `xml:",chardata"`
}

type NumeratorCoefficient struct {
	I     int     `xml:"i,attr"`
	Value float64 `xml:",chardata"`
}

type PolynomialCoefficient struct {
	Number int     `xml:"number,attr"`
	Value  float64 `xml:",chardata"`
}

type CoefficientNumerator struct {
	Number int     `xml:"number,attr"`
	Value  float64 `xml:",chardata"`
}

type CoefficientDenominator struct {
	Number int     `xml:"number,attr"`
	Value  float64 `xml:",chardata"`
}

type PoleZero struct {
	Number int `xml:"number,attr"`

	Real      Float `xml:"Real"`
	Imaginary Float `xml:"Imaginary"`
}

type CoefficientsType struct {
	ResourceId  string `xml:"resourceId,attr,omitempty"`
	Name        string `xml:"name,attr,omitempty"`
	Description string `xml:"description,attr,omitempty"`

	InputUnits             Units  `xml:"InputUnits"`
	OutputUnits            Units  `xml:"OutputUnits"`
	CfTransferFunctionType string `xml:"CfTransferFunctionType"`

	Numerators   []CoefficientNumerator   `xml:"Numerator,omitempty"`
	Denominators []CoefficientDenominator `xml:"Denominator,omitempty"`
}

type DecimationType struct {
	InputSampleRate float64 `xml:"InputSampleRate"`
	Factor          int     `xml:"Factor"`
	Offset          int     `xml:"Offset"`
	Delay           float64 `xml:"Delay"`
	Correction      float64 `xml:"Correction"`
}

type FirType struct {
	ResourceId  string `xml:"resourceId,attr,omitempty"`
	Name        string `xml:"name,attr,omitempty"`
	Description string `xml:"description,attr,omitempty"`

	InputUnits  Units  `xml:"InputUnits"`
	OutputUnits Units  `xml:"OutputUnits"`
	Symmetry    string `xml:"Symmetry"`

	NumeratorCoefficients []NumeratorCoefficient `xml:"NumeratorCoefficient"`
}

type PolesZerosType struct {
	ResourceId  string `xml:"resourceId,attr,omitempty"`
	Name        string `xml:"name,attr,omitempty"`
	Description string `xml:"description,attr,omitempty"`

	InputUnits  Units `xml:"InputUnits"`
	OutputUnits Units `xml:"OutputUnits"`

	PzTransferFunctionType string  `xml:"PzTransferFunctionType"`
	NormalizationFactor    float64 `xml:"NormalizationFactor"`
	NormalizationFrequency float64 `xml:"NormalizationFrequency"`

	Zeros []PoleZero `xml:"Zero"`
	Poles []PoleZero `xml:"Pole"`
}

// Gain ccalculates the poles and zeros response gain at a given frequency
func (pz PolesZerosType) Gain(freq float64) float64 {

	var w complex128
	switch pz.PzTransferFunctionType {
	case LaplaceRadiansSecondPzTransferFunction:
		w = complex(0.0, 2.0*math.Pi*freq)
	default:
		w = complex(0.0, freq)
	}

	h := complex(float64(1.0), float64(0.0))

	for _, zero := range pz.Zeros {
		h *= (w - complex(zero.Real.Value, zero.Imaginary.Value))
	}

	for _, pole := range pz.Poles {
		h /= (w - complex(pole.Real.Value, pole.Imaginary.Value))
	}

	return cmplx.Abs(h)
}

type PolynomialType struct {
	ResourceId  string `xml:"resourceId,attr,omitempty"`
	Name        string `xml:"name,attr,omitempty"`
	Description string `xml:"description,attr,omitempty"`

	InputUnits  Units `xml:"InputUnits"`
	OutputUnits Units `xml:"OutputUnits"`

	ApproximationType       string             `xml:"ApproximationType"`
	FrequencyLowerBound     float64            `xml:"FrequencyLowerBound"`
	FrequencyUpperBound     float64            `xml:"FrequencyUpperBound"`
	ApproximationLowerBound ApproximationBound `xml:"ApproximationLowerBound"`
	ApproximationUpperBound ApproximationBound `xml:"ApproximationUpperBound"`
	MaximumError            float64            `xml:"MaximumError"`

	Coefficients []PolynomialCoefficient `xml:"Coefficient,omitempty"`
}

func (p PolynomialType) Value(input float64) float64 {
	var value float64
	for n, c := range p.Coefficients {
		value += c.Value * math.Pow(input, float64(n))
	}
	return value
}

type StageGain struct {
	Value     float64 `xml:"Value"`
	Frequency float64 `xml:"Frequency"`
}

type ResponseStageType struct {
	Number     int    `xml:"number,attr"`
	ResourceId string `xml:"resourceId,attr,omitempty"`

	Coefficients *CoefficientsType `xml:"Coefficients,omitempty"`
	Decimation   *DecimationType   `xml:"Decimation,omitempty"`
	FIR          *FirType          `xml:"FIR,omitempty"`
	PolesZeros   *PolesZerosType   `xml:"PolesZeros,omitempty"`
	Polynomial   *PolynomialType   `xml:"Polynomial,omitempty"`

	StageGain *StageGain `xml:"StageGain,omitempty"`
}

// clone two responses to avoid shared backing arrays
func (r *ResponseStageType) Clone() (ResponseStageType, error) {

	var buff bytes.Buffer

	if err := gob.NewEncoder(&buff).Encode(r); err != nil {
		return ResponseStageType{}, err
	}

	var c ResponseStageType
	if err := gob.NewDecoder(&buff).Decode(&c); err != nil {
		return ResponseStageType{}, err
	}

	return c, nil
}

type InstrumentSensitivity struct {
	Value       float64 `xml:"Value"`
	Frequency   float64 `xml:"Frequency"`
	InputUnits  Units   `xml:"InputUnits"`
	OutputUnits Units   `xml:"OutputUnits"`
}

type InstrumentPolynomial struct {
	ResourceId  string `xml:"resourceId,attr,omitempty"`
	Name        string `xml:"name,attr"`
	Description string `xml:"description,attr,omitempty"`

	InputUnits  Units `xml:"InputUnits"`
	OutputUnits Units `xml:"OutputUnits"`

	ApproximationType       string             `xml:"ApproximationType"`
	FrequencyLowerBound     float64            `xml:"FrequencyLowerBound"`
	FrequencyUpperBound     float64            `xml:"FrequencyUpperBound"`
	ApproximationLowerBound ApproximationBound `xml:"ApproximationLowerBound"`
	ApproximationUpperBound ApproximationBound `xml:"ApproximationUpperBound"`
	MaximumError            float64            `xml:"MaximumError"`

	Coefficients []PolynomialCoefficient `xml:"Coefficient,omitempty"`
}

// ResponseType is a struct that mimics the StationXML ResponseType element, but is not constrained to a particular version.
type ResponseType struct {
	XMLName    xml.Name `xml:"Response"`
	ResourceId string   `xml:"resourceId,attr,omitempty"`

	InstrumentSensitivity *InstrumentSensitivity `xml:"InstrumentSensitivity,omitempty"`
	InstrumentPolynomial  *InstrumentPolynomial  `xml:"InstrumentPolynomial,omitempty"`

	Stages []ResponseStageType `xml:"Stage,omitempty"`

	// used for instrument polynomial calculations
	frequency float64
}

func NewResponseType(data []byte) (*ResponseType, error) {
	var s ResponseType
	if err := s.Unmarshal(data); err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *ResponseType) Frequency() float64 {
	return r.frequency
}

func (r *ResponseType) SetFrequency(freq float64) {
	r.frequency = freq
}

func (r *ResponseType) Scale() float64 {
	scale := 1.0
	for _, s := range r.Stages {
		if s.StageGain == nil || s.StageGain.Value == 0.0 {
			continue
		}
		scale *= s.StageGain.Value
	}
	return scale
}

func (r *ResponseType) PolynomialType() *PolynomialType {
	for _, s := range r.Stages {
		if s.Polynomial == nil {
			continue
		}
		return s.Polynomial
	}
	return nil
}

func (r *ResponseType) PolynomialCoefficients() []PolynomialCoefficient {
	var coeffs []PolynomialCoefficient

	if p := r.PolynomialType(); p != nil && len(p.Coefficients) > 0 {
		coeffs = append(coeffs, p.Coefficients[0])
	}

	if p := r.PolynomialType(); p != nil && len(p.Coefficients) > 1 {
		if scale := r.Scale(); scale != 0.0 {
			coeffs = append(coeffs, PolynomialCoefficient{
				Number: 1,
				Value:  p.Coefficients[1].Value / scale,
			})
		}
	}

	return coeffs
}

func (r *ResponseType) Unmarshal(data []byte) error {
	return xml.Unmarshal(data, r)
}

func (r ResponseType) Marshal() ([]byte, error) {
	body, err := xml.MarshalIndent(r, "", "  ")
	if err != nil {
		return nil, err
	}
	head := []byte(xml.Header)
	return append(head, append(body, '\n')...), nil
}
