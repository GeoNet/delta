package stationxml

// Corresponds to SEED blockette 57.
type Decimation struct {
	InputSampleRate Frequency
	Factor          int32
	Offset          int32
	Delay           Float
	Correction      Float
}

func (d Decimation) IsValid() error {

	if err := Validate(d.InputSampleRate); err != nil {
		return err
	}
	if err := Validate(d.Delay); err != nil {
		return err
	}
	if err := Validate(d.Correction); err != nil {
		return err
	}

	return nil
}
