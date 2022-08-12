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
	Triggered Type = 1 + iota
	Continuous
	Health
	Geophysical
	Weather
	Flag
	Synthesized
	Input
	Experimental
	Maintenance
	Beam
)

var ErrInvalidType = errors.New("unknown or invalid Type value")

type Type uint32

func ToType(s string) Type {
	switch s {
	case "TRIGGERED":
		return Triggered
	case "CONTINUOUS":
		return Continuous
	case "HEALTH":
		return Health
	case "GEOPHYSICAL":
		return Geophysical
	case "WEATHER":
		return Weather
	case "FLAG":
		return Flag
	case "SYNTHESIZED":
		return Synthesized
	case "INPUT":
		return Input
	case "EXPERIMENTAL":
		return Experimental
	case "MAINTENANCE":
		return Maintenance
	case "BEAM":
		return Beam
	default:
		return Type(0)
	}
}

func (v Type) String() string {
	switch v {
	case Triggered:
		return "TRIGGERED"
	case Continuous:
		return "CONTINUOUS"
	case Health:
		return "HEALTH"
	case Geophysical:
		return "GEOPHYSICAL"
	case Weather:
		return "WEATHER"
	case Flag:
		return "FLAG"
	case Synthesized:
		return "SYNTHESIZED"
	case Input:
		return "INPUT"
	case Experimental:
		return "EXPERIMENTAL"
	case Maintenance:
		return "MAINTENANCE"
	case Beam:
		return "BEAM"
	default:
		return ""
	}
}

func (v Type) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	return xml.Attr{Name: name, Value: v.String()}, nil
}

func (v *Type) UnmarshalXMLAttr(attr xml.Attr) error {
	switch attr.Value {
	case "TRIGGERED":
		*v = Triggered
	case "CONTINUOUS":
		*v = Continuous
	case "HEALTH":
		*v = Health
	case "GEOPHYSICAL":
		*v = Geophysical
	case "WEATHER":
		*v = Weather
	case "FLAG":
		*v = Flag
	case "SYNTHESIZED":
		*v = Synthesized
	case "INPUT":
		*v = Input
	case "EXPERIMENTAL":
		*v = Experimental
	case "MAINTENANCE":
		*v = Maintenance
	case "BEAM":
		*v = Beam
	default:
		return ErrInvalidApproximationType
	}

	return nil
}

func (v Type) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(v.String(), start)
}

func (v *Type) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {

	var s string
	err := d.DecodeElement(&s, &start)
	if err != nil {
		return err
	}

	switch s {
	case "TRIGGERED":
		*v = Triggered
	case "CONTINUOUS":
		*v = Continuous
	case "HEALTH":
		*v = Health
	case "GEOPHYSICAL":
		*v = Geophysical
	case "WEATHER":
		*v = Weather
	case "FLAG":
		*v = Flag
	case "SYNTHESIZED":
		*v = Synthesized
	case "INPUT":
		*v = Input
	case "EXPERIMENTAL":
		*v = Experimental
	case "MAINTENANCE":
		*v = Maintenance
	case "BEAM":
		*v = Beam
	default:
		return ErrInvalidApproximationType
	}

	return nil
}
