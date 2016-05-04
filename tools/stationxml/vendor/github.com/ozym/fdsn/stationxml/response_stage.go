package stationxml

// This complex type represents channel response and covers SEED blockettes 53 to 56.
type ResponseStage struct {
	Number Counter `xml:"number,attr"`

	// A choice of response types. There should be one response per stage.
	PolesZeros   *PolesZeros   `xml:",omitempty" json:",omitempty"`
	Coefficients *Coefficients `xml:",omitempty" json:",omitempty"`
	ResponseList *ResponseList `xml:,omitempty" json:",omitempty"`
	FIR          *FIR          `xml:",omitempty" json:",omitempty"`
	Polynomial   *Polynomial   `xml:",omitempty" json:",omitempty"`

	Decimation *Decimation `xml:",omitempty" json:",omitempty"`
	StageGain  Gain
}

func (r ResponseStage) IsValid() error {

	if r.PolesZeros != nil {
		if err := Validate(r.PolesZeros); err != nil {
			return err
		}
	}
	if r.Coefficients != nil {
		if err := Validate(r.Coefficients); err != nil {
			return err
		}
	}

	if r.ResponseList != nil {
		if err := Validate(r.ResponseList); err != nil {
			return err
		}
	}

	if r.FIR != nil {
		if err := Validate(r.FIR); err != nil {
			return err
		}
	}
	if r.Polynomial != nil {
		if err := Validate(r.Polynomial); err != nil {
			return err
		}
	}

	if r.Decimation != nil {
		if err := Validate(r.Decimation); err != nil {
			return err
		}
	}

	return nil
}
