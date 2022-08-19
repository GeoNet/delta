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
	AnalogRadiansSecondCfTransferFunction CfTransferFunctionType = 1 + iota
	AnalogHertzCfTransferFunction
	DigitalCfTransferFunction
)

var ErrInvalidCfTransferFunctionType = errors.New("unknown or invalid CfTransferFunctionType value")

type CfTransferFunctionType uint32

func ToCfTransferFunctionType(s string) CfTransferFunctionType {
	switch s {
	case "ANALOG (RADIANS/SECOND)":
		return AnalogRadiansSecondCfTransferFunction
	case "ANALOG (HERTZ)":
		return AnalogHertzCfTransferFunction
	case "DIGITAL":
		return DigitalCfTransferFunction
	default:
		return CfTransferFunctionType(0)
	}
}

func (v CfTransferFunctionType) String() string {
	switch v {
	case AnalogRadiansSecondCfTransferFunction:
		return "ANALOG (RADIANS/SECOND)"
	case AnalogHertzCfTransferFunction:
		return "ANALOG (HERTZ)"
	case DigitalCfTransferFunction:
		return "DIGITAL"
	default:
		return ""
	}
}

func (v CfTransferFunctionType) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	return xml.Attr{Name: name, Value: v.String()}, nil
}

func (v *CfTransferFunctionType) UnmarshalXMLAttr(attr xml.Attr) error {
	switch attr.Value {
	case "ANALOG (RADIANS/SECOND)":
		*v = AnalogRadiansSecondCfTransferFunction
	case "ANALOG (HERTZ)":
		*v = AnalogHertzCfTransferFunction
	case "DIGITAL":
		*v = DigitalCfTransferFunction
	default:
		return ErrInvalidApproximationType
	}

	return nil
}

func (v CfTransferFunctionType) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(v.String(), start)
}

func (v *CfTransferFunctionType) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {

	var s string
	err := d.DecodeElement(&s, &start)
	if err != nil {
		return err
	}

	switch s {
	case "ANALOG (RADIANS/SECOND)":
		*v = AnalogRadiansSecondCfTransferFunction
	case "ANALOG (HERTZ)":
		*v = AnalogHertzCfTransferFunction
	case "DIGITAL":
		*v = DigitalCfTransferFunction
	default:
		return ErrInvalidApproximationType
	}

	return nil
}
