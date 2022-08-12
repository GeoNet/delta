package stationxml

/*
 *  WARNING: CODE GENERATED AUTOMATICALLY.
 *
 *  Use "go generate" to update these files.
 *
 *  THIS FILE SHOULD NOT BE EDITED BY HAND.
 *
 */

import (
	"encoding/xml"
	"time"
)

// only include down to seconds
const DateTimeFormat = "2006-01-02T15:04:05"

type DateTime struct {
	time.Time `xml:",chardata"`
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
