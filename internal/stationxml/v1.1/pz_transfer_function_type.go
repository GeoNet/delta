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
	LaplaceRadiansSecondPzTransferFunction PzTransferFunctionType = 1 + iota
	LaplaceHertzPzTransferFunction
	DigitalZTransformPzTransferFunction
)

var ErrInvalidPzTransferFunctionType = errors.New("unknown or invalid PzTransferFunctionType value")

type PzTransferFunctionType uint32

func ToPzTransferFunctionType(s string) PzTransferFunctionType {
	switch s {
	case "LAPLACE (RADIANS/SECOND)":
		return LaplaceRadiansSecondPzTransferFunction
	case "LAPLACE (HERTZ)":
		return LaplaceHertzPzTransferFunction
	case "DIGITAL (Z-TRANSFORM)":
		return DigitalZTransformPzTransferFunction
	default:
		return PzTransferFunctionType(0)
	}
}

func (v PzTransferFunctionType) String() string {
	switch v {
	case LaplaceRadiansSecondPzTransferFunction:
		return "LAPLACE (RADIANS/SECOND)"
	case LaplaceHertzPzTransferFunction:
		return "LAPLACE (HERTZ)"
	case DigitalZTransformPzTransferFunction:
		return "DIGITAL (Z-TRANSFORM)"
	default:
		return ""
	}
}

func (v PzTransferFunctionType) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	return xml.Attr{Name: name, Value: v.String()}, nil
}

func (v *PzTransferFunctionType) UnmarshalXMLAttr(attr xml.Attr) error {
	switch attr.Value {
	case "LAPLACE (RADIANS/SECOND)":
		*v = LaplaceRadiansSecondPzTransferFunction
	case "LAPLACE (HERTZ)":
		*v = LaplaceHertzPzTransferFunction
	case "DIGITAL (Z-TRANSFORM)":
		*v = DigitalZTransformPzTransferFunction
	default:
		return ErrInvalidApproximationType
	}

	return nil
}

func (v PzTransferFunctionType) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(v.String(), start)
}

func (v *PzTransferFunctionType) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {

	var s string
	err := d.DecodeElement(&s, &start)
	if err != nil {
		return err
	}

	switch s {
	case "LAPLACE (RADIANS/SECOND)":
		*v = LaplaceRadiansSecondPzTransferFunction
	case "LAPLACE (HERTZ)":
		*v = LaplaceHertzPzTransferFunction
	case "DIGITAL (Z-TRANSFORM)":
		*v = DigitalZTransformPzTransferFunction
	default:
		return ErrInvalidApproximationType
	}

	return nil
}
