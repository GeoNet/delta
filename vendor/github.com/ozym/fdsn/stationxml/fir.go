package stationxml

// Response: FIR filter. Corresponds to SEED blockette 61. FIR filters are
// also commonly documented using the CoefficientsType element.
type FIR struct {
	BaseFilter

	Symmetry              Symmetry
	NumeratorCoefficients []NumeratorCoefficient `xml:"NumeratorCoefficient,omitempty" json:",omitempty"`
}

func (f FIR) IsValid() error {

	if err := Validate(f.BaseFilter); err != nil {
		return err
	}

	if err := Validate(f.Symmetry); err != nil {
		return err
	}

	return nil
}
