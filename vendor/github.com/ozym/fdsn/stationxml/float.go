package stationxml

type Float struct {
	UncertaintyDouble
	Unit string `xml:"unit,attr,omitempty" json:",omitempty"`

	Value float64 `xml:",chardata"`
}

func (f Float) IsValid() error {

	if err := Validate(f.UncertaintyDouble); err != nil {
		return err
	}

	return nil
}
