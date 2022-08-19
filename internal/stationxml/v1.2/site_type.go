package stationxml

/*
 *  WARNING: CODE GENERATED AUTOMATICALLY.
 *
 *  Use "go generate" to update these files.
 *
 *  THIS FILE SHOULD NOT BE EDITED BY HAND.
 *
 */

type SiteType struct {
	Name string `xml:"Name"`

	Description string `xml:"Description,omitempty"`

	Town string `xml:"Town,omitempty"`

	County string `xml:"County,omitempty"`

	Region string `xml:"Region,omitempty"`

	Country string `xml:"Country,omitempty"`
}
