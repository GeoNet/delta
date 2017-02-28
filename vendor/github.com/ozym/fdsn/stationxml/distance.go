package stationxml

// Extension of Float for distances, elevations, and depths.
type Distance struct {
	Float
}

func (d Distance) IsValid() error {

	if err := Validate(d.Float); err != nil {
		return err
	}

	return nil
}
