package stationxml

/*
 *  WARNING: CODE GENERATED AUTOMATICALLY.
 *
 *  Use "go generate" to update these files.
 *
 *  THIS FILE SHOULD NOT BE EDITED BY HAND.
 *
 */

type CommentType struct {
	Id CounterType `xml:"id,attr,omitempty"`

	Value string `xml:"Value"`

	BeginEffectiveTime DateTime `xml:"BeginEffectiveTime,omitempty"`

	EndEffectiveTime DateTime `xml:"EndEffectiveTime,omitempty"`

	Author []PersonType `xml:"Author,omitempty"`
}
