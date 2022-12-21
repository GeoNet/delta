package stationxml

import (
	"time"
)

// Equipment describes a StationXML Equipment element
type Equipment struct {
	Type             string
	Description      string
	Manufacturer     string
	Model            string
	SerialNumber     string
	InstallationDate time.Time
	RemovalDate      time.Time
	Response         string
}

// Stream forms the main part of an individual StationXML Channel element
type Stream struct {
	Code string

	SamplingRate float64
	Triggered    bool
	Types        string

	Vertical float64
	Azimuth  float64
	Dip      float64

	Datalogger Equipment
	Sensor     Equipment

	StartDate time.Time
	EndDate   time.Time

	Response *ResponseType
}

// Channel forms the main part of a set of StationXML Channel elements
type Channel struct {
	LocationCode string

	Latitude  float64
	Longitude float64
	Elevation float64
	Survey    string
	Datum     string

	Streams []Stream
}

// Station forms the main part of a StationXML Station element.
type Station struct {
	Code        string
	Name        string
	Description string

	Latitude  float64
	Longitude float64
	Elevation float64
	Survey    string
	Datum     string

	StartDate time.Time
	EndDate   time.Time

	CreationDate    time.Time
	TerminationDate time.Time

	Channels []Channel
}

// Network forms the main part of a StationXML Network element.
type Network struct {
	Code        string
	Description string
	Restricted  bool

	Stations []Station
}

// Start returns the earliest Station start time in a Network.
func (n Network) Start() time.Time {
	var start time.Time
	for _, s := range n.Stations {
		if start.IsZero() || s.StartDate.Before(start) {
			start = s.StartDate
		}
	}
	return start
}

// End returns the latest Station end time in a Network.
func (n Network) End() time.Time {
	var end time.Time
	for _, s := range n.Stations {
		if end.IsZero() || s.EndDate.After(end) {
			end = s.EndDate
		}
	}
	return end
}

// External maps between an External Network and individal Networks.
type External struct {
	Code        string
	Description string
	Restricted  bool

	Networks []Network
}

// Start returns the earliest Station start time in an External Network.
func (e External) Start() time.Time {
	var start time.Time
	for _, n := range e.Networks {
		if t := n.Start(); start.IsZero() || t.Before(start) {
			start = t
		}
	}
	return start
}

// End returns the latest Station end time in an External Network.
func (e External) End() time.Time {
	var end time.Time
	for _, n := range e.Networks {
		if t := n.End(); end.IsZero() || t.After(end) {
			end = t
		}
	}
	return end
}

// Root describes the standard StationXML layout which can be used as the barebones for building version specific encoders.
type Root struct {
	Source string
	Sender string
	Module string
	Create bool

	Externals []External
}

// ExternalCode returns the network code of the first External entry in the Root structure, this is aimed at building single entry file names.
func (r Root) ExternalCode() string {
	for _, e := range r.Externals {
		return e.Code
	}
	return ""
}

// NetworkCode returns the network code of the first Network entry in the Root structure, this is aimed at building single entry file names.
func (r Root) NetworkCode() string {
	for _, e := range r.Externals {
		for _, n := range e.Networks {
			return n.Code
		}
	}
	return ""
}

// StationCode returns the station code of the first Station entry in the Root structure, this is aimed at building single entry file names.
func (r Root) StationCode() string {
	for _, e := range r.Externals {
		for _, n := range e.Networks {
			for _, s := range n.Stations {
				return s.Code
			}
		}
	}
	return ""
}

// Single builds a Root structure for the given station code.
func (r Root) Single(code string) (Root, bool) {
	for _, e := range r.Externals {
		for _, n := range e.Networks {
			for _, s := range n.Stations {
				if s.Code != code {
					continue
				}

				root := Root{
					Source: r.Source,
					Sender: r.Sender,
					Module: r.Module,
					Create: r.Create,

					Externals: []External{{
						Code: e.Code,
						Networks: []Network{{
							Code:     n.Code,
							Stations: []Station{s},
						}},
					}},
				}

				return root, true
			}
		}
	}

	return Root{}, false
}
