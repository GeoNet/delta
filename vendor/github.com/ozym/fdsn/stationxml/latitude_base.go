package stationxml

import (
	"fmt"
)

// Base latitude type. Because of the limitations of schema, defining this type and then extending
// it to create the real latitude type is the only way to restrict values while adding datum as an attribute.
type LatitudeBase struct {
	Float
}

func (l LatitudeBase) IsValid() error {

	if l.Unit != "" && l.Unit != "DEGREES" {
		return fmt.Errorf("invalid latitude unit: %s", l.Unit)
	}
	if l.Value < -90 || l.Value > 90 {
		return fmt.Errorf("longitude outside range: %g", l.Value)
	}

	return nil
}
