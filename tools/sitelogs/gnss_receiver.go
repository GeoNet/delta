package main

import (
	"encoding/xml"
)

type GnssReceiver struct {
	XMLName xml.Name `xml:"gnssReceiver"`

	ReceiverType             string `xml:"equip:receiverType"`
	SatelliteSystem          string `xml:"equip:satelliteSystem"`
	SerialNumber             string `xml:"equip:serialNumber"`
	FirmwareVersion          string `xml:"equip:firmwareVersion"`
	ElevationCutoffSetting   string `xml:"equip:elevationCutoffSetting"`
	DateInstalled            string `xml:"equip:dateInstalled"`
	DateRemoved              string `xml:"equip:dateRemoved"`
	TemperatureStabilization string `xml:"equip:temperatureStabilization"`
	Notes                    string `xml:"equip:notes"`
}

type GnssReceiverInput struct {
	XMLName xml.Name `xml:"gnssReceiver"`

	ReceiverType             string `xml:"receiverType"`
	SatelliteSystem          string `xml:"satelliteSystem"`
	SerialNumber             string `xml:"serialNumber"`
	FirmwareVersion          string `xml:"firmwareVersion"`
	ElevationCutoffSetting   string `xml:"elevationCutoffSetting"`
	DateInstalled            string `xml:"dateInstalled"`
	DateRemoved              string `xml:"dateRemoved"`
	TemperatureStabilization string `xml:"temperatureStabilization"`
	Notes                    string `xml:"notes"`
}

type GnssReceivers []GnssReceiver

func (g GnssReceivers) Len() int { return len(g) }

func (g GnssReceivers) Swap(i, j int) { g[i], g[j] = g[j], g[i] }
func (g GnssReceivers) Less(i, j int) bool {
	if g[i].DateInstalled < g[j].DateInstalled {
		return true
	}
	return false
}
