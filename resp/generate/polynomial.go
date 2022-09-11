package main

import (
	"sort"

	"github.com/GeoNet/delta/internal/stationxml/v1.1"
)

type Polynomial struct {
	Gain                    float64 `yaml:"gain"`
	ApproximationType       string  `yaml:"approximationtype"`
	FrequencyLowerBound     float64 `yaml:"frequencylowerbound"`
	FrequencyUpperBound     float64 `yaml:"frequencyupperbound"`
	ApproximationLowerBound float64 `yaml:"approximationlowerbound"`
	ApproximationUpperBound float64 `yaml:"approximationupperbound"`
	MaximumError            float64 `yaml:"maximumerror"`
	Notes                   string  `yaml:"notes"`

	Coefficients []float64 `yaml:"coefficients"`
}

func (p Polynomial) ApproximationTypeLookup() string {
	switch p.ApproximationType {
	case "MACLAURIN":
		return "ApproximationTypeMaclaurin"
	default:
		return "ApproximationTypeUnknown"
	}
}

func (p Polynomial) Approximation() stationxml.ApproximationType {
	switch p.ApproximationType {
	case "MACLAURIN":
		return stationxml.MaclaurinApproximation
	default:
		return stationxml.ApproximationType(0)
	}
}

func (p Polynomial) Polynomial(filter string, stage ResponseStage) stationxml.PolynomialType {

	var coeffs []stationxml.Coefficient
	for n, c := range p.Coefficients {
		coeffs = append(coeffs, stationxml.Coefficient{
			Number: stationxml.CounterType(n),
			FloatNoUnitType: stationxml.FloatNoUnitType{
				Value: c,
			},
		})
	}

	// there can only be one polynomail
	return stationxml.PolynomialType{
		BaseFilterType: stationxml.BaseFilterType{
			ResourceId:  "Polynomial#" + filter,
			InputUnits:  stationxml.UnitsType{Name: stage.InputUnits},
			OutputUnits: stationxml.UnitsType{Name: stage.OutputUnits},
		},
		ApproximationType: func() *stationxml.ApproximationType {
			if v := stationxml.ToApproximationType(p.ApproximationType); v > 0 {
				return &v
			}
			return nil
		}(),
		FrequencyLowerBound:     stationxml.FrequencyType{FloatType: stationxml.FloatType{Value: p.FrequencyLowerBound}},
		FrequencyUpperBound:     stationxml.FrequencyType{FloatType: stationxml.FloatType{Value: p.FrequencyUpperBound}},
		ApproximationLowerBound: p.ApproximationLowerBound,
		ApproximationUpperBound: p.ApproximationUpperBound,
		MaximumError:            p.MaximumError,
		Coefficient:             coeffs,
	}
}

func (p Polynomial) Stage(filter string, stage ResponseStage, count int) stationxml.ResponseStageType {

	var coeffs []stationxml.Coefficient
	for n, c := range p.Coefficients {
		coeffs = append(coeffs, stationxml.Coefficient{
			Number: stationxml.CounterType(n),
			FloatNoUnitType: stationxml.FloatNoUnitType{
				Value: c,
			},
		})
	}

	return stationxml.ResponseStageType{
		Number: stationxml.CounterType(count),
		Polynomial: &stationxml.PolynomialType{
			BaseFilterType: stationxml.BaseFilterType{
				ResourceId:  "Polynomial#" + filter,
				InputUnits:  stationxml.UnitsType{Name: stage.InputUnits},
				OutputUnits: stationxml.UnitsType{Name: stage.OutputUnits},
			},
			ApproximationType: func() *stationxml.ApproximationType {
				v := p.Approximation()
				return &v
			}(),
			FrequencyLowerBound:     stationxml.FrequencyType{FloatType: stationxml.FloatType{Value: p.FrequencyLowerBound}},
			FrequencyUpperBound:     stationxml.FrequencyType{FloatType: stationxml.FloatType{Value: p.FrequencyUpperBound}},
			ApproximationLowerBound: p.ApproximationLowerBound,
			ApproximationUpperBound: p.ApproximationUpperBound,
			MaximumError:            p.MaximumError,
			Coefficient:             coeffs,
		},
	}

}

type polynomialMap map[string]Polynomial

func (p polynomialMap) Keys() []string {
	var keys []string
	for k := range p {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (p polynomialMap) Values() []Polynomial {
	var values []Polynomial
	for _, k := range p.Keys() {
		values = append(values, p[k])
	}
	return values
}
