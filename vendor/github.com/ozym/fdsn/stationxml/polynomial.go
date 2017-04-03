package stationxml

// Response: expressed as a polynomial (allows non-linear sensors to be described).
// Corresponds to SEED blockette 62. Can be used to describe a stage of acquisition or a complete system.
type Polynomial struct {
	BaseFilter

	ApproximationType       ApproximationType
	FrequencyLowerBound     Frequency
	FrequencyUpperBound     Frequency
	ApproximationLowerBound string
	ApproximationUpperBound string
	MaximumError            float64

	Coefficients []Coefficient `xml:"Coefficient,omitempty" json:",omitempty"`
}

func (p Polynomial) IsValid() error {

	if err := Validate(p.BaseFilter); err != nil {
		return err
	}

	if err := Validate(p.ApproximationType); err != nil {
		return err
	}

	if err := Validate(p.FrequencyLowerBound); err != nil {
		return err
	}
	if err := Validate(p.FrequencyUpperBound); err != nil {
		return err
	}

	return nil
}
