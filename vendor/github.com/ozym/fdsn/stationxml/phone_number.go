package stationxml

import (
	"fmt"
	"regexp"
)

type PhoneNumber struct {
	Description string `xml:"description,attr,omitempty" json:",omitempty"`

	CountryCode int32 `xml:",omitempty" json:",omitempty"`
	AreaCode    int32

	// Pattern "[0-9]+-[0-9]+"
	PhoneNumber string
}

func (p PhoneNumber) IsValid() error {

	if !(len(p.PhoneNumber) > 0) {
		return fmt.Errorf("empty phone number")
	}

	if p.AreaCode == 0 {
		return fmt.Errorf("no area code")
	}

	if !(regexp.MustCompile(`^[0-9]+-[0-9]+$`).MatchString(p.PhoneNumber)) {
		return fmt.Errorf("bad phone number: %s", p.PhoneNumber)
	}

	return nil
}
