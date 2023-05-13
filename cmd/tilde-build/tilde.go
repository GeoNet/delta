package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"sort"
	"time"
)

type Float struct {
	Value string `xml:",chardata"`
}

func toFloat(s string) *Float {
	if s == "-0" {
		return toFloat("0")
	}
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

func (s Station) Less(station Station) bool {
	return s.Code < station.Code
}

func (s *Station) Sort() {
	sort.Slice(s.Sensors, func(i, j int) bool { return s.Sensors[i].Less(s.Sensors[j]) })
}

type Domain struct {
	Name        string `xml:"name,attr"`
	Description string `xml:"description,attr,omitempty"`

	Stations []Station `xml:"Station,omitempty"`
}

func (d Domain) Less(domain Domain) bool {
	return d.Name < domain.Name
}

func (d *Domain) Sort() {
	for _, s := range d.Stations {
		s.Sort()
	}

	sort.Slice(d.Stations, func(i, j int) bool { return d.Stations[i].Less(d.Stations[j]) })
}

type Tilde struct {
	XMLName xml.Name `xml:"TildeXML"`

	Domains []Domain `xml:"Domain"`
}

func (t *Tilde) Sort() {

	for _, d := range t.Domains {
		d.Sort()
	}

	sort.Slice(t.Domains, func(i, j int) bool { return t.Domains[i].Less(t.Domains[j]) })
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
