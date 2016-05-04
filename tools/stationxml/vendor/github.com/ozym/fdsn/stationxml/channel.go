package stationxml

import (
	"fmt"
)

// Equivalent to SEED blockette 52 and parent element for the related the response blockettes.
type Channel struct {
	BaseNode

	// URI of any type of external report, such as data quality reports.
	ExternalReferences []ExternalReference `xml:"ExternalReference,omitempty" json:",omitempty"`

	LocationCode string `xml:"locationCode,attr"`

	// Latitude coordinate of this channel's sensor.
	Latitude Latitude

	//Longitude coordinate of this channel's sensor.
	Longitude Longitude

	// Elevation of the sensor.
	Elevation Distance

	// The local depth or overburden of the instrument's location.
	// For downhole instruments, the depth of the instrument under the surface ground level.
	// For underground vaults, the distance from the instrument to the local ground level above.
	Depth Distance

	// Azimuth of the sensor in degrees from north, clockwise.
	Azimuth *Azimuth `xml:",omitempty" json:",omitempty"`

	// Dip of the instrument in degrees, down from horizontal
	Dip *Dip `xml:",omitempty" json:",omitempty"`

	// The type of data this channel collects. Corresponds to channel flags in SEED blockette 52.
	// The SEED volume producer could use the first letter of an Output value as the SEED channel flag.
	Types []Type `xml:"Type,omitempty" json:",omitempty"`

	// This is a group of elements that represent sample rate.
	// If this group is included, then SampleRate, which is the sample rate in samples per second, is required.
	// SampleRateRatio, which is expressed as a ratio of number of samples in a number of seconds, is optional.
	// If both are included, SampleRate should be considered more definitive.
	SampleRateGroup

	// The storage format of the recorded data (e.g. SEED).
	StorageFormat string
	// A tolerance value, measured in seconds per sample, used as a threshold for time
	// error detection in data from the channel.
	ClockDrift *ClockDrift `xml:",omitempty" json:",omitempty"`

	CalibrationUnits *Units     `xml:",omitempty" json:",omitempty"`
	Sensor           *Equipment `xml:",omitempty" json:",omitempty"`
	PreAmplifier     *Equipment `xml:",omitempty" json:",omitempty"`
	DataLogger       *Equipment `xml:",omitempty" json:",omitempty"`
	Equipment        *Equipment `xml:",omitempty" json:",omitempty"`
	Response         *Response  `xml:",omitempty" json:",omitempty"`
}

func (c Channel) IsValid() error {

	if err := Validate(c.BaseNode); err != nil {
		return err
	}

	if err := Validate(c.Latitude); err != nil {
		return err
	}
	if err := Validate(c.Longitude); err != nil {
		return err
	}
	if err := Validate(c.Elevation); err != nil {
		return err
	}
	if err := Validate(c.Depth); err != nil {
		return err
	}

	if err := Validate(c.Dip); err != nil {
		return err
	}
	if err := Validate(c.Azimuth); err != nil {
		return err
	}

	for _, t := range c.Types {
		if err := Validate(t); err != nil {
			return err
		}
	}

	if err := Validate(c.SampleRateGroup); err != nil {
		return nil
	}

	if !(len(c.StorageFormat) > 0) {
		return fmt.Errorf("empty code element")
	}

	if c.ClockDrift != nil {
		if err := Validate(c.ClockDrift); err != nil {
			return nil
		}
	}

	if c.CalibrationUnits != nil {
		if err := Validate(c.CalibrationUnits); err != nil {
			return nil
		}
	}

	if c.Sensor != nil {
		if err := Validate(c.Sensor); err != nil {
			return nil
		}
	}
	if c.PreAmplifier != nil {
		if err := Validate(c.PreAmplifier); err != nil {
			return nil
		}
	}
	if c.DataLogger != nil {
		if err := Validate(c.DataLogger); err != nil {
			return nil
		}
	}
	if c.Equipment != nil {
		if err := Validate(c.Equipment); err != nil {
			return nil
		}
	}
	if c.Response != nil {
		if err := Validate(c.Response); err != nil {
			return nil
		}
	}

	return nil
}
