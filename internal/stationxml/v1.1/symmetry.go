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
	NoneSymmetry Symmetry = 1 + iota
	EvenSymmetry
	OddSymmetry
)

var ErrInvalidSymmetry = errors.New("unknown or invalid Symmetry value")

type Symmetry uint32

func ToSymmetry(s string) Symmetry {
	switch s {
	case "NONE":
		return NoneSymmetry
	case "EVEN":
		return EvenSymmetry
	case "ODD":
		return OddSymmetry
	default:
		return Symmetry(0)
	}
}

func (v Symmetry) String() string {
	switch v {
	case NoneSymmetry:
		return "NONE"
	case EvenSymmetry:
		return "EVEN"
	case OddSymmetry:
		return "ODD"
	default:
		return ""
	}
}

func (v Symmetry) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	return xml.Attr{Name: name, Value: v.String()}, nil
}

func (v *Symmetry) UnmarshalXMLAttr(attr xml.Attr) error {
	switch attr.Value {
	case "NONE":
		*v = NoneSymmetry
	case "EVEN":
		*v = EvenSymmetry
	case "ODD":
		*v = OddSymmetry
	default:
		return ErrInvalidApproximationType
	}

	return nil
}

func (v Symmetry) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(v.String(), start)
}

func (v *Symmetry) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {

	var s string
	err := d.DecodeElement(&s, &start)
	if err != nil {
		return err
	}

	switch s {
	case "NONE":
		*v = NoneSymmetry
	case "EVEN":
		*v = EvenSymmetry
	case "ODD":
		*v = OddSymmetry
	default:
		return ErrInvalidApproximationType
	}

	return nil
}
