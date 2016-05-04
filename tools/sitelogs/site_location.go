package main

import (
	"encoding/xml"
)

type ApproximatePositionITRF struct {
	XMLName xml.Name `xml:"mi:approximatePositionITRF"`

	XCoordinateInMeters string `xml:"mi:xCoordinateInMeters"`
	YCoordinateInMeters string `xml:"mi:yCoordinateInMeters"`
	ZCoordinateInMeters string `xml:"mi:zCoordinateInMeters"`
	LatitudeNorth       string `xml:"mi:latitude-North"`
	LongitudeEast       string `xml:"mi:longitude-East"`
	ElevationMEllips    string `xml:"mi:elevation-m_ellips."`
}

type SiteLocation struct {
	XMLName xml.Name `xml:"siteLocation"`

	City                    string `xml:"mi:city"`
	State                   string `xml:"mi:state"`
	Country                 string `xml:"mi:country"`
	TectonicPlate           string `xml:"mi:tectonicPlate"`
	ApproximatePositionITRF ApproximatePositionITRF
	Notes                   string `xml:"mi:notes"`
}
