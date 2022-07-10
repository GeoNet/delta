package stationxml

/*
 *  WARNING: CODE GENERATED AUTOMATICALLY.
 *
 *  Use "go generate" to update these files.
 *
 *  THIS FILE SHOULD NOT BE EDITED BY HAND.
 *
 */

type SensitivityType struct {
	GainType

	InputUnits UnitsType `xml:"InputUnits"`

	OutputUnits UnitsType `xml:"OutputUnits"`

	FrequencyStart *float64 `xml:"FrequencyStart"`

	FrequencyEnd *float64 `xml:"FrequencyEnd"`

	FrequencyDBVariation *float64 `xml:"FrequencyDBVariation"`
}
