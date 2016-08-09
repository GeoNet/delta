package main

import (
	"sort"
)

type ResponseStage struct {
	Type        string  `yaml:"type"`
	Lookup      string  `yaml:"lookup"`
	Frequency   float64 `yaml:"frequency"`
	SampleRate  float64 `yaml:"samplerate"`
	Decimate    int32   `yaml:"decimate"`
	Gain        float64 `yaml:"gain"`
	Scale       float64 `yaml:"scale"`
	Correction  float64 `yaml:"correction"`
	Delay       float64 `yaml:"delay"`
	InputUnits  string  `yaml:"inputunits"`
	OutputUnits string  `yaml:"outputunits"`
}

type filterMap map[string][]ResponseStage

func (f filterMap) Keys() []string {
	var keys []string
	for k, _ := range f {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (f filterMap) Values() [][]ResponseStage {
	var values [][]ResponseStage
	for _, k := range f.Keys() {
		values = append(values, f[k])
	}
	return values
}
