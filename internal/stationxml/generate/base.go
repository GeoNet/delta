package main

import (
	"io"
	"text/template"
)

const baseTemplate = `
package stationxml

import (
	"bytes"
	"encoding/xml"
)

const (
	fixFrom = {{bt}}xmlns:xsi="xsi"{{bt}}
	fixTo = {{bt}}xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"{{bt}}
)

// there's a golang missing feature https://github.com/golang/go/issues/13400 related to names spaces,
// this is a crude workaround to simply swap in and out the bad marshalling.
func fix(from []byte) []byte {
	return bytes.ReplaceAll(from, []byte(fixFrom), []byte(fixTo))
}

// add a workaround base URI type
type AnyURI string

// FDSNStationXML is a wrapper around the RootType to help with XML marshalling.
type FDSNStationXML struct {
	XMLName   xml.Name {{bt}}xml:"http://www.fdsn.org/xml/station/1 FDSNStationXML"{{bt}}
	SchemaLocation string {{bt}}xml:"xsi schemaLocation,attr"{{bt}}

	RootType
}

// NewFDSNStationXML returns an FDSNStationXML value that is wrapping the given RootType.
func NewFDSNStationXML(root RootType) FDSNStationXML{
	return FDSNStationXML{
		SchemaLocation: "http://www.fdsn.org/xml/station/1 http://www.fdsn.org/xml/station/fdsn-station-1.2.xsd",
		RootType: root,
	}
}

func (x FDSNStationXML) Marshal() ([]byte, error) {
	s, err := xml.Marshal(x)
	if err != nil {
		return nil, err
	}
	return append([]byte(xml.Header), append(fix(s), '\n')...), nil
}

func (x FDSNStationXML) MarshalIndent(prefix, indent string) ([]byte, error) {
	s, err := xml.MarshalIndent(x, prefix, indent)
	if err != nil {
		return nil, err
	}
	return append([]byte(xml.Header), append(fix(s), '\n')...), nil
}
`

const testTemplate = `
package stationxml

import (
	"testing"
)

const exampleXML ={{bt}}<?xml version="1.0" encoding="UTF-8"?>
<FDSNStationXML xmlns="http://www.fdsn.org/xml/station/1" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.fdsn.org/xml/station/1 http://www.fdsn.org/xml/station/fdsn-station-{{.Version}}.xsd" schemaVersion="{{.Version}}">
  <Source>Source</Source>
  <Sender>Sender</Sender>
  <Module>Module</Module>
  <ModuleURI>ModuleURI</ModuleURI>
  <Created>2022-07-05T21:54:16.7130Z</Created>
</FDSNStationXML>
{{bt}}

func TestSchema(t *testing.T) {

	base := FDSNStationXML{
		SchemaLocation: "http://www.fdsn.org/xml/station/1 http://www.fdsn.org/xml/station/fdsn-station-1.2.xsd",
		RootType: RootType{
			SchemaVersion: {{.Version}},
			Source: "Source",
			Sender: "Sender",
			Module: "Module",
			ModuleURI: "ModuleURI",
			Created: MustParseDateTime("2022-07-05T21:54:16.7130Z"),
		},
	}

	data, err := base.MarshalIndent("", "  ")
	if err != nil {
		t.Fatal(err)
	}

	if s := string(data); s != exampleXML {
		t.Errorf("invalid test xml, expected %s, got %s", exampleXML, s)
	}
}

func TestNewSchema(t *testing.T) {

	base := NewFDSNStationXML(RootType{
		SchemaVersion: {{.Version}},
		Source: "Source",
		Sender: "Sender",
		Module: "Module",
		ModuleURI: "ModuleURI",
		Created: MustParseDateTime("2022-07-05T21:54:16.7130Z"),
	})

	data, err := base.MarshalIndent("", "  ")
	if err != nil {
		t.Fatal(err)
	}

	if s := string(data); s != exampleXML {
		t.Errorf("invalid test xml, expected %s, got %s", exampleXML, s)
	}
}
`

const docTemplate = `
{{range $c := .Comments -}}
// {{$c}}
{{end -}}
package stationxml
`

const dateTimeTemplate = `package stationxml

import (
	"encoding/xml"
	"fmt"
	"time"
)

const DateTimeFormat = "2006-01-02T15:04:05.0000Z"

type DateTime struct {
	time.Time {{bt}}xml:",chardata"{{bt}}
}

func Now() DateTime {
	return DateTime{Time: time.Now()}
}

func NewDateTime(t time.Time) DateTime {
	return DateTime{
		Time: t,
	}
}

func ParseDateTime(s string) (DateTime, error) {
	x, err := time.Parse(DateTimeFormat, s)
	return DateTime{x}, err
}

func MustParseDateTime(s string) DateTime {
	x, err := time.Parse(DateTimeFormat, s)
	if err != nil {
		panic(err)
	}
	return DateTime{x}
}
func MustParseDateTimePtr(s string) *DateTime {
	x := MustParseDateTime(s)
	return &x
}

func (t DateTime) IsValid() error {
	if !t.Time.IsZero() && t.Time.Year() < 1880 {
		return fmt.Errorf("incorrect date: %s", t.String())
	}
	return nil
}

func (t DateTime) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if t.Time.Year() < 1880 {
		return e.EncodeElement(nil, start)
	}
	return e.EncodeElement(t.Time.Format(DateTimeFormat), start)
}

func (t *DateTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {

	var s string
	err := d.DecodeElement(&s, &start)
	if err != nil {
		return err
	}

	x, err := time.Parse(DateTimeFormat, s)
	if err != nil {
		return nil
	}
	*t = DateTime{x}

	return nil
}

func (t DateTime) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	if t.Time.Year() < 1880 {
		return xml.Attr{}, nil
	}

	return xml.Attr{Name: name, Value: t.Time.Format(DateTimeFormat)}, nil
}

func (t *DateTime) UnmarshalXMLAttr(attr xml.Attr) error {

	x, err := time.Parse(DateTimeFormat, attr.Value)
	if err != nil {
		return nil
	}
	*t = DateTime{x}

	return nil
}
`

func (s Schema) Render(w io.Writer, tmpl string) error {
	t, err := template.New("base").Funcs(
		template.FuncMap{
			"bt": func() string { return "`" },
		},
	).Parse(tmpl)
	if err != nil {
		return err
	}
	if err := t.Execute(w, s); err != nil {
		return err
	}
	return nil
}
