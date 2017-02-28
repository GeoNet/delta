package stationxml

// Complex numbers used as poles or zeros in channel response.
type PoleZero struct {
	Number uint32 `xml:"number,attr"`

	Real      FloatNoUnit
	Imaginary FloatNoUnit
}

func (pz PoleZero) IsValid() error {

	if err := Validate(pz.Real); err != nil {
		return err
	}
	if err := Validate(pz.Imaginary); err != nil {
		return err
	}

	return nil
}
