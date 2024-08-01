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
	OpenRestrictedStatus RestrictedStatusType = 1 + iota
	ClosedRestrictedStatus
	PartialRestrictedStatus
)

var ErrInvalidRestrictedStatusType = errors.New("unknown or invalid RestrictedStatusType value")

type RestrictedStatusType uint32

func ToRestrictedStatusType(s string) RestrictedStatusType {
	switch s {
	case "open":
		return OpenRestrictedStatus
	case "closed":
		return ClosedRestrictedStatus
	case "partial":
		return PartialRestrictedStatus
	default:
		return RestrictedStatusType(0)
	}
}

func (v RestrictedStatusType) String() string {
	switch v {
	case OpenRestrictedStatus:
		return "open"
	case ClosedRestrictedStatus:
		return "closed"
	case PartialRestrictedStatus:
		return "partial"
	default:
		return ""
	}
}

func (v RestrictedStatusType) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	return xml.Attr{Name: name, Value: v.String()}, nil
}

func (v *RestrictedStatusType) UnmarshalXMLAttr(attr xml.Attr) error {
	switch attr.Value {
	case "open":
		*v = OpenRestrictedStatus
	case "closed":
		*v = ClosedRestrictedStatus
	case "partial":
		*v = PartialRestrictedStatus
	default:
		return ErrInvalidApproximationType
	}

	return nil
}

func (v RestrictedStatusType) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(v.String(), start)
}

func (v *RestrictedStatusType) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {

	var s string
	err := d.DecodeElement(&s, &start)
	if err != nil {
		return err
	}

	switch s {
	case "open":
		*v = OpenRestrictedStatus
	case "closed":
		*v = ClosedRestrictedStatus
	case "partial":
		*v = PartialRestrictedStatus
	default:
		return ErrInvalidApproximationType
	}

	return nil
}
