package main

import (
	"encoding/xml"
)

type GnssAntenna struct {
	XMLName xml.Name `xml:"gnssAntenna"`

	AntennaType            string `xml:"equip:antennaType"`
	SerialNumber           string `xml:"equip:serialNumber"`
	AntennaReferencePoint  string `xml:"equip:antennaReferencePoint"`
	MarkerArpUpEcc         string `xml:"equip:marker-arpUpEcc."`
	MarkerArpNorthEcc      string `xml:"equip:marker-arpNorthEcc."`
	MarkerArpEastEcc       string `xml:"equip:marker-arpEastEcc."`
	AlignmentFromTrueNorth string `xml:"equip:alignmentFromTrueNorth"`
	AntennaRadomeType      string `xml:"equip:antennaRadomeType"`
	RadomeSerialNumber     string `xml:"equip:radomeSerialNumber"`
	AntennaCableType       string `xml:"equip:antennaCableType"`
	AntennaCableLength     string `xml:"equip:antennaCableLength"`
	DateInstalled          string `xml:"equip:dateInstalled"`
	DateRemoved            string `xml:"equip:dateRemoved"`
	Notes                  string `xml:"equip:notes"`
}

type GnssAntennaInput struct {
	XMLName xml.Name `xml:"gnssAntenna"`

	AntennaType            string `xml:"antennaType"`
	SerialNumber           string `xml:"serialNumber"`
	AntennaReferencePoint  string `xml:"antennaReferencePoint"`
	MarkerArpUpEcc         string `xml:"marker-arpUpEcc."`
	MarkerArpNorthEcc      string `xml:"marker-arpNorthEcc."`
	MarkerArpEastEcc       string `xml:"marker-arpEastEcc."`
	AlignmentFromTrueNorth string `xml:"alignmentFromTrueNorth"`
	AntennaRadomeType      string `xml:"antennaRadomeType"`
	RadomeSerialNumber     string `xml:"radomeSerialNumber"`
	AntennaCableType       string `xml:"antennaCableType"`
	AntennaCableLength     string `xml:"antennaCableLength"`
	DateInstalled          string `xml:"dateInstalled"`
	DateRemoved            string `xml:"dateRemoved"`
	Notes                  string `xml:"notes"`
}

type GnssAntennas []GnssAntenna

func (g GnssAntennas) Len() int      { return len(g) }
func (g GnssAntennas) Swap(i, j int) { g[i], g[j] = g[j], g[i] }
func (g GnssAntennas) Less(i, j int) bool {
	return g[i].DateInstalled < g[j].DateInstalled
}
