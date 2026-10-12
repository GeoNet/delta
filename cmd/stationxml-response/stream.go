package main

import (
	"time"
)

// Stream holds the simplified response information.
// It assumes that there is only one poles and zeros
// entry as well as a list of decimation filters
// that have been used to generate the output samples.
// The actual filters are represented by names to allow
// for reusing repeated values.
type Stream struct {
	Srcname string `json:"srcname"`

	Sensor     string `json:"sensor"`
	Datalogger string `json:"datalogger"`

	InputUnits  string `json:"input_units"`
	OutputUnits string `json:"output_units"`

	Sensitivity float64 `json:"sensitivity"`
	Frequency   float64 `json:"frequency"`
	SampleRate  float64 `json:"sample_rate"`

	PolesZeros        string             `json:"poles_zeros"`
	DecimationFilters []DecimationFilter `json:"decimation_filters"`

	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}
