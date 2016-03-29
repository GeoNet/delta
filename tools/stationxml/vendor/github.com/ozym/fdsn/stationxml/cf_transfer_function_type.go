package stationxml

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
)

// The type of data this channel collects. Corresponds to
// channel flags in SEED blockette 52. The SEED volume producer could
// use the first letter of an Output value as the SEED channel flag.
type CfTransferFunctionType uint

const (
	CfFunctionUnknown CfTransferFunctionType = iota
	CfFunctionAnalogRadiansPerSecond
	CfFunctionAnalogHertz
	CfFunctionDigital
	cfFunctionEnd
)

var cfFunctionTypeLookup = []string{
	CfFunctionUnknown:                "UNKNOWN",
	CfFunctionAnalogRadiansPerSecond: "ANALOG (RADIANS/SECOND)",
	CfFunctionAnalogHertz:            "ANALOG (HERTZ)",
	CfFunctionDigital:                "DIGITAL",
}

var cfFunctionTypeMap = map[string]CfTransferFunctionType{
	"UNKNOWN":                 CfFunctionUnknown,
	"ANALOG (RADIANS/SECOND)": CfFunctionAnalogRadiansPerSecond,
	"ANALOG (HERTZ)":          CfFunctionAnalogHertz,
	"DIGITAL":                 CfFunctionDigital,
}

func (f CfTransferFunctionType) String() string {

	if f < cfFunctionEnd {
		return cfFunctionTypeLookup[f]
	}

	return ""
}

func (f CfTransferFunctionType) IsValid() error {

	if !(f < cfFunctionEnd) {
		return fmt.Errorf("invalid function entry: %d", f)
	}

	return nil
}

func (f CfTransferFunctionType) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if !(f < cfFunctionEnd) {
		return fmt.Errorf("invalid function entry: %d", f)
	}
	return e.EncodeElement(cfFunctionTypeLookup[f], start)
}

func (f *CfTransferFunctionType) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {

	var s string
	err := d.DecodeElement(&s, &start)
	if err != nil {
		return err
	}

	if _, ok := cfFunctionTypeMap[s]; !ok {
		return fmt.Errorf("invalid function: %s", s)
	}

	*f = cfFunctionTypeMap[s]

	return nil
}

func (f CfTransferFunctionType) MarshalJSON() ([]byte, error) {
	if !(f < cfFunctionEnd) {
		return nil, fmt.Errorf("invalid type: %d", f)
	}
	return json.Marshal(cfFunctionTypeLookup[f])
}

func (f *CfTransferFunctionType) UnmarshalJSON(data []byte) error {
	var b []byte
	err := json.Unmarshal(data, b)
	if err != nil {
		return err
	}
	s := string(b)

	if _, ok := cfFunctionTypeMap[s]; !ok {
		return fmt.Errorf("invalid type: %s", s)
	}

	*f = cfFunctionTypeMap[s]

	return nil
}
