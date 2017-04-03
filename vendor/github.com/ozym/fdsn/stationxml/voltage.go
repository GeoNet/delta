package stationxml

import (
	"fmt"
)

// A time value in seconds.
type Voltage struct {
	Float
}

func (v Voltage) IsValid() error {

	if err := Validate(v.Float); err != nil {
		return err
	}

	if v.Unit != "" && v.Unit != "VOLTS" {
		return fmt.Errorf("invalid unit: %s", v.Unit)
	}

	return nil
}
