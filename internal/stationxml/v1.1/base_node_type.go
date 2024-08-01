package stationxml

/*
 *  WARNING: CODE GENERATED AUTOMATICALLY.
 *
 *  Use "go generate" to update these files.
 *
 *  THIS FILE SHOULD NOT BE EDITED BY HAND.
 *
 */

type BaseNodeType struct {
	Code string `xml:"code,attr"`

	StartDate DateTime `xml:"startDate,attr,omitempty"`

	EndDate DateTime `xml:"endDate,attr,omitempty"`

	RestrictedStatus RestrictedStatusType `xml:"restrictedStatus,attr,omitempty"`

	AlternateCode string `xml:"alternateCode,attr,omitempty"`

	HistoricalCode string `xml:"historicalCode,attr,omitempty"`

	Description string `xml:"Description,omitempty"`

	Identifier []IdentifierType `xml:"Identifier,omitempty"`

	Comment []CommentType `xml:"Comment,omitempty"`

	DataAvailability *DataAvailabilityType `xml:"DataAvailability,omitempty"`
}
