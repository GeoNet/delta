package main

import (
	"encoding/xml"
)

type GnssMetSensor struct {
	XMLName xml.Name `xml:"gnssMetSensor"`

	Manufacturer         string
	MetSensorModel       string
	SerialNumber         string
	DataSamplingInterval string
	Accuracy             string
	HeightDifftoAnt      string
	CalibrationDate      string
	EffectiveDates       string
	Notes                string
}

type GnssMetSensors []GnssMetSensor

func (g GnssMetSensors) Len() int      { return len(g) }
func (g GnssMetSensors) Swap(i, j int) { g[i], g[j] = g[j], g[i] }
func (g GnssMetSensors) Less(i, j int) bool {
	if g[i].EffectiveDates < g[j].EffectiveDates {
		return true
	}
	return false
}
