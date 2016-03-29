package stationxml

import (
	"fmt"
)

type Frequency struct {
	Float
}

func (f Frequency) IsValid() error {

	if err := Validate(f.Float); err != nil {
		return err
	}

	if f.Float.Unit != "" && f.Float.Unit != "HERTZ" {
		return fmt.Errorf("invalid unit: %s", f.Unit)
	}

	return nil
}
