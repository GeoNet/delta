package stationxml

import (
	"fmt"
)

// A tolerance value, measured in seconds per sample, used as a
// threshold for time error detection in data from the channel.
type ClockDrift struct {
	Float
}

func (c ClockDrift) IsValid() error {

	if err := Validate(c.Float); err != nil {
		return err
	}

	if c.Unit != "" && c.Unit != "SECONDS/SAMPLE" {
		return fmt.Errorf("invalid clock drift unit: %s", c.Unit)
	}

	return nil
}
