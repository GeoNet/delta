package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"time"
)

type Sensor struct {
	Code        string `xml:"code,attr,omitempty" json:"code,omitempty"`
	Model       string `xml:"model,attr,omitempty" json:"model,omitempty"`
	Make        string `xml:"make,attr,omitempty" json:"make,omitempty"`
	Type        string `xml:"type,attr,omitempty" json:"type,omitempty"`
	Channels    string `xml:"channels,attr,omitempty" json:"channels,omitempty"`
	Description string `xml:"description,attr,omitempty" json:"description,omitempty"`
	Property    string `xml:"property,attr,omitempty" json:"property,omitempty"`
	Aspect      string `xml:"aspect,attr,omitempty" json:"aspect,omitempty"`

	Azimuth float64 `xml:"azimuth,attr,omitempty" json:"azimuth,omitempty"`
	Dip     float64 `xml:"dip,attr,omitempty" json:"dip,omitempty"`
	Method  string  `xml:"method,attr,omitempty" json:"method,omitempty"`

	Vertical float64 `xml:"vertical,attr,omitempty" json:"vertical,omitempty"`
	North    float64 `xml:"north,attr,omitempty" json:"north,omitempty"`
	East     float64 `xml:"east,attr,omitempty" json:"east,omitempty"`

	StartDate time.Time `xml:"startDate,attr,omitempty" json:"start-date,omitempty"`
	EndDate   time.Time `xml:"endDate,attr,omitempty" json:"end-date,omitempty"`
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
	Code string `xml:"code,attr,omitempty" json:"code,omitempty"`
	Name string `xml:"name,attr,omitempty" json:"name,omitempty"`

	Latitude       float64 `xml:"latitude,attr,omitempty" json:"latitude,omitempty"`
	Longitude      float64 `xml:"longitude,attr,omitempty" json:"longitude,omitempty"`
	Elevation      float64 `xml:"elevation,attr,omitempty" json:"elevation,omitempty"`
	Depth          float64 `xml:"depth,attr,omitempty" json:"depth,omitempty"`
	Datum          string  `xml:"datum,attr,omitempty" json:"datum,omitempty"`
	Survey         string  `xml:"survey,attr,omitempty" json:"survey,omitempty"`
	RelativeHeight float64 `xml:"relativeHeight,attr,omitempty" json:"relative-height,omitempty"`

	StartDate time.Time `xml:"startDate,attr,omitempty" json:"start-date,omitempty"`
	EndDate   time.Time `xml:"endDate,attr,omitempty" json:"end-date,omitempty"`

	Features []Sensor `xml:"Feature,omitempty" json:"feature,omitempty"`
	Sensors  []Sensor `xml:"Sensor,omitempty" json:"sensor,omitempty"`
}

func (s Site) Less(site Site) bool {
	return s.Code < site.Code
}

type Station struct {
	Code     string `xml:"code,attr" json:"code,omitempty"`
	Network  string `xml:"network,attr,omitempty" json:"network,omitempty"`
	External string `xml:"external,attr,omitempty" json:"external,omitempty"`

	Name        string    `xml:"name,attr,omitempty" json:"name,omitempty"`
	Description string    `xml:"description,attr,omitempty" json:"description,omitempty"`
	StartDate   time.Time `xml:"startDate,attr,omitempty" json:"start-date,omitempty"`
	EndDate     time.Time `xml:"endDate,attr,omitempty" json:"end-date,omitempty"`

	Latitude  float64 `xml:"latitude,attr,omitempty" json:"latitude,omitempty"`
	Longitude float64 `xml:"longitude,attr,omitempty" json:"longitude,omitempty"`
	Elevation float64 `xml:"elevation,attr,omitempty" json:"elevation,omitempty"`
	Depth     float64 `xml:"depth,attr,omitempty" json:"depth,omitempty"`
	Datum     string  `xml:"datum,attr,omitempty" json:"datum,omitempty"`

	Sites []Site `xml:"Site,omitempty" json:"site,omitempty"`
}

func (s Station) Less(station Station) bool {
	switch {
	case s.Code < station.Code:
		return true
	case s.Code > station.Code:
		return false
	case s.Network < station.Network:
		return true
	default:
		return false
	}
}

type Mark struct {
	Code    string `xml:"code,attr" json:"code,omitempty"`
	Network string `xml:"network,attr,omitempty" json:"network,omitempty"`

	Name        string `xml:"name,attr,omitempty" json:"name,omitempty"`
	Description string `xml:"description,attr,omitempty" json:"description,omitempty"`
	DomesNumber string `xml:"domesNumber,attr,omitempty" json:"domes-number,omitempty"`

	Latitude           float64 `xml:"latitude,attr,omitempty" json:"latitude,omitempty"`
	Longitude          float64 `xml:"longitude,attr,omitempty" json:"longitude,omitempty"`
	Elevation          float64 `xml:"elevation,attr,omitempty" json:"elevation,omitempty"`
	GroundRelationship float64 `xml:"groundRelationship,attr,omitempty" json:"ground-relationship,omitempty"`
	Datum              string  `xml:"datum,attr,omitempty" json:"datum,omitempty"`

	MarkType        string  `xml:"markType,attr,omitempty" json:"mark-type,omitempty"`
	MonumentType    string  `xml:"monumentType,attr,omitempty" json:"monument-type,omitempty"`
	FoundationType  string  `xml:"foundationType,attr,omitempty" json:"foundation-type,omitempty"`
	FoundationDepth float64 `xml:"foundationDepth,attr,omitempty" json:"foundation-depth,omitempty"`
	Bedrock         string  `xml:"bedrock,attr,omitempty" json:"bedrock,omitempty"`
	Geology         string  `xml:"geology,attr,omitempty" json:"geology,omitempty"`

	StartDate time.Time `xml:"startDate,attr,omitempty" json:"start-date,omitempty"`
	EndDate   time.Time `xml:"endDate,attr,omitempty" json:"end-date,omitempty"`

	Antennas  []Sensor `xml:"Antenna,omitempty" json:"antenna,omitempty"`
	Receivers []Sensor `xml:"Receiver,omitempty" json:"receiver,omitempty"`
}

func (m Mark) Less(mark Mark) bool {
	return m.Code < mark.Code
}

type View struct {
	Code        string `xml:"code,attr" json:"code,omitempty"`
	Label       string `xml:"label,attr,omitempty" json:"label,omitempty"`
	Description string `xml:"description,attr,omitempty" json:"description,omitempty"`

	Azimuth float64 `xml:"azimuth,attr,omitempty" json:"azimuth,omitempty"`
	Method  string  `xml:"method,attr,omitempty" json:"method,omitempty"`
	Dip     float64 `xml:"dip,attr,omitempty" json:"dip,omitempty"`

	StartDate time.Time `xml:"startDate,attr,omitempty" json:"start-date,omitempty"`
	EndDate   time.Time `xml:"endDate,attr,omitempty" json:"end-date,omitempty"`

	Sensors []Sensor `xml:"Sensor,omitempty" json:"sensor,omitempty"`
}

func (v View) Less(view View) bool {
	return v.Code < view.Code
}

type Mount struct {
	Code     string `xml:"code,attr" json:"code,omitempty"`
	Network  string `xml:"network,attr,omitempty" json:"network,omitempty"`
	External string `xml:"external,attr,omitempty" json:"external,omitempty"`

	Name        string `xml:"name,attr,omitempty" json:"name,omitempty"`
	Description string `xml:"description,attr,omitempty" json:"description,omitempty"`
	Mount       string `xml:"mount,attr,omitempty" json:"mount,omitempty"`

	Latitude  float64 `xml:"latitude,attr" json:"latitude,omitempty"`
	Longitude float64 `xml:"longitude,attr" json:"longitude,omitempty"`
	Elevation float64 `xml:"elevation,attr" json:"elevation,omitempty"`
	Datum     string  `xml:"datum,attr,omitempty" json:"datum,omitempty"`

	StartDate time.Time `xml:"startDate,attr,omitempty" json:"start-date,omitempty"`
	EndDate   time.Time `xml:"endDate,attr,omitempty" json:"end-date,omitempty"`

	Views []View `xml:"View,omitempty" json:"view,omitempty"`
}

func (m Mount) Less(mount Mount) bool {
	return m.Code < mount.Code
}

type Group struct {
	Name string `xml:"name,omitempty,attr" json:"name,omitempty"`

	Marks    []Mark    `xml:"Mark,omitempty" json:"marks,omitempty"`
	Mounts   []Mount   `xml:"Mount,omitempty" json:"mounts,omitempty"`
	Samples  []Station `xml:"Sample,omitempty" json:"samples,omitempty"`
	Stations []Station `xml:"Station,omitempty" json:"stations,omitempty"`
}

type Network struct {
	XMLName xml.Name `xml:"SensorXML"`

	Groups []Group `xml:"Group,omitempty" json:"group,omitempty"`

	Stations []Station `xml:"Station,omitempty" json:"station,omitempty"`
	Marks    []Mark    `xml:"Mark,omitempty" json:"mark,omitempty"`
	Buoys    []Station `xml:"Buoy,omitempty" json:"buoy,omitempty"`
	Mounts   []Mount   `xml:"Mount,omitempty" json:"mount,omitempty"`
	Samples  []Station `xml:"Sample,omitempty" json:"sample,omitempty"`
}

func (n *Network) EncodeXML(wr io.Writer, prefix, indent string) error {
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

func (n *Network) EncodeJSON(wr io.Writer) error {
	// remap network to avoid xml specific details
	remap := struct {
		Groups   []Group   `json:"group,omitempty"`
		Stations []Station `json:"station,omitempty"`
		Marks    []Mark    `json:"mark,omitempty"`
		Buoys    []Station `json:"buoy,omitempty"`
		Mounts   []Mount   `json:"mount,omitempty"`
		Samples  []Station `json:"sample,omitempty"`
	}{
		Groups:   n.Groups,
		Stations: n.Stations,
		Marks:    n.Marks,
		Buoys:    n.Buoys,
		Mounts:   n.Mounts,
		Samples:  n.Samples,
	}

	return json.NewEncoder(wr).Encode(remap)
}
