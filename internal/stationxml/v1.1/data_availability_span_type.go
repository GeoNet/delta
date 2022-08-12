package stationxml

/*
 *  WARNING: CODE GENERATED AUTOMATICALLY.
 *
 *  Use "go generate" to update these files.
 *
 *  THIS FILE SHOULD NOT BE EDITED BY HAND.
 *
 */

type DataAvailabilitySpanType struct {
	Start DateTime `xml:"start,attr"`

	End DateTime `xml:"end,attr"`

	NumberSegments int `xml:"numberSegments,attr"`

	MaximumTimeTear float64 `xml:"maximumTimeTear,attr,omitempty"`
}
