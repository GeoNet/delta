package stationxml

import (
	"encoding/xml"
	"fmt"
)

type ApproximationType uint

const (
	ApproximationTypeUnknown ApproximationType = iota
	ApproximationTypeMaclaurin
	approximationTypeEnd
)

var approximationTypeLookup = []string{
	ApproximationTypeUnknown:   "UNKNOWN",
	ApproximationTypeMaclaurin: "MACLAURIN",
}

var approximationTypeMap = map[string]ApproximationType{
	"UNKNOWN":   ApproximationTypeUnknown,
	"MACLAURIN": ApproximationTypeMaclaurin,
}

// The type of data this channel collects. Corresponds to
// channel flags in SEED blockette 52. The SEED volume producer could
// use the first letter of an Output value as the SEED channel flag.
func (f ApproximationType) String() string {

	if f < approximationTypeEnd {
		return approximationTypeLookup[f]
	}

	return ""
}

func (f ApproximationType) IsValid() error {

	if !(f < approximationTypeEnd) {
		return fmt.Errorf("invalid approximation type: %d", f)
	}

	return nil
}

func (f ApproximationType) MarshalXML(e *xml.Encoder, start xml.StartElement) error {

	if f < approximationTypeEnd {
		return e.EncodeElement(approximationTypeLookup[f], start)
	}

	return fmt.Errorf("invalid approximation type: %d", f)
}

func (f *ApproximationType) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {

	var s string
	err := d.DecodeElement(&s, &start)
	if err != nil {
		return err
	}

	if _, ok := approximationTypeMap[s]; !ok {
		return fmt.Errorf("invalid function: %s", s)
	}

	*f = approximationTypeMap[s]

	return nil
}
