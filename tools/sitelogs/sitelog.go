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

	FormInformation        FormInformation
	SiteIdentification     SiteIdentification
	SiteLocation           SiteLocation
	GnssReceivers          []GnssReceiver
	GnssAntennas           []GnssAntenna
	GnssHumiditySensors    []GnssMetSensor `xml:"humiditySensor,omitempty"`
	GnssPressureSensors    []GnssMetSensor `xml:"pressureSensor,omitempty"`
	GnssTemperatureSensors []GnssMetSensor `xml:"temperatureSensor,omitempty"`
	ContactAgency          Agency          `xml:"contactAgency"`
	ResponsibleAgency      *Agency         `xml:"responsibleAgency,omitempty"`
	MoreInformation        MoreInformation
}

type SiteLogInput struct {
	XMLName xml.Name `xml:"igsSiteLog"`

	FormInformation        FormInformationInput    `xml:"formInformation"`
	SiteIdentification     SiteIdentificationInput `xml:"siteIdentification"`
	SiteLocation           SiteLocationInput       `xml:"siteLocation"`
	GnssReceivers          []GnssReceiverInput     `xml:"gnssReceiver"`
	GnssAntennas           []GnssAntennaInput      `xml:"gnssAntenna"`
	GnssHumiditySensors    []GnssMetSensor         `xml:"gnssHumiditySensor"`
	GnssTemperatureSensors []GnssMetSensor         `xml:"gnssTemperatureSensor"`
	GnssPressureSensors    []GnssMetSensor         `xml:"gnssPressureSensor"`
	ContactAgency          AgencyInput             `xml:"contactAgency"`
	ResponsibleAgency      *AgencyInput            `xml:"responsibleAgency,omitempty"`
	MoreInformation        MoreInformationInput    `xml:"moreInformation"`
}

func (s SiteLogInput) SiteLog() *SiteLog {
	return &SiteLog{
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
		GnssHumiditySensors:    s.GnssHumiditySensors,
		GnssTemperatureSensors: s.GnssTemperatureSensors,
		GnssPressureSensors:    s.GnssPressureSensors,
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
