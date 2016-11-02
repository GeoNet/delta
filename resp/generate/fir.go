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
