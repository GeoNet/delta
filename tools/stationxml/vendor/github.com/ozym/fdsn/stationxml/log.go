package stationxml

// Container for log entries.
type Log struct {
	Entries []Comment `xml:"Entry,omitempty" json:",omitempty"`
}

func (l Log) IsValid() error {

	for _, c := range l.Entries {
		if err := Validate(c); err != nil {
			return err
		}
	}

	return nil
}
