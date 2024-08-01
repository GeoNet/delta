package stationxml

/*
 *  WARNING: CODE GENERATED AUTOMATICALLY.
 *
 *  Use "go generate" to update these files.
 *
 *  THIS FILE SHOULD NOT BE EDITED BY HAND.
 *
 */

type PoleZeroType struct {
	Number int `xml:"number,attr,omitempty"`

	Real FloatNoUnitType `xml:"Real"`

	Imaginary FloatNoUnitType `xml:"Imaginary"`
}
