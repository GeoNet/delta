package main

import (
	"bytes"
	"encoding/xml"
	"log"
	"math"
	"sort"

	"golang.org/x/net/html/charset"

	"github.com/GeoNet/delta/internal/stationxml/v1.0"
)

type Decoder10 struct{}

func (d Decoder10) Decode(verbose bool, data []byte) ([]byte, error) {

	var root stationxml.FDSNStationXML
	decoder := xml.NewDecoder(bytes.NewReader(data))
	decoder.CharsetReader = charset.NewReaderLabel
	if err := decoder.Decode(&root); err != nil {
		return nil, err
	}

	for n := range root.Network {
		net := &root.Network[n]

		if net.TotalNumberStations != nil {
			if verbose {
				log.Println("found network TotalNumberStations")
			}
			net.TotalNumberStations = nil
		}

		if net.SelectedNumberStations != nil {
			if verbose {
				log.Println("found network SelectedNumberStations")
			}
			net.SelectedNumberStations = nil
		}

		for s := range net.Station {
			sta := &net.Station[s]

			if sta.TotalNumberChannels != nil {
				if verbose {
					log.Println("found station TotalNumberChannels")
				}
				sta.TotalNumberChannels = nil
			}
			if sta.SelectedNumberChannels != nil {
				if verbose {
					log.Println("found station SelectedNumberChannels")
				}
				sta.SelectedNumberChannels = nil
			}

			for c := range root.Network[n].Station[s].Channel {
				cha := &sta.Channel[c]

				if cha.ClockDrift != nil {
					if verbose {
						log.Println("found channel ClockDrift")
					}
					cha.ClockDrift = nil
				}
				if cha.DataLogger != nil {
					log.Println("found DataLogger")
					cha.DataLogger = nil
				}
				if cha.Sensor != nil {
					cha.Sensor = nil
				}

				for i := range cha.Comment {
					cha.Comment[i].Id = stationxml.CounterType(i + 1)
				}
				if cha.Response == nil {
					continue
				}
				for v := range cha.Response.Stage {
					stage := &cha.Response.Stage[v]
					if stage.PolesZeros == nil {
						continue
					}
					sort.Slice(stage.PolesZeros.Zero, func(i, j int) bool {
						switch {
						case stage.PolesZeros.Zero[i].Real.Value < stage.PolesZeros.Zero[j].Real.Value:
							return true
						case stage.PolesZeros.Zero[i].Real.Value > stage.PolesZeros.Zero[j].Real.Value:
							return false
						case stage.PolesZeros.Zero[i].Imaginary.Value < stage.PolesZeros.Zero[j].Imaginary.Value:
							return true
						default:
							return false
						}
					})
					for i := range stage.PolesZeros.Zero {
						stage.PolesZeros.Zero[i].Number = i + 1
					}
					sort.Slice(stage.PolesZeros.Pole, func(i, j int) bool {
						switch {
						case stage.PolesZeros.Pole[i].Real.Value < stage.PolesZeros.Pole[j].Real.Value:
							return true
						case stage.PolesZeros.Pole[i].Real.Value > stage.PolesZeros.Pole[j].Real.Value:
							return false
						case stage.PolesZeros.Pole[i].Imaginary.Value < stage.PolesZeros.Pole[j].Imaginary.Value:
							return true
						default:
							return false
						}
					})
					for i := range stage.PolesZeros.Pole {
						stage.PolesZeros.Pole[i].Number = i + 1
					}
					if verbose {
						//log.Println("found poles zeros Name")
						//log.Println("found poles zeros ResourceId")
					}
					stage.PolesZeros.Name = ""
					stage.PolesZeros.ResourceId = ""
				}
				for v := range cha.Response.Stage {
					stage := &cha.Response.Stage[v]
					if stage.Polynomial == nil {
						continue
					}
					if verbose {
						//log.Println("found polynomial Name")
						//log.Println("found polynomial ResourceId")
					}
					stage.Polynomial.Name = ""
					stage.Polynomial.ResourceId = ""
				}
				for v := range cha.Response.Stage {
					stage := &cha.Response.Stage[v]
					if stage.FIR == nil {
						continue
					}
					if verbose {
						//log.Println("found FIR Name")
						//log.Println("found FIR ResourceId")
					}
					stage.FIR.Name = ""
					stage.FIR.ResourceId = ""
				}
				for v := range cha.Response.Stage {
					stage := &cha.Response.Stage[v]

					stage.StageGain.Value = math.Round(stage.StageGain.Value)
				}
				for v := range cha.Response.Stage {
					stage := &cha.Response.Stage[v]
					if stage.Coefficients == nil {
						continue
					}
					if verbose {
						//log.Println("found Coefficients Name")
						//log.Println("found Coefficients ResourceId")
					}
					stage.Coefficients.Name = ""
					stage.Coefficients.ResourceId = ""
				}
				for v := range cha.Response.Stage {
					stage := &cha.Response.Stage[v]
					if stage.Decimation == nil {
						continue
					}
					stage.Decimation.Delay.Value = 0.0
					stage.Decimation.Correction.Value = 0.0
				}
				if cha.Response.InstrumentSensitivity != nil {
					cha.Response.InstrumentSensitivity.Value = math.Round(cha.Response.InstrumentSensitivity.Value)
				}
			}

			for i := range sta.Comment {
				sta.Comment[i].Id = stationxml.CounterType(i + 1)
			}

			net.StartDate = stationxml.DateTime{}
			net.EndDate = stationxml.DateTime{}
		}
		sort.Slice(root.Network[n].Station, func(i, j int) bool {
			return root.Network[n].Station[i].BaseNodeType.Code < root.Network[n].Station[j].BaseNodeType.Code
		})

	}

	if root.Created != nil {
		root.Created = nil
	}

	sort.Slice(root.Network, func(i, j int) bool {
		return root.Network[i].BaseNodeType.Code < root.Network[j].BaseNodeType.Code
	})

	res, err := xml.MarshalIndent(root, "", "  ")
	if err != nil {
		return nil, err
	}

	return append([]byte(xml.Header+"\n"), res...), nil
}
