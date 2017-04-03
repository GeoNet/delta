package stationxml

import (
	"fmt"
)

// A base node type for derivation from: Network, Station and Channel types.
type BaseNode struct {
	Code             string           `xml:"code,attr"`
	StartDate        *DateTime        `xml:"startDate,attr,omitempty" json:",omitempty"`
	EndDate          *DateTime        `xml:"endDate,attr,omitempty" json:",omitempty"`
	RestrictedStatus RestrictedStatus `xml:"restrictedStatus,attr,omitempty" json:",omitempty"`

	// A code used for display or association, alternate to the SEED-compliant code.
	AlternateCode string `xml:"alternateCode,attr,omitempty" json:",omitempty"`

	// A previously used code if different from the current code.
	HistoricalCode string `xml:"historicalCode,attr,omitempty" json:",omitempty"`

	Description string    `xml:"Description,omitempty" json:",omitempty"`
	Comments    []Comment `xml:"Comment,omitempty" json:",omitempty"`
}

func (b BaseNode) IsValid() error {

	if !(len(b.Code) > 0) {
		return fmt.Errorf("empty code element")
	}

	if b.StartDate != nil {
		if err := Validate(b.StartDate); err != nil {
			return fmt.Errorf("bad start date: %s", err)
		}
	}
	if b.EndDate != nil {
		if err := Validate(b.EndDate); err != nil {
			return fmt.Errorf("bad end date: %s", err)
		}
	}

	if err := Validate(b.RestrictedStatus); err != nil {
		return err
	}

	return nil
}
