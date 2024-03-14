package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"time"
)

type Sensor struct {
	Code        string `xml:"code,attr,omitempty"`
	Model       string `xml:"model,attr,omitempty"`
	Make        string `xml:"make,attr,omitempty"`
	Type        string `xml:"type,attr,omitempty"`
	Channels    string `xml:"channels,attr,omitempty"`
	Description string `xml:"description,attr,omitempty"`
	Property    string `xml:"property,attr,omitempty"`
	Aspect      string `xml:"aspect,attr,omitempty"`

	Azimuth float64 `xml:"azimuth,attr,omitempty"`
	Dip     float64 `xml:"dip,attr,omitempty"`
	Method  string  `xml:"method,attr,omitempty"`

	Vertical float64 `xml:"vertical,attr,omitempty"`
	North    float64 `xml:"north,attr,omitempty"`
	East     float64 `xml:"east,attr,omitempty"`

	StartDate time.Time `xml:"startDate,attr,omitempty"`
	EndDate   time.Time `xml:"endDate,attr,omitempty"`
}

func (s Sensor) Less(sensor Sensor) bool {
	switch {
	case s.StartDate.Before(sensor.StartDate):
		return true
	case s.StartDate.After(sensor.StartDate):
		return false
	case s.Code < sensor.Code:
		return true
	case s.Code > sensor.Code:
		return false
	case s.Make < sensor.Make:
		return true
	case s.Make > sensor.Make:
		return false
	case s.Model < sensor.Model:
		return true
	case s.Model > sensor.Model:
		return false
	case s.Property < sensor.Property:
		return true
	case s.Property > sensor.Property:
		return false
	case s.Aspect < sensor.Aspect:
		return false
	default:
		return false
	}
}

type Site struct {
	Code string `xml:"code,attr,omitempty"`
	Name string `xml:"name,attr,omitempty"`

	Latitude       float64 `xml:"latitude,attr,omitempty"`
	Longitude      float64 `xml:"longitude,attr,omitempty"`
	Elevation      float64 `xml:"elevation,attr,omitempty"`
	Depth          float64 `xml:"depth,attr,omitempty"`
	Datum          string  `xml:"datum,attr,omitempty"`
	Survey         string  `xml:"survey,attr,omitempty"`
	RelativeHeight float64 `xml:"relativeHeight,attr,omitempty"`

	StartDate time.Time `xml:"startDate,attr,omitempty"`
	EndDate   time.Time `xml:"endDate,attr,omitempty"`

	Sensors []Sensor `xml:"Sensor,omitempty"`
}

func (s Site) Less(site Site) bool {
	return s.Code < site.Code
}

type Station struct {
	Code        string    `xml:"code,attr"`
	Name        string    `xml:"name,attr,omitempty"`
	Network     string    `xml:"network,attr,omitempty"`
	Description string    `xml:"description,attr,omitempty"`
	StartDate   time.Time `xml:"startDate,attr,omitempty"`
	EndDate     time.Time `xml:"endDate,attr,omitempty"`

	Latitude  float64 `xml:"latitude,attr,omitempty"`
	Longitude float64 `xml:"longitude,attr,omitempty"`
	Elevation float64 `xml:"elevation,attr,omitempty"`
	Depth     float64 `xml:"depth,attr,omitempty"`
	Datum     string  `xml:"datum,attr,omitempty"`

	Sites []Site `xml:"Site,omitempty"`
}

func (s Station) Less(station Station) bool {
	return s.Code < station.Code
}

type Mark struct {
	Code        string `xml:"code,attr"`
	Name        string `xml:"name,attr,omitempty"`
	Network     string `xml:"network,attr,omitempty"`
	Description string `xml:"description,attr,omitempty"`
	DomesNumber string `xml:"domesNumber,attr,omitempty"`

	Latitude           float64 `xml:"latitude,attr,omitempty"`
	Longitude          float64 `xml:"longitude,attr,omitempty"`
	Elevation          float64 `xml:"elevation,attr,omitempty"`
	GroundRelationship float64 `xml:"groundRelationship,attr,omitempty"`
	Datum              string  `xml:"datum,attr,omitempty"`

	MarkType        string  `xml:"markType,attr,omitempty"`
	MonumentType    string  `xml:"monumentType,attr,omitempty"`
	FoundationType  string  `xml:"foundationType,attr,omitempty"`
	FoundationDepth float64 `xml:"foundationDepth,attr,omitempty"`
	Bedrock         string  `xml:"bedrock,attr,omitempty"`
	Geology         string  `xml:"geology,attr,omitempty"`

	StartDate time.Time `xml:"startDate,attr,omitempty"`
	EndDate   time.Time `xml:"endDate,attr,omitempty"`

	Antennas  []Sensor `xml:"Antenna,omitempty"`
	Receivers []Sensor `xml:"Receiver,omitempty"`
}

func (m Mark) Less(mark Mark) bool {
	return m.Code < mark.Code
}

type View struct {
	Code        string `xml:"code,attr"`
	Label       string `xml:"label,attr,omitempty"`
	Description string `xml:"description,attr,omitempty"`

	Azimuth float64 `xml:"azimuth,attr,omitempty"`
	Method  string  `xml:"method,attr,omitempty"`
	Dip     float64 `xml:"dip,attr,omitempty"`

	StartDate time.Time `xml:"startDate,attr,omitempty"`
	EndDate   time.Time `xml:"endDate,attr,omitempty"`

	Sensors []Sensor `xml:"Sensor,omitempty"`
}

func (v View) Less(view View) bool {
	return v.Code < view.Code
}

type Mount struct {
	Code        string `xml:"code,attr"`
	Name        string `xml:"name,attr,omitempty"`
	Network     string `xml:"network,attr,omitempty"`
	Description string `xml:"description,attr,omitempty"`
	Mount       string `xml:"mount,attr,omitempty"`

	Latitude  float64 `xml:"latitude,attr"`
	Longitude float64 `xml:"longitude,attr"`
	Elevation float64 `xml:"elevation,attr"`
	Datum     string  `xml:"datum,attr,omitempty"`

	StartDate time.Time `xml:"startDate,attr,omitempty"`
	EndDate   time.Time `xml:"endDate,attr,omitempty"`

	Views []View `xml:"View,omitempty"`
}

func (m Mount) Less(mount Mount) bool {
	return m.Code < mount.Code
}

type Network struct {
	XMLName xml.Name `xml:"SensorXML"`

	Stations []Station `xml:"Station,omitempty"`
	Marks    []Mark    `xml:"Mark,omitempty"`
	Buoys    []Station `xml:"Buoy,omitempty"`
	Mounts   []Mount   `xml:"Mount,omitempty"`
	Samples  []Station `xml:"Sample,omitempty"`
}

func (n *Network) Marshal(wr io.Writer) error {
	enc := xml.NewEncoder(wr)
	if _, err := fmt.Fprintf(wr, "%s", xml.Header); err != nil {
		return err
	}
	if err := enc.Encode(n); err != nil {
		return err
	}
	return nil
}

func (n *Network) MarshalIndent(wr io.Writer, prefix, indent string) error {
	enc := xml.NewEncoder(wr)
	enc.Indent(prefix, indent)
	if _, err := fmt.Fprintf(wr, "%s", xml.Header); err != nil {
		return err
	}
	if err := enc.Encode(n); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(wr, "\n"); err != nil {
		return err
	}
	return nil
}
