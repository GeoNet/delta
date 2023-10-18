package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"time"
)

type Float struct {
	Value string `xml:",chardata"`
}

func toFloat(s string) *Float {
	return &Float{
		Value: s,
	}
}

func toTimePtr(t time.Time) *time.Time {
	return &t
}

type Sensor struct {
	Code        string     `xml:"code,attr,omitempty"`
	Name        string     `xml:"name,attr,omitempty"`
	Sensor      string     `xml:"sensor,attr,omitempty"`
	Type        string     `xml:"type,attr,omitempty"`
	Start       *time.Time `xml:"startDate,attr,omitempty"`
	End         *time.Time `xml:"endDate,attr,omitempty"`
	Description string     `xml:"description,attr,omitempty"`

	Latitude       *Float `xml:"Latitude,omitempty"`
	Longitude      *Float `xml:"Longitude,omitempty"`
	Elevation      *Float `xml:"Elevation,omitempty"`
	Datum          string `xml:"Datum,omitempty"`
	RelativeHeight *Float `xml:"RelativeHeight,omitempty"`

	Unit string `xml:"Unit,omitempty"`
}

func (s Sensor) Less(sensor Sensor) bool {
	switch {
	case s.Code < sensor.Code:
		return true
	case s.Code > sensor.Code:
		return false
	case s.Start != nil && sensor.Start != nil:
		return s.Start.Before(*sensor.Start)
	case s.Start == nil && sensor.Start == nil:
		return false
	case s.Start == nil:
		return true
	default:
		return false
	}
}

type Station struct {
	Code        string     `xml:"code,attr"`
	Start       *time.Time `xml:"startDate,attr,omitempty"`
	End         *time.Time `xml:"endDate,attr,omitempty"`
	Description string     `xml:"description,attr,omitempty"`

	Latitude  *Float `xml:"Latitude,omitempty"`
	Longitude *Float `xml:"Longitude,omitempty"`
	Elevation *Float `xml:"Elevation,omitempty"`
	Datum     string `xml:"Datum,omitempty"`

	Sensors []Sensor `xml:"Sensor,omitempty"`
}

type Domain struct {
	Name        string `xml:"name,attr"`
	Description string `xml:"description,attr,omitempty"`

	Stations []Station `xml:"Station,omitempty"`
}

type Tilde struct {
	XMLName xml.Name `xml:"TildeXML"`

	Domains []Domain `xml:"Domain"`
}

func (t *Tilde) Marshal(wr io.Writer) error {
	enc := xml.NewEncoder(wr)
	if _, err := fmt.Fprintf(wr, "%s", xml.Header); err != nil {
		return err
	}
	if err := enc.Encode(t); err != nil {
		return err
	}
	return nil
}

func (t *Tilde) MarshalIndent(wr io.Writer, prefix, indent string) error {
	enc := xml.NewEncoder(wr)
	enc.Indent(prefix, indent)
	if _, err := fmt.Fprintf(wr, "%s", xml.Header); err != nil {
		return err
	}
	if err := enc.Encode(t); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(wr, "\n"); err != nil {
		return err
	}
	return nil
}
