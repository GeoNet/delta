package stationxml

/*
 *  WARNING: CODE GENERATED AUTOMATICALLY.
 *
 *  Use "go generate" to update these files.
 *
 *  THIS FILE SHOULD NOT BE EDITED BY HAND.
 *
 */

type ResponseStageType struct {
	Number CounterType `xml:"number,attr"`

	ResourceId string `xml:"resourceId,attr,omitempty"`

	PolesZeros *PolesZerosType `xml:"PolesZeros,omitempty"`

	Coefficients *CoefficientsType `xml:"Coefficients,omitempty"`

	ResponseList *ResponseListType `xml:"ResponseList,omitempty"`

	FIR *FIRType `xml:"FIR,omitempty"`

	Decimation *DecimationType `xml:"Decimation,omitempty"`

	StageGain *GainType `xml:"StageGain,omitempty"`

	Polynomial *PolynomialType `xml:"Polynomial,omitempty"`
}
