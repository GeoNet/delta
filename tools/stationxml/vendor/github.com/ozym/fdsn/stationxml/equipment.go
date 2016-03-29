package stationxml

type Equipment struct {
	// This field contains a string that should serve as a unique resource identifier.
	// This identifier can be interpreted differently depending on the datacenter/software
	// that generated the document. Also, we recommend to use something like
	// GENERATOR:Meaningful ID. As a common behaviour equipment with the same ID should
	// contains the same information/be derived from the same base instruments.
	ResourceId string `xml:"resourceId,attr,omitempty" json:",omitempty"`

	Type             string     `xml:",omitempty" json:",omitempty"`
	Description      string     `xml:",omitempty" json:",omitempty"`
	Manufacturer     string     `xml:",omitempty" json:",omitempty"`
	Vendor           string     `xml:",omitempty" json:",omitempty"`
	Model            string     `xml:",omitempty" json:",omitempty"`
	SerialNumber     string     `xml:",omitempty" json:",omitempty"`
	InstallationDate *DateTime  `xml:",omitempty" json:",omitempty"`
	RemovalDate      *DateTime  `xml:",omitempty" json:",omitempty"`
	CalibrationDates []DateTime `xml:"CalibrationDate,omitempty" json:",omitempty"`
}

func (e Equipment) IsValid() error {

	if e.InstallationDate != nil {
		if err := Validate(e.InstallationDate); err != nil {
			return err
		}
	}

	if e.RemovalDate != nil {
		if err := Validate(e.RemovalDate); err != nil {
			return err
		}
	}

	for _, c := range e.CalibrationDates {
		if err := Validate(c); err != nil {
			return err
		}
	}

	return nil
}
