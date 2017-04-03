package stationxml

import (
	"encoding/xml"
	"fmt"
)

// The type of data this channel collects. Corresponds to
// channel flags in SEED blockette 52. The SEED volume producer could
// use the first letter of an Output value as the SEED channel flag.
type PzTransferFunctionType uint

const (
	PZFunctionUnknown PzTransferFunctionType = iota
	PZFunctionLaplaceRadiansPerSecond
	PZFunctionLaplaceHertz
	PZFunctionLaplaceZTransform
	pzFunctionEnd
)

var pzFunctionLookup = []string{
	PZFunctionUnknown:                 "UNKNOWN",
	PZFunctionLaplaceRadiansPerSecond: "LAPLACE (RADIANS/SECOND)",
	PZFunctionLaplaceHertz:            "LAPLACE (HERTZ)",
	PZFunctionLaplaceZTransform:       "DIGITAL (Z-TRANSFORM)",
}

var pzFunctionMap = map[string]PzTransferFunctionType{
	"UNKNOWN":                  PZFunctionUnknown,
	"LAPLACE (RADIANS/SECOND)": PZFunctionLaplaceRadiansPerSecond,
	"LAPLACE (HERTZ)":          PZFunctionLaplaceHertz,
	"DIGITAL (Z-TRANSFORM)":    PZFunctionLaplaceZTransform,
}

func (f PzTransferFunctionType) IsValid() error {

	if !(f < pzFunctionEnd) {
		return fmt.Errorf("invalid transfer function type: %d", f)
	}

	return nil
}

func (f PzTransferFunctionType) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if !(f < pzFunctionEnd) {
		return fmt.Errorf("invalid function entry: %d", f)
	}
	return e.EncodeElement(pzFunctionLookup[f], start)
}

func (f *PzTransferFunctionType) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {

	var s string
	err := d.DecodeElement(&s, &start)
	if err != nil {
		return err
	}

	if _, ok := pzFunctionMap[s]; !ok {
		return fmt.Errorf("invalid function: %s", s)
	}

	*f = pzFunctionMap[s]

	return nil
}
