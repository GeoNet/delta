package stationxml

// This is a group of elements that represent sample rate.
// If this group is included, then SampleRate, which is the sample rate in samples per second, is required.
// SampleRateRatio, which is expressed as a ratio of number of samples in a number of seconds, is optional.
// If both are included, SampleRate should be considered more definitive.
type SampleRateGroup struct {
	SampleRate      SampleRate
	SampleRateRatio *SampleRateRatio `xml:",omitempty" json:",omitempty"`
}

func (s SampleRateGroup) IsValid() error {

	if err := Validate(s.SampleRate); err != nil {
		return err
	}

	return nil
}
