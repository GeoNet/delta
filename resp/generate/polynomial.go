package main

import (
	"sort"
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
