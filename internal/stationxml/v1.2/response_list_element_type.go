package stationxml

/*
 *  WARNING: CODE GENERATED AUTOMATICALLY.
 *
 *  Use "go generate" to update these files.
 *
 *  THIS FILE SHOULD NOT BE EDITED BY HAND.
 *
 */

type ResponseListElementType struct {
	Frequency FrequencyType `xml:"Frequency"`

	Amplitude FloatType `xml:"Amplitude"`

	Phase AngleType `xml:"Phase"`
}
