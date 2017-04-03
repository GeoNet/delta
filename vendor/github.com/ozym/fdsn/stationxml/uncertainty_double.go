package stationxml

import (
	"fmt"
)

// Expressing uncertainties or errors with a positive and a negative component.
// Both values should be given as positive integers, but minus_error is understood to actually be negative.
type UncertaintyDouble struct {
	PlusError  float64 `xml:"plusError,attr,omitempty" json:",omitempty"`
	MinusError float64 `xml:"minusError,attr,omitempty" json:",omitempty"`
}

func (f UncertaintyDouble) IsValid() error {

	if f.PlusError < 0.0 {
		return fmt.Errorf("float plus error shouldn't be negative: %g", f.PlusError)
	}
	if f.MinusError < 0.0 {
		return fmt.Errorf("float minus error shouldn't be negative: %g", f.MinusError)
	}

	return nil
}
