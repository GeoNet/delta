package stationxml

type FloatNoUnit struct {
	UncertaintyDouble

	Value float64 `xml:",chardata"`
}

func (f FloatNoUnit) IsValid() error {

	if err := Validate(f.UncertaintyDouble); err != nil {
		return err
	}

	return nil
}
