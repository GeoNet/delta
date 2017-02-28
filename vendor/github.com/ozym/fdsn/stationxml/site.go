package stationxml

import (
	"fmt"
)

// Description of a site location using name and optional geopolitical boundaries (country, city, etc.).
type Site struct {
	// The commonly used name of this station, equivalent to the SEED blockette 50, field 9.
	Name string

	// A longer description of the location of this station, e.g.
	// "NW corner of Yellowstone National Park" or "20 miles west of Highway 40."
	Description string `xml:",omitempty" json:",omitempty"`

	// The town or city closest to the station.
	Town   string `xml:",omitempty" json:",omitempty"`
	County string `xml:",omitempty" json:",omitempty"`

	// The state, province, or region of this site.
	Region  string `xml:",omitempty" json:",omitempty"`
	Country string `xml:",omitempty" json:",omitempty"`
}

func (s Site) IsValid() error {

	if !(len(s.Name) > 0) {
		return fmt.Errorf("empty site name")
	}

	return nil
}
