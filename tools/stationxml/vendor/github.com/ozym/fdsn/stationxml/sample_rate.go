package stationxml

import (
	"fmt"
)

// Sample rate in samples per second.
type SampleRate struct {
	Float
}

func (s SampleRate) IsValid() error {

	if err := Validate(s.Float); err != nil {
		return err
	}

	if s.Float.Unit != "" && s.Float.Unit != "SAMPLES/S" {
		return fmt.Errorf("invalid sample rate unit: %s", s.Unit)
	}

	return nil
}
