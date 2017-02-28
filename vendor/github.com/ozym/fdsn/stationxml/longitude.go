package stationxml

// Type for latitude coordinate. min: -180 max: 180
type Longitude struct {
	LongitudeBase
	Datum string `xml:"datum,attr,omitempty" json:",omitempty"` // WGS84
}

func (l Longitude) IsValid() error {

	if err := Validate(l.LongitudeBase); err != nil {
		return err
	}

	return nil
}
