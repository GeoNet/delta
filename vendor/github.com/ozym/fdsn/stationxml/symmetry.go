package stationxml

import (
	"encoding/xml"
	"fmt"
)

// The type of data this channel collects. Corresponds to
// channel flags in SEED blockette 52. The SEED volume producer could
// use the first letter of an Output value as the SEED channel flag.
type Symmetry uint

const (
	SymmetryUnknown Symmetry = iota
	SymmetryNone
	SymmetryEven
	SymmetryOdd
	symmetryEnd
)

var symmetryLookup = []string{
	SymmetryUnknown: "UNKNOWN",
	SymmetryNone:    "NONE",
	SymmetryEven:    "EVEN",
	SymmetryOdd:     "ODD",
}

var symmetryMap = map[string]Symmetry{
	"UNKNOWN": SymmetryUnknown,
	"NONE":    SymmetryNone,
	"EVEN":    SymmetryEven,
	"ODD":     SymmetryOdd,
}

func (s Symmetry) IsValid() error {

	if !(s < symmetryEnd) {
		return fmt.Errorf("invalid symmetry entry: %d", s)
	}

	return nil
}

func (s Symmetry) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if !(s < symmetryEnd) {
		return fmt.Errorf("invalid symmetry entry: %d", s)
	}
	return e.EncodeElement(symmetryLookup[s], start)
}

func (s *Symmetry) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {

	var t string
	err := d.DecodeElement(&t, &start)
	if err != nil {
		return err
	}

	if _, ok := symmetryMap[t]; !ok {
		return fmt.Errorf("invalid symmetry: %s", t)
	}

	*s = symmetryMap[t]

	return nil
}
