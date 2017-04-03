package stationxml

import (
	"fmt"
)

type Angle struct {
	Float
}

func (a Angle) IsValid() error {

	if err := Validate(a.Float); err != nil {
		return err
	}

	if a.Unit != "" && a.Unit != "DEGREES" {
		return fmt.Errorf("invalid unit: %s", a.Unit)
	}
	if a.Value < -360 || a.Value > 360 {
		return fmt.Errorf("angle outside range: %g", a.Value)
	}

	return nil
}
