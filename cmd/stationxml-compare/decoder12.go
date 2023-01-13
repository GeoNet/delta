package main

import (
	"bytes"
	"encoding/xml"
	"math"
	"sort"

	"golang.org/x/net/html/charset"

	"github.com/GeoNet/delta/internal/stationxml/v1.2"
)

type Decoder12 struct{}

func (d Decoder12) Decode(verbose bool, data []byte) ([]byte, error) {

	var root stationxml.FDSNStationXML
	decoder := xml.NewDecoder(bytes.NewReader(data))
	decoder.CharsetReader = charset.NewReaderLabel
	if err := decoder.Decode(&root); err != nil {
		return nil, err
	}

	for n := range root.Network {
		net := &root.Network[n]

		if net.TotalNumberStations != nil {
			net.TotalNumberStations = nil
		}

		if net.SelectedNumberStations != nil {
			net.SelectedNumberStations = nil
		}

		for s := range net.Station {
			sta := &net.Station[s]

			if sta.TotalNumberChannels != nil {
				sta.TotalNumberChannels = nil
			}
			if sta.SelectedNumberChannels != nil {
				sta.SelectedNumberChannels = nil
			}

			for c := range root.Network[n].Station[s].Channel {
				cha := &sta.Channel[c]

				if cha.ClockDrift != nil {
					cha.ClockDrift = nil
				}
				if cha.DataLogger != nil {
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
					stage.PolesZeros.Name = ""
					stage.PolesZeros.ResourceId = ""
				}
				for v := range cha.Response.Stage {
					stage := &cha.Response.Stage[v]
					if stage.Polynomial == nil {
						continue
					}
					stage.Polynomial.Name = ""
					stage.Polynomial.ResourceId = ""
					for i := range stage.Polynomial.Coefficient {
						stage.Polynomial.Coefficient[i].Number = stationxml.CounterType(i + 1)
					}
				}
				for v := range cha.Response.Stage {
					stage := &cha.Response.Stage[v]
					if stage.FIR == nil {
						continue
					}
					stage.FIR.Name = ""
					stage.FIR.ResourceId = ""
				}
				for v := range cha.Response.Stage {
					stage := &cha.Response.Stage[v]
					if stage.StageGain == nil {
						continue
					}

					stage.StageGain.Value = math.Round(stage.StageGain.Value)
				}
				for v := range cha.Response.Stage {
					stage := &cha.Response.Stage[v]
					if stage.Coefficients == nil {
						continue
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
				if cha.Response.InstrumentPolynomial != nil {
					cha.Response.InstrumentPolynomial.Name = ""
					cha.Response.InstrumentPolynomial.ResourceId = ""
					for i := range cha.Response.InstrumentPolynomial.Coefficient {
						cha.Response.InstrumentPolynomial.Coefficient[i].Number = stationxml.CounterType(i + 1)
					}
				}
			}

			for i := range sta.Comment {
				sta.Comment[i].Id = stationxml.CounterType(i + 1)
			}
			sort.Slice(root.Network[n].Station[s].Channel, func(i, j int) bool {
				switch {
				case root.Network[n].Station[s].Channel[i].LocationCode < root.Network[n].Station[s].Channel[j].LocationCode:
					return true
				case root.Network[n].Station[s].Channel[i].LocationCode > root.Network[n].Station[s].Channel[j].LocationCode:
					return false
				case root.Network[n].Station[s].Channel[i].StartDate.Before(root.Network[n].Station[s].Channel[j].StartDate.Time):
					return true
				case root.Network[n].Station[s].Channel[i].StartDate.After(root.Network[n].Station[s].Channel[j].StartDate.Time):
					return false
				case root.Network[n].Station[s].Channel[i].BaseNodeType.Code < root.Network[n].Station[s].Channel[j].BaseNodeType.Code:
					return true
				default:
					return false
				}
			})

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
