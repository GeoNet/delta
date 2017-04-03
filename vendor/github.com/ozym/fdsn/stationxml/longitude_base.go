package stationxml

import (
	"fmt"
)

type LongitudeBase struct {
	Float
}

func (l LongitudeBase) IsValid() error {

	if l.Unit != "" && l.Unit != "DEGREES" {
		return fmt.Errorf("invalid longitude unit: %s", l.Unit)
	}
	if l.Value < -180 || l.Value > 180 {
		return fmt.Errorf("longitude outside range: %g", l.Value)
	}

	return nil
}
