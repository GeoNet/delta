package stationxml

import (
	"fmt"
)

type Operator struct {
	Agencies []string `xml:"Agency,omitempty" json:",omitempty"`
	Contacts []Person `xml:"Contact,omitempty" json:",omitempty"`
	WebSites []AnyURI `xml:"WebSite,omitempty" json:",omitempty"`
}

func (o Operator) IsValid() error {

	for _, a := range o.Agencies {
		if !(len(a) > 0) {
			return fmt.Errorf("empty operator agency")
		}
	}
	for _, c := range o.Contacts {
		if err := Validate(c); err != nil {
			return err
		}
	}
	for _, w := range o.WebSites {
		if !(len(w) > 0) {
			return fmt.Errorf("empty websites uri")
		}
	}

	return nil
}
