package stationxml

// Sensitivity and frequency ranges.
// The FrequencyRangeGroup is an optional construct that defines a pass band in Hertz
// (FrequencyStart and FrequencyEnd) in which the SensitivityValue is valid within the
// number of decibels specified in FrequencyDBVariation.
type Sensitivity struct {
	// Scalar sensitivity gain and the frequency at which it is valid.
	Gain

	// The units of the data as input from the perspective of data acquisition.
	// After correcting data for this response, these would be the resulting units.
	InputUnits Units `xml:"InputUnits"`
	// The units of the data as output from the perspective of data acquisition.
	// These would be the units of the data prior to correcting for this response.
	OutputUnits Units `xml:"OutputUnits"`
	// The frequency range for which the SensitivityValue is valid within the dB variation specified.
	FrequencyRangeGroups []FrequencyRangeGroup `xml:"FrequencyRangeGroup,omitempty" json:",omitempty"`
}

func (s Sensitivity) IsValid() error {

	if err := Validate(s.InputUnits); err != nil {
		return err
	}
	if err := Validate(s.OutputUnits); err != nil {
		return err
	}

	return nil
}
