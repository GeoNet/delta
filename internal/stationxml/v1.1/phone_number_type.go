package stationxml

/*
 *  WARNING: CODE GENERATED AUTOMATICALLY.
 *
 *  Use "go generate" to update these files.
 *
 *  THIS FILE SHOULD NOT BE EDITED BY HAND.
 *
 */

type PhoneNumberType struct {
	Description string `xml:"description,attr,omitempty"`

	CountryCode int `xml:"CountryCode,omitempty"`

	AreaCode int `xml:"AreaCode"`

	PhoneNumber PhoneNumber `xml:"PhoneNumber"`
}
