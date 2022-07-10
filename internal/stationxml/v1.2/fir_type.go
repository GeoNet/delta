package stationxml

/*
 *  WARNING: CODE GENERATED AUTOMATICALLY.
 *
 *  Use "go generate" to update these files.
 *
 *  THIS FILE SHOULD NOT BE EDITED BY HAND.
 *
 */

type FIRType struct {
	BaseFilterType

	Symmetry Symmetry `xml:"Symmetry"`

	NumeratorCoefficient []NumeratorCoefficient `xml:"NumeratorCoefficient,omitempty"`
}
