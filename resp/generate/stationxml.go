package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	stationxml "github.com/GeoNet/delta/internal/stationxml/v1.1"
)

// ResponseType mimics stationxml.ResponseType but adds an XMLName tag to correctly marshal the type name.
type ResponseType struct {
	XMLName               xml.Name                       `xml:"Response"`
	ResourceId            string                         `xml:"resourceId,attr,omitempty"`
	InstrumentSensitivity *stationxml.SensitivityType    `xml:"InstrumentSensitivity,omitempty"`
	InstrumentPolynomial  *stationxml.PolynomialType     `xml:"InstrumentPolynomial,omitempty"`
	Stage                 []stationxml.ResponseStageType `xml:"Stage,omitempty"`
}

// adjust gains based on current value
func runningGain(current, gain float64) float64 {
	switch {
	case current == 0.0:
		return gain
	case gain == 0.0:
		return current
	default:
		return current * gain
	}
}

// adjust units based on whether the current one has been set
func runningUnits(current, units string) string {
	if current == "" {
		return units
	}
	return current
}

// adds xml header to front of file
func writeXMLFile(path string, data []byte) error {
	const newline byte = '\n'

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.Write(append([]byte(xml.Header), append(data, newline)...)); err != nil {
		return err
	}

	return nil
}

// generate xml files in a given directory which should exists.
func (g Generate) build(dir string, dataloggers map[string]DataloggerModel, sensors map[string]SensorModel) error {

	// where to store xml output prior to writing
	files := make(map[string][]byte)

	for _, response := range g.ResponseMap {
		for _, sensor := range response.Sensors {
			for _, label := range sensor.Sensors {
				model, ok := sensors[label]
				if !ok {
					log.Printf("warning, unknown sensor: %s, skipping", label)
					continue
				}

				// NRL -- use the frequency of the gain element if posisble

				var stages []stationxml.ResponseStageType
				for _, filter := range sensor.Filters {
					m, ok := g.FilterMap[filter]
					if !ok {
						return fmt.Errorf("error, invalid filter: %s", filter)
					}

					// can only be one for a sensor response
					var sensitivity *stationxml.SensitivityType
					var polynomial *stationxml.PolynomialType

					// build response stages, usually only one for a sensor
					for count, stage := range m {
						switch stage.Type {
						case "paz":
							paz := g.PAZ(stage.Lookup)
							if paz == nil {
								return fmt.Errorf("invalid paz filter: %s", stage.Lookup)
							}

							stages = append(stages, paz.Stage(filter, stage, count+1, stage.Frequency, stage.Gain))
							if s := paz.Sensitivity(stage); sensitivity == nil {
								sensitivity = &s
							}
						case "poly":
							poly := g.Polynomial(stage.Lookup)
							if poly == nil {
								return fmt.Errorf("error, invalid poly filter: %s", stage.Lookup)
							}

							stages = append(stages, poly.Stage(filter, stage, count+1))
							if p := poly.Polynomial(filter, stage); polynomial == nil {
								polynomial = &p
							}
						default:
							return fmt.Errorf("error, unknown filter type: %s", stage.Type)
						}
					}

					resp := ResponseType{
						InstrumentSensitivity: sensitivity,
						InstrumentPolynomial:  polynomial,
						Stage:                 stages,
					}

					data, err := xml.MarshalIndent(resp, "", "  ")
					if err != nil {
						return err
					}

					files[g.SensorName(model, sensor, label)] = data
				}

				for _, datalogger := range response.Dataloggers {
					for _, label := range datalogger.Dataloggers {
						model, ok := dataloggers[label]
						if !ok {
							log.Printf("warning, unknown datalogger: %s, skipping", label)
							continue
						}

						final := datalogger.SampleRate
						if final < 0 {
							final = -1.0 / final
						}

						// NRL -- use the frequency of one half the Nyquist (FR/2) / 2
						freq := final / 4

						var stages []stationxml.ResponseStageType

						var gain, rate float64
						var inputUnits, outputUnits string
						for _, filter := range datalogger.Filters {
							m, ok := g.FilterMap[filter]
							if !ok {
								return fmt.Errorf("error, unknown filter: %s", filter)
							}
							for count, stage := range m {
								switch stage.Type {
								case "a2d":
									rate = stage.SampleRate

									stages = append(stages, A2D(filter, stage, count+1, freq, rate))

									// initial sample rates and units
									if v := float64(stage.Decimate); v > 1.0 {
										rate /= v
									}

									gain = runningGain(gain, stage.Gain)
									inputUnits = runningUnits(inputUnits, stage.InputUnits)
									outputUnits = stage.OutputUnits

								case "paz":
									paz := g.PAZ(stage.Lookup)
									if paz == nil {
										return fmt.Errorf("invalid paz filter: %s", stage.Lookup)
									}

									scale := stage.Gain * paz.Gain(freq) / paz.Gain(stage.Frequency)
									stages = append(stages, paz.Stage(filter, stage, count+1, freq, scale))

									gain = runningGain(gain, scale)
									inputUnits = runningUnits(inputUnits, stage.InputUnits)
									outputUnits = stage.OutputUnits

								case "fir":
									fir := g.FIR(stage.Lookup)
									if fir == nil {
										return fmt.Errorf("error, unknown fir filter %s", stage.Lookup)
									}

									stages = append(stages, fir.Stage(filter, stage, count+1, freq, rate))

									if v := float64(fir.Decimation); v > 1.0 {
										rate /= v
									}

									gain = runningGain(gain, stage.Gain)
									inputUnits = runningUnits(inputUnits, stage.InputUnits)
									outputUnits = stage.OutputUnits
								}
							}
						}

						resp := ResponseType{
							InstrumentSensitivity: &stationxml.SensitivityType{
								GainType: stationxml.GainType{
									Value: func() float64 {
										if gain != 0.0 {
											return gain
										}
										return 1.0
									}(),
									Frequency: freq,
								},
								InputUnits:  stationxml.UnitsType{Name: inputUnits},
								OutputUnits: stationxml.UnitsType{Name: outputUnits},
							},
							Stage: stages,
						}

						data, err := xml.MarshalIndent(resp, "", "  ")
						if err != nil {
							return err
						}

						files[g.DataloggerName(model, datalogger, strings.ReplaceAll(label, " ", ""), datalogger.SampleRate)] = data
					}
				}
			}
		}
	}

	for k, v := range files {
		if err := writeXMLFile(filepath.Join(dir, k+".xml"), v); err != nil {
			return err
		}
	}

	return nil
}
