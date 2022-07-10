package stationxml

/*
 *  WARNING: CODE GENERATED AUTOMATICALLY.
 *
 *  Use "go generate" to update these files.
 *
 *  THIS FILE SHOULD NOT BE EDITED BY HAND.
 *
 */

type PolynomialType struct {
	BaseFilterType

	ApproximationType *ApproximationType `xml:"ApproximationType,omitempty"`

	FrequencyLowerBound FrequencyType `xml:"FrequencyLowerBound"`

	FrequencyUpperBound FrequencyType `xml:"FrequencyUpperBound"`

	ApproximationLowerBound float64 `xml:"ApproximationLowerBound"`

	ApproximationUpperBound float64 `xml:"ApproximationUpperBound"`

	MaximumError float64 `xml:"MaximumError"`

	Coefficient []Coefficient `xml:"Coefficient"`
}
