package stationxml

/*
 *  WARNING: CODE GENERATED AUTOMATICALLY.
 *
 *  Use "go generate" to update these files.
 *
 *  THIS FILE SHOULD NOT BE EDITED BY HAND.
 *
 */

type DataAvailabilityType struct {
	Extent *DataAvailabilityExtentType `xml:"Extent,omitempty"`

	Span []DataAvailabilitySpanType `xml:"Span,omitempty"`
}
