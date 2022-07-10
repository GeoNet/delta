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
	"errors"
)

const (
	MaclaurinApproximation ApproximationType = 1 + iota
)

var ErrInvalidApproximationType = errors.New("unknown or invalid ApproximationType value")

type ApproximationType uint32

func ToApproximationType(s string) ApproximationType {
	switch s {
	case "MACLAURIN":
		return MaclaurinApproximation
	default:
		return ApproximationType(0)
	}
}

func (v ApproximationType) String() string {
	switch v {
	case MaclaurinApproximation:
		return "MACLAURIN"
	default:
		return ""
	}
}

func (v ApproximationType) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	return xml.Attr{Name: name, Value: v.String()}, nil
}

func (v *ApproximationType) UnmarshalXMLAttr(attr xml.Attr) error {
	switch attr.Value {
	case "MACLAURIN":
		*v = MaclaurinApproximation
	default:
		return ErrInvalidApproximationType
	}

	return nil
}

func (v ApproximationType) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(v.String(), start)
}

func (v *ApproximationType) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {

	var s string
	err := d.DecodeElement(&s, &start)
	if err != nil {
		return err
	}

	switch s {
	case "MACLAURIN":
		*v = MaclaurinApproximation
	default:
		return ErrInvalidApproximationType
	}

	return nil
}
