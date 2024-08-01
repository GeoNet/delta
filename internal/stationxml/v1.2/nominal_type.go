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
	NominalNominal NominalType = 1 + iota
	CalculatedNominal
)

var ErrInvalidNominalType = errors.New("unknown or invalid NominalType value")

type NominalType uint32

func ToNominalType(s string) NominalType {
	switch s {
	case "NOMINAL":
		return NominalNominal
	case "CALCULATED":
		return CalculatedNominal
	default:
		return NominalType(0)
	}
}

func (v NominalType) String() string {
	switch v {
	case NominalNominal:
		return "NOMINAL"
	case CalculatedNominal:
		return "CALCULATED"
	default:
		return ""
	}
}

func (v NominalType) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	return xml.Attr{Name: name, Value: v.String()}, nil
}

func (v *NominalType) UnmarshalXMLAttr(attr xml.Attr) error {
	switch attr.Value {
	case "NOMINAL":
		*v = NominalNominal
	case "CALCULATED":
		*v = CalculatedNominal
	default:
		return ErrInvalidApproximationType
	}

	return nil
}

func (v NominalType) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(v.String(), start)
}

func (v *NominalType) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {

	var s string
	err := d.DecodeElement(&s, &start)
	if err != nil {
		return err
	}

	switch s {
	case "NOMINAL":
		*v = NominalNominal
	case "CALCULATED":
		*v = CalculatedNominal
	default:
		return ErrInvalidApproximationType
	}

	return nil
}
