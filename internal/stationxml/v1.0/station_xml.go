package stationxml

/*
 *  WARNING: CODE GENERATED AUTOMATICALLY.
 *
 *  Use "go generate" to update these files.
 *
 *  THIS FILE SHOULD NOT BE EDITED BY HAND.
 *
 */

import (
	"encoding/xml"
)

// for use when building the root element
const SchemaVersion = 1.0

type FDSNStationXML struct {
	XMLName xml.Name `xml:"http://www.fdsn.org/xml/station/1 FDSNStationXML"`

	RootType
}

func (x FDSNStationXML) Marshal() ([]byte, error) {
	s, err := xml.Marshal(x)
	if err != nil {
		return nil, err
	}
	return append([]byte(xml.Header), append(s, '\n')...), nil
}

func (x FDSNStationXML) MarshalIndent(prefix, indent string) ([]byte, error) {
	s, err := xml.MarshalIndent(x, prefix, indent)
	if err != nil {
		return nil, err
	}
	return append([]byte(xml.Header), append(s, '\n')...), nil
}
