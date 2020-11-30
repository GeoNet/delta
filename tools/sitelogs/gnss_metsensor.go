package main

import (
//"encoding/xml"
)

type GnssMetSensor struct {
	MetSensorModel                  string  `xml:"equip:type"`
	Manufacturer                    string  `xml:"equip:manufacturer"`
	SerialNumber                    string  `xml:"equip:serialNumber"`
	HeightDifftoAnt                 string  `xml:"equip:heightDiffToAntenna"`
	CalibrationDate                 string  `xml:"equip:calibrationDate"`
	EffectiveDates                  string  `xml:"equip:effectiveDates"`
	DataSamplingInterval            float64 `xml:"equip:dataSamplingInterval"`
	AccuracyPercentRelativeHumidity float64 `xml:"equip:accuracy-percentRelativeHumidity,omitempty"`
	AccuracyHPa                     float64 `xml:"equip:accuracy-hPa,omitempty"`
	AccuracyDegreesCelcius          float64 `xml:"equip:accuracy-degreesCelcius,omitempty"`
	Aspiration                      string  `xml:"equip:aspiration"`
	Notes                           string  `xml:"equip:notes"`
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
