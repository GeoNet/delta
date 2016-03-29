package stationxml

import (
	"encoding/xml"
	"fmt"
)

type NameSpace string
type SchemaVersion string

var FDSNNameSpace NameSpace = "http://www.fdsn.org/xml/station/1"
var FDSNSchemaVersion SchemaVersion = "1.0"

// FDSNStationXML represents the FDSN StationXML schema's root type.
//
// Designed as an XML representation of SEED metadata, the schema maps to
// the most important and commonly used structures of SEED 2.4. When definitions and usage are
// underdefined the SEED manual should be referred to for clarification.
//
// Top-level type for Station XML. Required field are Source (network ID of the institution sending
// the message) and one or more Network containers or one or more Station containers.
type FDSNStationXML struct {
	NameSpace NameSpace `xml:"xmlns,attr"`

	// The schema version compatible with the document.
	SchemaVersion SchemaVersion `xml:"schemaVersion,attr"`

	// Network ID of the institution sending the message.
	Source string `xml:"Source"`

	// Name of the institution sending this message.
	Sender string `xml:",omitempty" json:",omitempty"`

	//Name of the software module that generated this document.
	Module string `xml:",omitempty" json:",omitempty"`

	// This is the address of the query that generated the document, or,
	// if applicable, the address of the software that generated this document.
	ModuleURI AnyURI `xml:",omitempty" json:",omitempty"`

	Created DateTime `xml:"Created"`

	Networks []Network `xml:"Network,omitempty" json:",omitempty"`
}

func NewFDSNStationXML(source, sender, module string, uri AnyURI, networks []Network) FDSNStationXML {
	return FDSNStationXML{
		NameSpace:     FDSNNameSpace,
		SchemaVersion: FDSNSchemaVersion,
		Source:        source,
		Sender:        sender,
		Module:        module,
		ModuleURI:     uri,
		Networks:      networks,
		Created:       Now(),
	}
}

func (x FDSNStationXML) IsValid() error {

	if x.NameSpace != FDSNNameSpace {
		return fmt.Errorf("wrong name space: %s", x.NameSpace)
	}
	if x.SchemaVersion != FDSNSchemaVersion {
		return fmt.Errorf("wrong schema version: %s", x.SchemaVersion)
	}

	if !(len(x.Source) > 0) {
		return fmt.Errorf("empty source element")
	}

	if x.Created.IsZero() {
		return fmt.Errorf("created date should not be zero")
	}

	if err := Validate(x.Created); err != nil {
		return err
	}

	for _, n := range x.Networks {
		if err := Validate(n); err != nil {
			return err
		}
	}

	return nil
}

func (x FDSNStationXML) Marshal() ([]byte, error) {
	h := xml.Header // []byte(FDSN_XML_HEADER)
	s, err := xml.Marshal(x)
	if err != nil {
		return nil, err
	}
	return append([]byte(h), append(s, '\n')...), nil
}
