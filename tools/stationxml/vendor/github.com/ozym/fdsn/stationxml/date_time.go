package stationxml

import (
	"encoding/xml"
	"fmt"
	"time"
)

type DateTime struct {
	time.Time `xml:",chardata,omitempty"`
}

func Now() DateTime {
	return DateTime{Time: time.Now()}
}

const DateTimeFormat = "2006-01-02T15:04:05"

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
