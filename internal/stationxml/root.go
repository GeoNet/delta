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

// External maps between an External Network and individal Networks.
type External struct {
	Code        string
	Description string
	Restricted  bool

	StartDate time.Time
	EndDate   time.Time

	Networks []Network
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
						Code:        e.Code,
						Description: e.Description,
						StartDate:   e.StartDate,
						EndDate:     e.EndDate,
						Networks: []Network{{
							Code:        n.Code,
							Description: n.Description,
							Stations:    []Station{s},
						}},
					}},
				}

				return root, true
			}
		}
	}

	return Root{}, false
}
