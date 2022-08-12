package stationxml

/*
 *  WARNING: CODE GENERATED AUTOMATICALLY.
 *
 *  Use "go generate" to update these files.
 *
 *  THIS FILE SHOULD NOT BE EDITED BY HAND.
 *
 */

type NetworkType struct {
	BaseNodeType

	Operator []OperatorType `xml:"Operator,omitempty"`

	TotalNumberStations *CounterType `xml:"TotalNumberStations,omitempty"`

	SelectedNumberStations *CounterType `xml:"SelectedNumberStations,omitempty"`

	Station []StationType `xml:"Station,omitempty"`
}
