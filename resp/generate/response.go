package main

import (
	"sort"
)

type Sensor struct {
	Sensors  []string `yaml:"sensors"`
	Filters  []string `yaml:"filters"`
	Channels string   `yaml:"channels"`
	Reversed bool     `yaml:"reversed"`
}

type Datalogger struct {
	Dataloggers   []string `yaml:"dataloggers"`
	Type          string   `yaml:"type"`
	Label         string   `yaml:"label"`
	SampleRate    float64  `yaml:"samplerate"`
	Frequency     float64  `yaml:"frequency"`
	StorageFormat string   `yaml:"storageformat"`
	ClockDrift    float64  `yaml:"clockdrift"`
	Filters       []string `yaml:"filters"`
	Reversed      bool     `yaml:"reversed"`
}

type Response struct {
	Sensors     []Sensor     `yaml:"sensors"`
	Dataloggers []Datalogger `yaml:"dataloggers"`
}

type responseMap map[string]Response

func (r responseMap) Keys() []string {
	var keys []string
	for k := range r {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (r responseMap) Values() []Response {
	var values []Response
	for _, k := range r.Keys() {
		values = append(values, r[k])
	}
	return values
}
