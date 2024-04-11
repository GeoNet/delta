package main

import (
	"encoding/json"
	"io"
	"time"
)

// Config describes the processing settings for a geomag stream.
type Config struct {
	// Srcname can be used as a stream key
	Srcname string `json:"srcname"`
	// Network is the expected network code as stored.
	Network string `json:"network"`
	// Station is the expected station code as stored.
	Station string `json:"station"`
	// Location is the expected site location code as stored.
	Location string `json:"location"`
	// Channel is the expected channel code as stored.
	Channel string `json:"channel"`
	// ScaleBias is the offset that needs to be added to each data sample.
	ScaleBias float64 `json:"scale_bias"`
	// ScaleFactor is the value that needs to be multiplied to each data sample.
	ScaleFactor float64 `json:"scale_factor"`
	// InputUnits describes the units for the input signal.
	InputUnits string `json:"input_units"`
	// OutputUnits describes the units for the output after scaling.
	OutputUnits string `json:"output_units"`
	// Start is the time when the scale factors are valid.
	Start time.Time `json:"start"`
	// End is the time when the scale factors are no longer valid.
	End time.Time `json:"end"`
}

// Less can be used for sorting Config slices.
func (c Config) Less(config Config) bool {
	switch {
	case c.Srcname < config.Srcname:
		return true
	case c.Srcname > config.Srcname:
		return false
	case c.Start.Before(config.Start):
		return true
	default:
		return false
	}
}

// Encode will write JSON encoded output of a Config slice.
func Encode(wr io.Writer, d []Config) error {

	// build an encoder
	enc := json.NewEncoder(wr)
	enc.SetIndent("", "  ")

	// do the encoding
	return enc.Encode(d)
}
