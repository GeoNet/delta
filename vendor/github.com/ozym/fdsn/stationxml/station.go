package stationxml

import (
	"fmt"
)

// This type represents a Station epoch.
// It is common to only have a single station epoch with the station's creation
// and termination dates as the epoch start and end dates.
type Station struct {
	BaseNode

	Latitude  Latitude  `xml:"Latitude"`
	Longitude Longitude `xml:"Longitude"`
	Elevation Distance  `xml:"Elevation"`

	// These fields describe the location of the station using geopolitical
	// entities (country, city, etc.).
	Site Site `xml:"Site"`

	// Type of vault, e.g. WWSSN, tunnel, transportable array, etc.
	Vault string `xml:",omitempty" json:",omitempty"`

	// Type of rock and/or geologic formation.
	Geology string `xml:",omitempty" json:",omitempty"`

	// Equipment used by all channels at a station.
	Equipments []Equipment `xml:"Equipment,omitempty" json:",omitempty"`

	// An operating agency and associated contact persons.
	// If there multiple operators, each one should be encapsulated within an Operator tag.
	// Since the Contact element is a generic type that represents any contact person,
	// it also has its own optional Agency element.
	Operators []Operator `xml:"Operator,omitempty" json:",omitempty"`

	// Date and time (UTC) when the station was first installed.
	CreationDate DateTime `xml:"CreationDate"`

	// Date and time (UTC) when the station was terminated or will be terminated.
	// A blank value should be assumed to mean that the station is still active.
	TerminationDate *DateTime `xml:",omitempty" json:",omitempty"`

	// Total number of channels recorded at this station.
	TotalNumberChannels Counter `xml:",omitempty" json:",omitempty"`

	// Number of channels recorded at this station and selected by the query
	// that produced this document.
	SelectedNumberChannels Counter `xml:",omitempty" json:",omitempty"`

	// URI of any type of external report, such as IRIS data reports or dataless SEED volumes.
	ExternalReferences []ExternalReference `xml:"ExternalReference,omitempty" json:",omitempty"`

	Channels []Channel `xml:"Channel,omitempty" json:",omitempty"`
}

func (s Station) IsValid() error {

	if err := Validate(s.BaseNode); err != nil {
		return err
	}

	if err := Validate(s.Latitude); err != nil {
		return err
	}
	if err := Validate(s.Longitude); err != nil {
		return err
	}
	if err := Validate(s.Elevation); err != nil {
		return err
	}
	if err := Validate(s.Site); err != nil {
		return err
	}

	for _, e := range s.Equipments {
		if err := Validate(e); err != nil {
			return err
		}
	}

	for _, o := range s.Operators {
		if err := Validate(o); err != nil {
			return err
		}
	}

	if s.CreationDate.IsZero() {
		return fmt.Errorf("missing creation date")
	}

	if err := Validate(s.CreationDate); err != nil {
		return err
	}

	if s.TerminationDate != nil {
		if err := Validate(s.TerminationDate); err != nil {
			return err
		}
	}

	for _, x := range s.ExternalReferences {
		if err := Validate(x); err != nil {
			return err
		}
	}

	for _, c := range s.Channels {
		if err := Validate(c); err != nil {
			return err
		}
	}

	return nil
}
