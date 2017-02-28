package stationxml

import (
	"fmt"
)

// A time value in seconds.
type Second struct {
	Float
}

func (s Second) IsValid() error {

	if err := Validate(s.Float); err != nil {
		return err
	}

	if s.Unit != "" && s.Unit != "SECONDS" {
		return fmt.Errorf("invalid unit: %s", s.Unit)
	}

	return nil
}
