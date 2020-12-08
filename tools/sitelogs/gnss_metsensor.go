package main

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
