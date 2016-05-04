package stationxml

import (
	"fmt"
)

// A type to document units. Corresponds to SEED blockette 34.
type Units struct {
	// Name of units, e.g. "M/S", "V", "PA".
	Name string
	// Description of units, e.g. "Velocity in meters per second", "Volts", "Pascals".
	Description string `xml:",omitempty" json:",omitempty"`
}

func (u Units) IsValid() error {

	if !(len(u.Name) > 0) {
		return fmt.Errorf("empty units name")
	}

	return nil
}
