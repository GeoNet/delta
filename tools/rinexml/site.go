package main

import (
	"encoding/xml"
)

type SiteXML struct {
	XMLName xml.Name

	Mark     MarkXML     `xml:"mark"`
	Location LocationXML `xml:"location"`

	Sessions []CGPSSessionXML `xml:"cgps-session,omitempty"`
}

type MarkXML struct {
	GeodeticCode string `xml:"geodetic-code"`
	DomesNumber  string `xml:"domes-number"`
}

type LocationXML struct {
	Latitude  float64 `xml:"latitude"`
	Longitude float64 `xml:"longitude"`
	Height    float64 `xml:"height"`
	Datum     string  `xml:"datum"`
}

type Number struct {
	Units string `xml:"unit,attr"`
	Value string `xml:",chardata"`
}

type OperatorXML struct {
	Name   string `xml:"name"`
	Agency string `xml:"agency"`
}

type RinexXML struct {
	HeaderCommentName string `xml:"header-comment-name"`
	HeaderCommentText string `xml:"header-comment-text"`
}

type DownloadNameFormatXML struct {
	Type    string `xml:"type,attr"`
	Year    string `xml:"year,omitempty"`
	YearDay string `xml:"year-day,omitempty"`
	Month   string `xml:"month,omitempty"`
	Day     string `xml:"day,omitempty"`
	Hour    string `xml:"hour,omitempty"`
	Minute  string `xml:"minute,omitempty"`
	Second  string `xml:"second,omitempty"`
}

type FirmwareHistoryXML struct {
	StartTime string `xml:"start-time"`
	StopTime  string `xml:"stop-time"`
	Version   string `xml:"version"`
}

type ReceiverXML struct {
	SerialNumber      string               `xml:"serial-number"`
	IGSDesignation    string               `xml:"igs-designation"`
	FirmwareHistories []FirmwareHistoryXML `xml:"firmware-history"`
}

type CGPSAntennaXML struct {
	SerialNumber   string `xml:"serial-number"`
	IGSDesignation string `xml:"igs-designation"`
}

type InstalledCGPSAntennaXML struct {
	Height      Number         `xml:"height"`
	OffsetEast  Number         `xml:"offset-east"`
	OffsetNorth Number         `xml:"offset-north"`
	Radome      string         `xml:"radome"`
	CGPSAntenna CGPSAntennaXML `xml:"cgps-antenna"`
}

type MetSensor struct {
	Model      string `xml:"model"`
	Type       string `xml:"type"`
	HrAccuracy string `xml:"hr>accuracy"`
	PrAccuracy string `xml:"pr>accuracy"`
	TdAccuracy string `xml:"td>accuracy"`
}

type CGPSSessionXML struct {
	StartTime            string                  `xml:"start-time"`
	StopTime             string                  `xml:"stop-time"`
	ObservationInterval  Number                  `xml:"observation-interval"`
	Operator             OperatorXML             `xml:"operator"`
	Rinex                RinexXML                `xml:"rinex"`
	DataFormat           string                  `xml:"data-format"`
	DownloadNameFormat   DownloadNameFormatXML   `xml:"download-name-format"`
	Receiver             ReceiverXML             `xml:"receiver"`
	InstalledCGPSAntenna InstalledCGPSAntennaXML `xml:"installed-cgps-antenna"`
	MetSensor            *MetSensor              `xml:"met-sensor,omitempty"`
}

func NewSiteXML(mark MarkXML, location LocationXML, sessions []CGPSSessionXML) SiteXML {
	return SiteXML{
		XMLName:  xml.Name{Local: "SITE"},
		Mark:     mark,
		Location: location,
		Sessions: sessions,
	}
}

func (x SiteXML) Marshal() ([]byte, error) {
	s, err := xml.MarshalIndent(x, "", "  ")
	if err != nil {
		return nil, err
	}
	return append([]byte(xml.Header), append(s, []byte{'\n', '\n'}...)...), nil
}
