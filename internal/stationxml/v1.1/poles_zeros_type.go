package stationxml

/*
 *  WARNING: CODE GENERATED AUTOMATICALLY.
 *
 *  Use "go generate" to update these files.
 *
 *  THIS FILE SHOULD NOT BE EDITED BY HAND.
 *
 */

type PolesZerosType struct {
	BaseFilterType

	PzTransferFunctionType PzTransferFunctionType `xml:"PzTransferFunctionType"`

	NormalizationFactor float64 `xml:"NormalizationFactor,omitempty"`

	NormalizationFrequency FrequencyType `xml:"NormalizationFrequency"`

	Zero []PoleZeroType `xml:"Zero,omitempty"`

	Pole []PoleZeroType `xml:"Pole,omitempty"`
}
