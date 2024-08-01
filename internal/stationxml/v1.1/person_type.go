package stationxml

/*
 *  WARNING: CODE GENERATED AUTOMATICALLY.
 *
 *  Use "go generate" to update these files.
 *
 *  THIS FILE SHOULD NOT BE EDITED BY HAND.
 *
 */

type PersonType struct {
	Name []string `xml:"Name,omitempty"`

	Agency []string `xml:"Agency,omitempty"`

	Email []EmailType `xml:"Email,omitempty"`

	Phone []PhoneNumberType `xml:"Phone,omitempty"`
}
