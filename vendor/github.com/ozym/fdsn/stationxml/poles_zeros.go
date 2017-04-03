package stationxml

type PolesZeros struct {
	BaseFilter

	PzTransferFunctionType PzTransferFunctionType
	NormalizationFactor    float64
	NormalizationFrequency Frequency
	Zeros                  []PoleZero `xml:"Zero,omitempty" json:",omitempty"`
	Poles                  []PoleZero `xml:"Pole,omitempty" json:",omitempty"`
}

func (pz PolesZeros) IsValid() error {

	if err := Validate(pz.BaseFilter); err != nil {
		return err
	}

	if err := Validate(pz.PzTransferFunctionType); err != nil {
		return err
	}
	if err := Validate(pz.NormalizationFrequency); err != nil {
		return err
	}

	for _, z := range pz.Zeros {
		if err := Validate(z); err != nil {
			return err
		}
	}

	for _, p := range pz.Poles {
		if err := Validate(p); err != nil {
			return err
		}
	}

	return nil
}
