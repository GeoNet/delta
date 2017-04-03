package stationxml

// Response: coefficients for FIR filter.
// Laplace transforms or IIR filters can be expressed using type as well but the
// PolesAndZerosType should be used instead. Corresponds to SEED blockette 54.
type Coefficients struct {
	BaseFilter

	CfTransferFunctionType CfTransferFunctionType

	Numerators   []Float `xml:"Numerator,omitempty" json:",omitempty"`
	Denominators []Float `xml:"Denominator,omitempty" json:",omitempty"`
}

func (c Coefficients) IsValid() error {

	if err := Validate(c.BaseFilter); err != nil {
		return err
	}

	if err := Validate(c.CfTransferFunctionType); err != nil {
		return err
	}

	for _, n := range c.Numerators {
		if err := Validate(n); err != nil {
			return err
		}
	}

	for _, d := range c.Denominators {
		if err := Validate(d); err != nil {
			return err
		}
	}

	return nil
}
