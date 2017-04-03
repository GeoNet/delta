package stationxml

import (
	"fmt"
)

// Representation of a person's contact information.
// A person can belong to multiple agencies and have multiple email addresses and phone numbers.
type Person struct {
	Names        []string      `xml:"Name,omitempty" json:",omitempty"`
	Agencies     []string      `xml:"Agency,omitempty" json:",omitempty"`
	Email        []Email       `xml:"Email,omitempty" json:",omitempty"`
	PhoneNumbers []PhoneNumber `xml:"Phone,omitempty" json:",omitempty"`
}

func (p Person) IsValid() error {

	for _, n := range p.Names {
		if !(len(n) > 0) {
			return fmt.Errorf("empty person name")
		}
	}
	for _, a := range p.Agencies {
		if !(len(a) > 0) {
			return fmt.Errorf("empty person agency")
		}
	}
	for _, e := range p.Email {
		if err := Validate(e); err != nil {
			return err
		}
	}
	for _, x := range p.PhoneNumbers {
		if err := Validate(x); err != nil {
			return err
		}
	}

	return nil
}
