package stationxml

// Response: list of frequency, amplitude and phase values. Corresponds to SEED blockette 55.
type ResponseList struct {
	BaseFilter

	ResponseListElements []ResponseListElement `xml:"ResponseListElement,omitempty" json:",omitempty"`
}

func (r ResponseList) IsValid() error {

	if err := Validate(r.BaseFilter); err != nil {
		return err
	}

	for _, e := range r.ResponseListElements {
		if err := Validate(e); err != nil {
			return err
		}
	}

	return nil
}
