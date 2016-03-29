package stationxml

import (
	"fmt"
)

// Representation of floating-point numbers used as measurements. min: 0, max: 360
type Azimuth struct {
	Float
}

func (a Azimuth) IsValid() error {

	if err := Validate(a.Float); err != nil {
		return err
	}

	if a.Unit != "" && a.Unit != "DEGREES" {
		return fmt.Errorf("azimuth invalid unit: %s", a.Unit)
	}
	if a.Value < 0 || a.Value > 360 {
		return fmt.Errorf("azimuth outside range: %g", a.Value)
	}

	return nil
}
