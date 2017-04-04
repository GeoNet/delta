package main

import (
	"sort"
	"strings"
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

type firMap map[string]FIR

func (f firMap) Keys() []string {
	var keys []string
	for k, _ := range f {
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
