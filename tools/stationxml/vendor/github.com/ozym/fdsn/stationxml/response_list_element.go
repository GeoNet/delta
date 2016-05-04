package stationxml

type ResponseListElement struct {
	Frequency Frequency
	Amplitude Float
	Phase     Angle
}

func (r ResponseListElement) IsValid() error {

	if err := Validate(r.Frequency); err != nil {
		return err
	}
	if err := Validate(r.Amplitude); err != nil {
		return err
	}
	if err := Validate(r.Phase); err != nil {
		return err
	}

	return nil
}
