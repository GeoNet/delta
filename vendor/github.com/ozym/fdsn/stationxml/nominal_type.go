package stationxml

import (
	"encoding/xml"
	"fmt"
)

type NominalType bool

const (
	Nominal    NominalType = true
	Calculated NominalType = false
)

func (n NominalType) String() string {
	if n == Nominal {
		return "NOMINAL"
	} else {
		return "CALCULATED"
	}
}

func (n NominalType) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if n == Nominal {
		return e.EncodeElement("NOMINAL", start)
	} else {
		return e.EncodeElement("CALCULATED", start)
	}
}

func (n *NominalType) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {

	var s string
	err := d.DecodeElement(&s, &start)
	if err != nil {
		return err
	}

	switch s {
	case "NOMINAL":
		*n = Nominal
	case "CALCULATED":
		*n = Calculated
	default:
		return fmt.Errorf("invalid nominal type: %s", s)
	}

	return nil
}
