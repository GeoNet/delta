package stationxml

import (
	"fmt"
)

// The BaseFilter Type is derived by all filters
type BaseFilter struct {
	// Same meaning as Equipment resourceId
	ResourceId string `xml:"resourceId,attr"`

	// A name given to this filter.
	Name string `xml:"name,attr"`

	Description string `xml:",omitempty" json:",omitempty"`

	// The units of the data as input from the perspective of data acquisition.
	// After correcting data for this response, these would be the resulting units.
	InputUnits Units

	// The units of the data as output from the perspective of data acquisition.
	// These would be the units of the data prior to correcting for this response.
	OutputUnits Units
}

func (f BaseFilter) IsValid() error {

	if !(len(f.ResourceId) > 0) {
		return fmt.Errorf("empty filter resourceid")
	}
	if !(len(f.Name) > 0) {
		return fmt.Errorf("empty fir name")
	}

	if err := Validate(f.InputUnits); err != nil {
		return err
	}
	if err := Validate(f.OutputUnits); err != nil {
		return err
	}

	return nil
}
