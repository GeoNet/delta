// Package tides provides an embeddable mechanism for providing tidal constituents for given gauge sites
// that can be used for tide predication.
package tides

import (
	"fmt"
	"strings"
)

//go:generate bash -c "go run generate/*.go | gofmt -s > auto.go; test -s auto.go || rm auto.go"

// Constitutent stores the calculated amplotude and lag for a named tidal phase, the amplitude
// and lag are given in units of meters and degrees respectively.
type Constituent struct {
	Name      string
	Amplitude float64
	Lag       float64
}

// String provides a standard representation of the tidal constituent.
func (c Constituent) String() string {
	return fmt.Sprintf("%s/%g/%g", c.Name, c.Amplitude, c.Lag)
}

// Tide provides the general parameters needed to predict tides at a given site and the
// associated tidal consitituents. The TimeZone, Latitude, and Longitude are expected
// to be the parameters used to generated the tidal prediction constituents and may
// differ from the geographic values recorded elsewhere.
type Tide struct {
	Code      string
	Network   string
	Number    string
	TimeZone  float64
	Latitude  float64
	Longitude float64
	Crex      string

	Constituents []Constituent
}

// Zone provides a conversion of the time zone parameter for use with common tidal
// prediction software.
func (t Tide) Zone() float64 {
	return (360.0 - t.TimeZone) / 15.0
}

// Lookup will return a Tide pointer for the requested site code.
// A nil pointer will be returned if a code cannot be found.
func Lookup(code string) *Tide {
	if t, ok := _tides[strings.ToUpper(code)]; ok {
		return &t
	}
	return nil
}
