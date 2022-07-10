package stationxml

/*
 *  WARNING: CODE GENERATED AUTOMATICALLY.
 *
 *  Use "go generate" to update these files.
 *
 *  THIS FILE SHOULD NOT BE EDITED BY HAND.
 *
 */

type RootType struct {
	SchemaVersion float64 `xml:"schemaVersion,attr"`

	Source string `xml:"Source"`

	Sender string `xml:"Sender,omitempty"`

	Module string `xml:"Module,omitempty"`

	Created *DateTime `xml:"Created"`

	Network []NetworkType `xml:"Network"`
}
