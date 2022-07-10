package stationxml

/*
 *  WARNING: CODE GENERATED AUTOMATICALLY.
 *
 *  Use "go generate" to update these files.
 *
 *  THIS FILE SHOULD NOT BE EDITED BY HAND.
 *
 */

type FloatType struct {
	Unit string `xml:"unit,attr,omitempty"`

	PlusError float64 `xml:"plusError,attr,omitempty"`

	MinusError float64 `xml:"minusError,attr,omitempty"`

	Value float64 `xml:",chardata"`
}
