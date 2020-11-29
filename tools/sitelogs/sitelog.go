package main

import (
	"bytes"
	"encoding/xml"
	"strings"
	//	"fmt"
)

var equipNameSpace = "http://sopac.ucsd.edu/ns/geodesy/doc/igsSiteLog/equipment/2004"
var contactNameSpace = "http://sopac.ucsd.edu/ns/geodesy/doc/igsSiteLog/contact/2004"
var miNameSpace = "http://sopac.ucsd.edu/ns/geodesy/doc/igsSiteLog/monumentInfo/2004"
var liNameSpace = "http://sopac.ucsd.edu/ns/geodesy/doc/igsSiteLog/localInterferences/2004"
var xmlNameSpace = "http://sopac.ucsd.edu/ns/geodesy/doc/igsSiteLog/2011"
var xsiNameSpace = "http://www.w3.org/2001/XMLSchema-instance"
var schemaLocation = "http://sopac.ucsd.edu/ns/geodesy/doc/igsSiteLog/2011 http://sopac.ucsd.edu/ns/geodesy/doc/igsSiteLog/2011/igsSiteLog.xsd"

/*
xmlns:equip="http://sopac.ucsd.edu/ns/geodesy/doc/igsSiteLog/equipment/2004"
xmlns:contact="http://sopac.ucsd.edu/ns/geodesy/doc/igsSiteLog/contact/2004"
xmlns:mi="http://sopac.ucsd.edu/ns/geodesy/doc/igsSiteLog/monumentInfo/2004"
xmlns:li="http://sopac.ucsd.edu/ns/geodesy/doc/igsSiteLog/localInterferences/2004"
xmlns="http://sopac.ucsd.edu/ns/geodesy/doc/igsSiteLog/2011"
xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
xsi:schemaLocation="http://sopac.ucsd.edu/ns/geodesy/doc/igsSiteLog/2011 http://sopac.ucsd.edu/ns/geodesy/doc/igsSiteLog/2011/igsSiteLog.xsd"
*/

/*
<geo:GeodesyML
gml:id="GEO_1"
xmlns:gml="http://www.opengis.net/gml/3.2"
xmlns:geo="urn:xml-gov-au:icsm:egeodesy:0.2"
xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
xmlns:gmd="http://www.isotc211.org/2005/gmd"
xmlns:gco="http://www.isotc211.org/2005/gco"
xmlns:xlink="http://www.w3.org/1999/xlink"
*/

type SrsName struct {
	SrsName string `xml:"urn:xml-gov-au:icsm:egeodesy:0.3 srsName,attr"`
}

// FDSNStationXML represents the FDSN StationXML schema's root type.
//
// Designed as an XML representation of SEED metadata, the schema maps to
// the most important and commonly used structures of SEED 2.4. When definitions and usage are
// underdefined the SEED manual should be referred to for clarification.
//
// Top-level type for Station XML. Required field are Source (network ID of the institution sending
// the message) and one or more Network containers or one or more Station containers.
type SiteLog struct {
	XMLName xml.Name `xml:"igsSiteLog"`

	//<igsSiteLog xmlns:equip="http://sopac.ucsd.edu/ns/geodesy/doc/igsSiteLog/equipment/2004" xmlns:contact="http://sopac.ucsd.edu/ns/geodesy/doc/igsSiteLog/contact/2004" xmlns:mi="http://sopac.ucsd.edu/ns/geodesy/doc/igsSiteLog/monumentInfo/2004" xmlns:li="http://sopac.ucsd.edu/ns/geodesy/doc/igsSiteLog/localInterferences/2004" xmlns="http://sopac.ucsd.edu/ns/geodesy/doc/igsSiteLog/2011" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://sopac.ucsd.edu/ns/geodesy/doc/igsSiteLog/2011 http://sopac.ucsd.edu/ns/geodesy/doc/igsSiteLog/2011/igsSiteLog.xsd">

	EquipNameSpace   string `xml:"xmlns:equip,attr"`
	ContactNameSpace string `xml:"xmlns:contact,attr"`
	MiNameSpace      string `xml:"xmlns:mi,attr"`
	LiNameSpace      string `xml:"xmlns:li,attr"`
	XmlNameSpace     string `xml:"xmlns,attr"`
	XsiNameSpace     string `xml:"xmlns:xsi,attr"`
	SchemaLocation   string `xml:"xsi:schemaLocation,attr"`

	FormInformation    FormInformation
	SiteIdentification SiteIdentification
	SiteLocation       SiteLocation
	GnssReceivers      []GnssReceiver
	GnssAntennas       []GnssAntenna
	GnssMetSensors     []GnssMetSensor
	ContactAgency      Agency  `xml:"contactAgency"`
	ResponsibleAgency  *Agency `xml:"responsibleAgency,omitempty"`
	MoreInformation    MoreInformation
}

type SiteLogInput struct {
	XMLName xml.Name `xml:"igsSiteLog"`

	FormInformation    FormInformationInput    `xml:"formInformation"`
	SiteIdentification SiteIdentificationInput `xml:"siteIdentification"`
	SiteLocation       SiteLocationInput       `xml:"siteLocation"`
	GnssReceivers      []GnssReceiverInput     `xml:"gnssReceiver"`
	GnssAntennas       []GnssAntennaInput      `xml:"gnssAntenna"`
	GnssMetSensors     []GnssMetSensor         `xml:"gnssMetSensor"`
	ContactAgency      AgencyInput             `xml:"contactAgency"`
	ResponsibleAgency  *AgencyInput            `xml:"responsibleAgency,omitempty"`
	MoreInformation    MoreInformationInput    `xml:"moreInformation"`
}

func (s SiteLogInput) SiteLog() SiteLog {
	return SiteLog{
		EquipNameSpace:   equipNameSpace,
		ContactNameSpace: contactNameSpace,
		MiNameSpace:      miNameSpace,
		LiNameSpace:      liNameSpace,
		XmlNameSpace:     xmlNameSpace,
		XsiNameSpace:     xsiNameSpace,
		SchemaLocation:   schemaLocation,

		FormInformation:    FormInformation(s.FormInformation),
		SiteIdentification: SiteIdentification(s.SiteIdentification),
		SiteLocation: SiteLocation{
			City:                    s.SiteLocation.City,
			State:                   s.SiteLocation.State,
			Country:                 s.SiteLocation.Country,
			TectonicPlate:           s.SiteLocation.TectonicPlate,
			ApproximatePositionITRF: ApproximatePositionITRF(s.SiteLocation.ApproximatePositionITRF),
			Notes:                   s.SiteLocation.Notes,
		},
		GnssReceivers: func() []GnssReceiver {
			var receivers []GnssReceiver
			for _, r := range s.GnssReceivers {
				receivers = append(receivers, GnssReceiver(r))
			}
			return receivers
		}(),
		GnssAntennas: func() []GnssAntenna {
			var antennas []GnssAntenna
			for _, r := range s.GnssAntennas {
				antennas = append(antennas, GnssAntenna(r))
			}
			return antennas
		}(),
		GnssMetSensors: s.GnssMetSensors,
		ContactAgency: Agency{
			Agency:                s.ContactAgency.Agency,
			PreferredAbbreviation: s.ContactAgency.PreferredAbbreviation,
			MailingAddress:        s.ContactAgency.MailingAddress,
			PrimaryContact:        Contact(s.ContactAgency.PrimaryContact),
			SecondaryContact:      Contact(s.ContactAgency.SecondaryContact),
			Notes:                 s.ContactAgency.Notes,
		},
		ResponsibleAgency: func() *Agency {
			if s.ResponsibleAgency != nil {
				return &Agency{
					Agency:                s.ResponsibleAgency.Agency,
					PreferredAbbreviation: s.ResponsibleAgency.PreferredAbbreviation,
					MailingAddress:        s.ResponsibleAgency.MailingAddress,
					PrimaryContact:        Contact(s.ResponsibleAgency.PrimaryContact),
					SecondaryContact:      Contact(s.ResponsibleAgency.SecondaryContact),
					Notes:                 s.ResponsibleAgency.Notes,
				}
			}
			return nil
		}(),
		MoreInformation: MoreInformation(s.MoreInformation),
	}
}

/*
type Site struct {
	XMLName xml.Name `xml:"geo:Site"`
	Id      string   `xml:"gml:id,attr"`
	Type    Type
}

type Type struct {
	XMLName xml.Name `xml:"geo:type"`
	Type    string   `xml:",chardata"`
}

type Monument struct {
	XMLName xml.Name `xml:"geo:Monument"`
	Link    string   `xml:"xlink:href,attr"`
	Names   []Name   `xml:",omitempty"`
}

type Name struct {
	XMLName   xml.Name `xml:"gml:name"`
	CodeSpace string   `xml:"codeSpace,attr"`
	Name      string   `xml:",chardata"`
}

type Receiver struct {
	XMLName xml.Name `xml:"geo:gnssReceiver"`
	Id      string   `xml:"gml:id,attr"`

	ManufacturerSerialNumber string  `xml:"geo:manufacturerSerialNumber"`
	ReceiverType             string  `xml:"geo:receiverType"`
	SatelliteSystem          string  `xml:"geo:satelliteSystem"`
	SerialNumber             string  `xml:"geo:serialNumber"`
	FirmwareVersion          string  `xml:"geo:firmwareVersion"`
	ElevationCutoffSetting   float32 `xml:"geo:elevationCutoffSetting"`
	DateInstalled            string  `xml:"geo:dateInstalled"`
	DateRemoved              string  `xml:"geo:dateRemoved"`
	TemperatureStabilization string  `xml:"geo:temperatureStabilization"`
	Notes                    string  `xml:"geo:notes"`
}
*/

func NewSiteLog() SiteLog {
	return SiteLog{
		EquipNameSpace:   equipNameSpace,
		ContactNameSpace: contactNameSpace,
		MiNameSpace:      miNameSpace,
		LiNameSpace:      liNameSpace,
		XmlNameSpace:     xmlNameSpace,
		XsiNameSpace:     xsiNameSpace,
		SchemaLocation:   schemaLocation,
	}
}

//func NewGeodesyML(site Site, monument Monument, receivers []Receiver /*source, sender, module string, uri AnyURI, networks []Network*/) GeodesyML {
//	return GeodesyML{
//		GMLNameSpace: GMLNameSpace,
//		//GEONameSpace:   GEONameSpace,
//		XSINameSpace:   XSINameSpace,
//		GMDNameSpace:   GMDNameSpace,
//		GCONameSpace:   GCONameSpace,
//		XLINKNameSpace: XLINKNameSpace,
//		SchemaLocation: SchemaLocation,
//
//		Site:      site,
//		Monument:  monument,
//		Receivers: receivers,

/*
	Site: Site{
		Id: "Site_1",
		Type: Type{
			Type: "CORRS",
		},
		Monument: Monument{
			Link: "#MONUMENT_1",
		},
	},
*/

/*
	NameSpace:     FDSNNameSpace,
	SchemaVersion: FDSNSchemaVersion,
	Source:        source,
	Sender:        sender,
	Module:        module,
	ModuleURI:     uri,
	Networks:      networks,
	Created:       Now(),
*/
//	}
//}

//func (x GeodesyML) IsValid() error {

/*
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
*/

//	return nil
//}

func (x SiteLog) Strip() SiteLog {
	x.ContactAgency.MailingAddress = strings.ReplaceAll(x.ContactAgency.MailingAddress, "\n", " ")
	if x.ResponsibleAgency != nil {
		x.ResponsibleAgency.MailingAddress = strings.ReplaceAll(x.ResponsibleAgency.MailingAddress, "\n", " ")
	}
	x.MoreInformation.Notes = strings.ReplaceAll(x.MoreInformation.Notes, "\n", " ")
	return x
}

func (x SiteLog) Marshal() ([]byte, error) {
	h := xml.Header
	s, err := xml.MarshalIndent(x, "", "  ")
	if err != nil {
		return nil, err
	}
	return append([]byte(h), append(s, '\n')...), nil
}

func (x SiteLog) MarshalLegacy() ([]byte, error) {
	h := xml.Header

	h = strings.ReplaceAll(h, "\"", "'")

	x.ContactAgency.MailingAddress = strings.ReplaceAll(x.ContactAgency.MailingAddress, "\n", " ")
	if x.ResponsibleAgency != nil {
		x.ResponsibleAgency.MailingAddress = strings.ReplaceAll(x.ResponsibleAgency.MailingAddress, "\n", " ")
	}
	x.MoreInformation.Notes = strings.ReplaceAll(x.MoreInformation.Notes, "\n", " ")

	s, err := xml.MarshalIndent(x, "", "  ")
	if err != nil {
		return nil, err
	}
	s = bytes.ReplaceAll(s, []byte("&#39;"), []byte("'"))
	s = bytes.ReplaceAll(s, []byte("&#xA;"), []byte("\n"))
	return append([]byte(h), append(s, '\n')...), nil
}
