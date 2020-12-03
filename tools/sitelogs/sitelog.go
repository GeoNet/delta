package main

import (
	"bytes"
	"encoding/xml"
	"strings"
)

var equipNameSpace = "http://sopac.ucsd.edu/ns/geodesy/doc/igsSiteLog/equipment/2004"
var contactNameSpace = "http://sopac.ucsd.edu/ns/geodesy/doc/igsSiteLog/contact/2004"
var miNameSpace = "http://sopac.ucsd.edu/ns/geodesy/doc/igsSiteLog/monumentInfo/2004"
var liNameSpace = "http://sopac.ucsd.edu/ns/geodesy/doc/igsSiteLog/localInterferences/2004"
var xmlNameSpace = "http://sopac.ucsd.edu/ns/geodesy/doc/igsSiteLog/2011"
var xsiNameSpace = "http://www.w3.org/2001/XMLSchema-instance"
var schemaLocation = "http://sopac.ucsd.edu/ns/geodesy/doc/igsSiteLog/2011 http://sopac.ucsd.edu/ns/geodesy/doc/igsSiteLog/2011/igsSiteLog.xsd"

// Top-level type for SiteLlog XML.
type SiteLog struct {
	XMLName xml.Name `xml:"igsSiteLog"`

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
