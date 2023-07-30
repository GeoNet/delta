package stationxml

import (
	"encoding/xml"
	"sort"
	"time"

	"github.com/GeoNet/delta/internal/stationxml/v1.0"
)

type Encoder10 struct{}

func Encode10(r Root) ([]byte, error) {
	return Encoder10{}.MarshalRoot(r)
}

func (e Encoder10) toDateTime(at time.Time) stationxml.DateTime {
	if time.Since(at) > 0 {
		return stationxml.DateTime{Time: at}
	}
	return stationxml.DateTime{}
}

func (e Encoder10) toDateTimePtr(at time.Time) *stationxml.DateTime {
	t := stationxml.DateTime{Time: at}
	return &t
}

func (e Encoder10) toRestrictedStatus(restricted bool) stationxml.RestrictedStatusType {
	switch {
	case restricted:
		return stationxml.ClosedRestrictedStatus
	default:
		return stationxml.OpenRestrictedStatus
	}
}

func (e Encoder10) toSampleRateRatio(f float64) *stationxml.SampleRateRatioType {
	switch {
	case f > 1.0:
		return &stationxml.SampleRateRatioType{
			NumberSamples: int(f),
			NumberSeconds: 1,
		}
	case f > 0.0:
		return &stationxml.SampleRateRatioType{
			NumberSamples: 1,
			NumberSeconds: int(1.0 / f),
		}
	default:
		return nil
	}
}

func (e Encoder10) Response(response *ResponseType) *stationxml.ResponseType {
	var stages []stationxml.ResponseStageType

	for _, s := range response.Stages {

		var pz *stationxml.PolesZerosType
		if s.PolesZeros != nil {
			var zeros []stationxml.PoleZeroType
			for _, z := range s.PolesZeros.Zeros {
				zeros = append(zeros, stationxml.PoleZeroType{
					Number:    z.Number,
					Real:      stationxml.FloatNoUnitType{Value: z.Real.Value},
					Imaginary: stationxml.FloatNoUnitType{Value: z.Imaginary.Value},
				})
			}

			var poles []stationxml.PoleZeroType
			for _, p := range s.PolesZeros.Poles {
				poles = append(poles, stationxml.PoleZeroType{
					Number:    p.Number,
					Real:      stationxml.FloatNoUnitType{Value: p.Real.Value},
					Imaginary: stationxml.FloatNoUnitType{Value: p.Imaginary.Value},
				})
			}

			pz = &stationxml.PolesZerosType{
				BaseFilterType: stationxml.BaseFilterType{
					ResourceId:  s.PolesZeros.ResourceId,
					Name:        s.PolesZeros.Name,
					Description: s.PolesZeros.Description,
					InputUnits: stationxml.UnitsType{
						Name:        s.PolesZeros.InputUnits.Name,
						Description: s.PolesZeros.InputUnits.Description,
					},
					OutputUnits: stationxml.UnitsType{
						Name:        s.PolesZeros.OutputUnits.Name,
						Description: s.PolesZeros.OutputUnits.Description,
					},
				},

				PzTransferFunctionType: stationxml.ToPzTransferFunctionType(s.PolesZeros.PzTransferFunctionType),
				NormalizationFactor:    s.PolesZeros.NormalizationFactor,
				NormalizationFrequency: stationxml.FrequencyType{
					FloatType: stationxml.FloatType{Value: s.PolesZeros.NormalizationFrequency},
				},

				Zero: zeros,
				Pole: poles,
			}
		}

		var coeffs *stationxml.CoefficientsType
		if s.Coefficients != nil {
			var nums []stationxml.FloatType
			for _, f := range s.Coefficients.Numerators {
				nums = append(nums, stationxml.FloatType{Value: f.Value})
			}

			var denoms []stationxml.FloatType
			for _, f := range s.Coefficients.Denominators {
				denoms = append(denoms, stationxml.FloatType{Value: f.Value})
			}

			coeffs = &stationxml.CoefficientsType{
				BaseFilterType: stationxml.BaseFilterType{
					ResourceId:  s.Coefficients.ResourceId,
					Name:        s.Coefficients.Name,
					Description: s.Coefficients.Description,
					InputUnits: stationxml.UnitsType{
						Name:        s.Coefficients.InputUnits.Name,
						Description: s.Coefficients.InputUnits.Description,
					},
					OutputUnits: stationxml.UnitsType{
						Name:        s.Coefficients.OutputUnits.Name,
						Description: s.Coefficients.OutputUnits.Description,
					},
				},
				CfTransferFunctionType: stationxml.ToCfTransferFunctionType(s.Coefficients.CfTransferFunctionType),
				Numerator:              nums,
				Denominator:            denoms,
			}
		}

		var fir *stationxml.FIRType
		if s.FIR != nil {
			var coeffs []stationxml.NumeratorCoefficient
			for _, c := range s.FIR.NumeratorCoefficients {
				coeffs = append(coeffs, stationxml.NumeratorCoefficient{
					I:     c.I,
					Value: c.Value,
				})
			}

			fir = &stationxml.FIRType{
				BaseFilterType: stationxml.BaseFilterType{
					ResourceId:  s.FIR.ResourceId,
					Name:        s.FIR.Name,
					Description: s.FIR.Description,
					InputUnits: stationxml.UnitsType{
						Name:        s.FIR.InputUnits.Name,
						Description: s.FIR.InputUnits.Description,
					},
					OutputUnits: stationxml.UnitsType{
						Name:        s.FIR.OutputUnits.Name,
						Description: s.FIR.OutputUnits.Description,
					},
				},

				Symmetry: stationxml.ToSymmetry(s.FIR.Symmetry),

				NumeratorCoefficient: coeffs,
			}
		}

		var poly *stationxml.PolynomialType
		if s.Polynomial != nil {
			approx := stationxml.ToApproximationType(s.Polynomial.ApproximationType)

			var coeffs []stationxml.Coefficient
			for _, c := range s.Polynomial.Coefficients {
				coeffs = append(coeffs, stationxml.Coefficient{
					Number:          stationxml.CounterType(c.Number),
					FloatNoUnitType: stationxml.FloatNoUnitType{Value: c.Value},
				})
			}

			poly = &stationxml.PolynomialType{
				BaseFilterType: stationxml.BaseFilterType{
					ResourceId:  s.Polynomial.ResourceId,
					Name:        s.Polynomial.Name,
					Description: s.Polynomial.Description,
					InputUnits: stationxml.UnitsType{
						Name:        s.Polynomial.InputUnits.Name,
						Description: s.Polynomial.InputUnits.Description,
					},
					OutputUnits: stationxml.UnitsType{
						Name:        s.Polynomial.OutputUnits.Name,
						Description: s.Polynomial.OutputUnits.Description,
					},
				},

				FrequencyLowerBound: stationxml.FrequencyType{
					FloatType: stationxml.FloatType{Value: s.Polynomial.FrequencyLowerBound},
				},
				FrequencyUpperBound: stationxml.FrequencyType{
					FloatType: stationxml.FloatType{Value: s.Polynomial.FrequencyUpperBound},
				},
				ApproximationType:       &approx,
				ApproximationLowerBound: s.Polynomial.ApproximationLowerBound.Value,
				ApproximationUpperBound: s.Polynomial.ApproximationUpperBound.Value,
				MaximumError:            s.Polynomial.MaximumError,
				Coefficient:             coeffs,
			}
		}

		var dec *stationxml.DecimationType
		if s.Decimation != nil {
			dec = &stationxml.DecimationType{
				InputSampleRate: stationxml.FrequencyType{
					FloatType: stationxml.FloatType{Value: s.Decimation.InputSampleRate},
				},
				Factor:     s.Decimation.Factor,
				Offset:     s.Decimation.Offset,
				Delay:      stationxml.FloatType{Value: s.Decimation.Delay},
				Correction: stationxml.FloatType{Value: s.Decimation.Correction},
			}
		}

		var gain stationxml.GainType

		// in v1.0 we don't used a base PolynomialType
		if s.Polynomial != nil {

			// assume gain is related to the second coefficient
			value := 1.0
			if len(s.Polynomial.Coefficients) > 1 {
				value = 1.0 / s.Polynomial.Coefficients[1].Value
			}

			gain = stationxml.GainType{
				Value:     value,
				Frequency: response.frequency,
			}
		}

		if s.StageGain != nil {
			gain = stationxml.GainType{
				Value:     s.StageGain.Value,
				Frequency: s.StageGain.Frequency,
			}
		}

		stages = append(stages, stationxml.ResponseStageType{
			Number: stationxml.CounterType(s.Number),

			PolesZeros:   pz,
			Coefficients: coeffs,
			FIR:          fir,
			Polynomial:   poly,
			Decimation:   dec,
			StageGain:    gain,
		})
	}

	var sensitivity *stationxml.SensitivityType
	if response.InstrumentSensitivity != nil {
		sensitivity = &stationxml.SensitivityType{
			GainType: stationxml.GainType{
				Value:     response.InstrumentSensitivity.Value,
				Frequency: response.InstrumentSensitivity.Frequency,
			},
			InputUnits: stationxml.UnitsType{
				Name:        response.InstrumentSensitivity.InputUnits.Name,
				Description: response.InstrumentSensitivity.InputUnits.Description,
			},
			OutputUnits: stationxml.UnitsType{
				Name:        response.InstrumentSensitivity.OutputUnits.Name,
				Description: response.InstrumentSensitivity.OutputUnits.Description,
			},
		}
	}

	// in v1.0 we haven't used a base PolynomialType
	var polynomial *stationxml.PolynomialType
	if response.InstrumentPolynomial != nil {

		gain := 1.0
		for _, s := range response.Stages {
			if s.Polynomial != nil {
				continue
			}
			if s.StageGain == nil {
				continue
			}
			gain *= s.StageGain.Value
		}

		var value float64
		for _, s := range response.Stages {
			if s.Polynomial == nil {
				continue
			}
			if len(s.Polynomial.Coefficients) > 1 {
				value = gain / s.Polynomial.Coefficients[1].Value
			}
		}
		sensitivity = &stationxml.SensitivityType{
			GainType: stationxml.GainType{
				Value:     value,
				Frequency: response.frequency,
			},
			InputUnits: stationxml.UnitsType{
				Name:        response.InstrumentPolynomial.InputUnits.Name,
				Description: response.InstrumentPolynomial.InputUnits.Description,
			},
			OutputUnits: stationxml.UnitsType{
				Name:        response.InstrumentPolynomial.OutputUnits.Name,
				Description: response.InstrumentPolynomial.OutputUnits.Description,
			},
		}
	}

	return &stationxml.ResponseType{
		ResourceId: response.ResourceId,

		InstrumentSensitivity: sensitivity,
		InstrumentPolynomial:  polynomial,

		Stage: stages,
	}
}

func (e Encoder10) Channel(root Root, network Network, channel Channel) []stationxml.ChannelType {
	var channels []stationxml.ChannelType

	for _, stream := range channel.Streams {

		var response *stationxml.ResponseType
		if stream.Response != nil {
			response = e.Response(stream.Response)
		}

		var types []stationxml.Type
		switch {
		case stream.Triggered:
			types = append(types, stationxml.Triggered)
		default:
			types = append(types, stationxml.Continuous)
		}
		for _, t := range stream.Types {
			switch t {
			case 'G':
				types = append(types, stationxml.Geophysical)
			case 'W':
				types = append(types, stationxml.Weather)
			case 'H':
				types = append(types, stationxml.Health)
			}
		}

		var depth float64
		if stream.Vertical != 0.0 {
			depth = -stream.Vertical
		}

		datalogger := &stationxml.EquipmentType{
			ResourceId:       "Datalogger#" + stream.Datalogger.Model + ":" + stream.Datalogger.SerialNumber,
			Type:             stream.Datalogger.Type,
			Description:      stream.Datalogger.Description,
			Manufacturer:     stream.Datalogger.Manufacturer,
			Model:            stream.Datalogger.Model,
			SerialNumber:     stream.Datalogger.SerialNumber,
			InstallationDate: e.toDateTime(stream.Datalogger.InstallationDate),
			RemovalDate:      e.toDateTime(stream.Datalogger.RemovalDate),
		}

		channels = append(channels, stationxml.ChannelType{
			BaseNodeType: stationxml.BaseNodeType{
				Code:             stream.Code,
				RestrictedStatus: e.toRestrictedStatus(network.Restricted),
				StartDate:        e.toDateTime(stream.StartDate),
				EndDate:          e.toDateTime(stream.EndDate),
				Comment: []stationxml.CommentType{
					{
						Id:    stationxml.CounterType(1),
						Value: ToSiteSurvey(channel.Survey),
					},
					{
						Id:    stationxml.CounterType(2),
						Value: "Location is given in " + channel.Datum,
					},
					{
						Id:    stationxml.CounterType(3),
						Value: "Sensor orientation not known",
					},
				},
			},
			LocationCode: channel.LocationCode,
			Latitude: stationxml.LatitudeType{
				LatitudeBaseType: stationxml.LatitudeBaseType{
					FloatType: stationxml.FloatType{
						Value: channel.Latitude,
					},
				},
				Datum: channel.Datum,
			},
			Longitude: stationxml.LongitudeType{
				LongitudeBaseType: stationxml.LongitudeBaseType{
					FloatType: stationxml.FloatType{
						Value: channel.Longitude,
					},
				},
				Datum: channel.Datum,
			},
			Elevation:       stationxml.DistanceType{FloatType: stationxml.FloatType{Value: channel.Elevation}},
			Depth:           stationxml.DistanceType{FloatType: stationxml.FloatType{Value: depth}},
			Azimuth:         &stationxml.AzimuthType{FloatType: stationxml.FloatType{Value: stream.Azimuth}},
			Dip:             &stationxml.DipType{FloatType: stationxml.FloatType{Value: stream.Dip}},
			Type:            types,
			SampleRate:      stationxml.SampleRateType{FloatType: stationxml.FloatType{Value: stream.SamplingRate}},
			SampleRateRatio: e.toSampleRateRatio(stream.SamplingRate),
			Sensor: &stationxml.EquipmentType{
				ResourceId:       "Sensor#" + stream.Sensor.Model + ":" + stream.Sensor.SerialNumber,
				Type:             stream.Sensor.Type,
				Description:      stream.Sensor.Description,
				Manufacturer:     stream.Sensor.Manufacturer,
				Model:            stream.Sensor.Model,
				SerialNumber:     stream.Sensor.SerialNumber,
				InstallationDate: e.toDateTime(stream.Sensor.InstallationDate),
				RemovalDate:      e.toDateTime(stream.Sensor.RemovalDate),
			},
			DataLogger: datalogger,
			Response:   response,
		})
	}

	return channels
}

func (e Encoder10) Station(root Root, external External, network Network, station Station) stationxml.StationType {

	var channels []stationxml.ChannelType
	for _, c := range station.Channels {
		channels = append(channels, e.Channel(root, network, c)...)
	}

	sort.Slice(channels, func(i, j int) bool {
		switch {
		case channels[i].LocationCode < channels[j].LocationCode:
			return true
		case channels[i].LocationCode > channels[j].LocationCode:
			return false
		case channels[i].StartDate.Before(channels[j].StartDate.Time):
			return true
		case channels[i].StartDate.After(channels[j].StartDate.Time):
			return false
		case channels[i].SampleRate.Value > channels[j].SampleRate.Value:
			return true
		case channels[i].SampleRate.Value < channels[j].SampleRate.Value:
			return false
		case channels[i].Code < channels[j].Code:
			return true
		default:
			return false
		}
	})

	return stationxml.StationType{
		BaseNodeType: stationxml.BaseNodeType{
			Code:             station.Code,
			Description:      network.Description,
			RestrictedStatus: e.toRestrictedStatus(network.Restricted),
			StartDate:        e.toDateTime(station.StartDate),
			EndDate:          e.toDateTime(station.EndDate),
			Comment: []stationxml.CommentType{{
				Id:    stationxml.CounterType(1),
				Value: "Location is given in " + station.Datum,
			}},
		},
		Latitude: stationxml.LatitudeType{LatitudeBaseType: stationxml.LatitudeBaseType{
			FloatType: stationxml.FloatType{
				Value: station.Latitude,
			}},
			Datum: station.Datum,
		},
		Longitude: stationxml.LongitudeType{LongitudeBaseType: stationxml.LongitudeBaseType{
			FloatType: stationxml.FloatType{
				Value: station.Longitude,
			}},
			Datum: station.Datum,
		},
		Elevation: stationxml.DistanceType{
			FloatType: stationxml.FloatType{Value: station.Elevation},
		},
		Site: stationxml.SiteType{
			Name:        station.Name,
			Description: station.Description,
		},
		CreationDate:    e.toDateTime(station.CreationDate),
		TerminationDate: e.toDateTime(station.TerminationDate),
		Channel:         channels,
	}
}

func (e Encoder10) Network(root Root, external External) stationxml.NetworkType {

	var stations []stationxml.StationType
	for _, network := range external.Networks {
		for _, station := range network.Stations {
			stations = append(stations, e.Station(root, external, network, station))
		}
	}

	sort.Slice(stations, func(i, j int) bool {
		return stations[i].BaseNodeType.Code < stations[j].BaseNodeType.Code
	})

	return stationxml.NetworkType{
		BaseNodeType: stationxml.BaseNodeType{
			Code:        external.Code,
			Description: external.Description,
			RestrictedStatus: func() stationxml.RestrictedStatusType {
				switch external.Restricted {
				case true:
					return stationxml.ClosedRestrictedStatus
				default:
					return stationxml.OpenRestrictedStatus
				}
			}(),
			StartDate: e.toDateTime(external.Start()),
			EndDate:   e.toDateTime(external.End()),
		},
		Station: stations,
	}
}

func (e Encoder10) MarshalRoot(root Root) ([]byte, error) {

	type FDSNStationXML struct {
		stationxml.RootType

		NameSpace string `xml:"xmlns,attr"`
	}

	var created *stationxml.DateTime
	if root.Create {
		created = e.toDateTimePtr(time.Now().UTC())
	}

	var networks []stationxml.NetworkType
	for _, ext := range root.Externals {
		networks = append(networks, e.Network(root, ext))
	}

	sort.Slice(networks, func(i, j int) bool {
		return networks[i].BaseNodeType.Code < networks[j].BaseNodeType.Code
	})

	r := FDSNStationXML{
		NameSpace: "http://www.fdsn.org/xml/station/1",
		RootType: stationxml.RootType{
			SchemaVersion: 1.0,

			Source: root.Source,
			Sender: root.Sender,
			Module: root.Module,

			Network: networks,

			Created: created,
		},
	}

	h := xml.Header
	b, err := xml.MarshalIndent(r, "", "  ")
	if err != nil {
		return nil, err
	}

	return append([]byte(h), append(b, '\n')...), nil
}
