package stationxml

/*
 *  WARNING: CODE GENERATED AUTOMATICALLY.
 *
 *  Use "go generate" to update these files.
 *
 *  THIS FILE SHOULD NOT BE EDITED BY HAND.
 *
 */

type ResponseType struct {
	ResourceId string `xml:"resourceId,attr,omitempty"`

	InstrumentSensitivity *SensitivityType `xml:"InstrumentSensitivity,omitempty"`

	InstrumentPolynomial *PolynomialType `xml:"InstrumentPolynomial,omitempty"`

	Stage []ResponseStageType `xml:"Stage,omitempty"`
}
