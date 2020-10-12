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

type SiteIdentificationInput struct {
	XMLName xml.Name `xml:"siteIdentification"`

	SiteName               string `xml:"siteName"`
	FourCharacterID        string `xml:"fourCharacterID"`
	MonumentInscription    string `xml:"monumentInscription"`
	IersDOMESNumber        string `xml:"iersDOMESNumber"`
	CdpNumber              string `xml:"cdpNumber"`
	MonumentDescription    string `xml:"monumentDescription"`
	HeightOfTheMonument    string `xml:"heightOfTheMonument"`
	MonumentFoundation     string `xml:"monumentFoundation"`
	FoundationDepth        string `xml:"foundationDepth"`
	MarkerDescription      string `xml:"markerDescription"`
	DateInstalled          string `xml:"dateInstalled"`
	GeologicCharacteristic string `xml:"geologicCharacteristic"`
	BedrockType            string `xml:"bedrockType"`
	BedrockCondition       string `xml:"bedrockCondition"`
	FractureSpacing        string `xml:"fractureSpacing"`
	FaultZonesNearby       string `xml:"faultZonesNearby"`
	DistanceActivity       string `xml:"distance-Activity"`
	Notes                  string `xml:"notes"`
}
