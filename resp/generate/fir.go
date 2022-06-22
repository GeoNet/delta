package main

import (
	"sort"
	"strings"

	"github.com/GeoNet/delta/internal/stationxml/v1.1"
)

type FIR struct {
	Causal     bool      `yaml:"causal"`
	Symmetry   string    `yaml:"symmetry"`
	Decimation float64   `yaml:"decimation"`
	Gain       float64   `yaml:"gain"`
	Notes      string    `yaml:"notes"`
	Factors    []float64 `yaml:"factors"`
}

func (f FIR) SymmetryLookup() string {

	switch strings.ToUpper(f.Symmetry) {
	case "EVEN":
		return "SymmetryEven"
	case "ODD":
		return "SymmetryOdd"
	default:
		return "SymmetryNone"
	}
}

func (f FIR) SymmetryType() stationxml.Symmetry {

	switch strings.ToUpper(f.Symmetry) {
	case "EVEN":
		return stationxml.EvenSymmetry
	case "ODD":
		return stationxml.OddSymmetry
	default:
		return stationxml.NoneSymmetry
	}
}

func (f FIR) Correction(sps float64) float64 {

	switch strings.ToUpper(f.Symmetry) {
	case "EVEN":
		return float64(len(f.Factors)*2-1) / (2.0 * sps)
	case "ODD":
		return float64((len(f.Factors)-1)*2+1-1) / (2.0 * sps)
	default:
		return float64(len(f.Factors)-1) / (2.0 * sps)
	}
}

func (fir FIR) Stage(filter string, stage ResponseStage, count int, freq, rate float64) stationxml.ResponseStageType {

	var corr float64
	if v := float64(fir.Decimation); v > 1.0 {
		if corr = stage.Correction; corr == 0.0 {
			corr = fir.Correction(rate / v)
		}
	}

	var coeffs []stationxml.NumeratorCoefficient
	for n, c := range fir.Factors {
		coeffs = append(coeffs, stationxml.NumeratorCoefficient{
			I:     n + 1,
			Value: c,
		})
	}

	return stationxml.ResponseStageType{
		Number: stationxml.CounterType(count),
		FIR: &stationxml.FIRType{
			BaseFilterType: stationxml.BaseFilterType{
				ResourceId:  "FIR#" + filter,
				Name:        stage.Lookup,
				Description: "",

				InputUnits: stationxml.UnitsType{
					Name: stage.InputUnits,
				},
				OutputUnits: stationxml.UnitsType{
					Name: stage.OutputUnits,
				},
			},
			Symmetry:             fir.SymmetryType(),
			NumeratorCoefficient: coeffs,
		},
		Decimation: &stationxml.DecimationType{
			InputSampleRate: stationxml.FrequencyType{
				FloatType: stationxml.FloatType{
					Value: rate,
				},
			},
			Factor: int(fir.Decimation),
			Delay: stationxml.FloatType{
				Value: corr,
			},
			Correction: stationxml.FloatType{
				Value: corr,
			},
		},
		StageGain: &stationxml.GainType{
			Value: func() float64 {
				if stage.Gain != 0.0 {
					return stage.Gain
				}
				return 1.0
			}(),
			Frequency: freq,
		},
	}
}

type firMap map[string]FIR

func (f firMap) Keys() []string {
	var keys []string
	for k := range f {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (f firMap) Values() []FIR {
	var values []FIR
	for _, k := range f.Keys() {
		values = append(values, f[k])
	}
	return values
}
