package stationxml

/*
 *  WARNING: CODE GENERATED AUTOMATICALLY.
 *
 *  Use "go generate" to update these files.
 *
 *  THIS FILE SHOULD NOT BE EDITED BY HAND.
 *
 */

type DecimationType struct {
	InputSampleRate FrequencyType `xml:"InputSampleRate"`

	Factor int `xml:"Factor"`

	Offset int `xml:"Offset"`

	Delay FloatType `xml:"Delay"`

	Correction FloatType `xml:"Correction"`
}
