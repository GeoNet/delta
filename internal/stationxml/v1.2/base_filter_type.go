package stationxml

/*
 *  WARNING: CODE GENERATED AUTOMATICALLY.
 *
 *  Use "go generate" to update these files.
 *
 *  THIS FILE SHOULD NOT BE EDITED BY HAND.
 *
 */

type BaseFilterType struct {
	ResourceId string `xml:"resourceId,attr,omitempty"`

	Name string `xml:"name,attr,omitempty"`

	Description string `xml:"Description,omitempty"`

	InputUnits UnitsType `xml:"InputUnits"`

	OutputUnits UnitsType `xml:"OutputUnits"`
}
