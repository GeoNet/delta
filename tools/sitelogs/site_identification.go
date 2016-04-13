package main

import (
	"encoding/xml"
)

type SiteIdentification struct {
	XMLName xml.Name `xml:"siteIdentification"`

	SiteName               string `xml:"mi:siteName"`
	FourCharacterID        string `xml:"mi:fourCharacterID"`
	MonumentInscription    string `xml:"mi:monumentInscription"`
	IersDOMESNumber        string `xml:"mi:iersDOMESNumber"`
	CdpNumber              string `xml:"mi:cdpNumber"`
	MonumentDescription    string `xml:"mi:monumentDescription"`
	HeightOfTheMonument    string `xml:"mi:heightOfTheMonument"`
	MonumentFoundation     string `xml:"mi:monumentFoundation"`
	FoundationDepth        string `xml:"mi:foundationDepth"`
	MarkerDescription      string `xml:"mi:markerDescription"`
	DateInstalled          string `xml:"mi:dateInstalled"`
	GeologicCharacteristic string `xml:"mi:geologicCharacteristic"`
	BedrockType            string `xml:"mi:bedrockType"`
	BedrockCondition       string `xml:"mi:bedrockCondition"`
	FractureSpacing        string `xml:"mi:fractureSpacing"`
	FaultZonesNearby       string `xml:"mi:faultZonesNearby"`
	DistanceActivity       string `xml:"mi:distance-Activity"`
	Notes                  string `xml:"mi:notes"`
}
