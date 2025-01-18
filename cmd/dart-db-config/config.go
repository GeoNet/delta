package main

import (
	"encoding/json"
	"io"
	"time"
)

// Constituent values need by the de-tiding algorithm
type Constituent struct {
	Name      string  `json:"name"`
	Number    int     `json:"number"`
	Amplitude float64 `json:"amplitude"`
	Lag       float64 `json:"lag"`
}

// Detide parameters need by the de-tiding algorithm
type Detide struct {
	TimeZone float64 `json:"timezone"`
	Latitude float64 `json:"latitude"`

	Constituents []Constituent `json:"constituents,omitempty"`
}

// Deployment documents an individual DART operational period with optional de-tiding parameters.
type Deployment struct {
	Network          string        `json:"network"`
	Buoy             string        `json:"buoy"`
	Location         string        `json:"location"`
	Latitude         float64       `json:"latitude"`
	Longitude        float64       `json:"longitude"`
	Depth            float64       `json:"depth"`
	Detide           *Detide       `json:"detide,omitempty"`
	TimingCorrection time.Duration `json:"timing-correction,omitempty"`
	Start            time.Time     `json:"start"`
	End              time.Time     `json:"end"`
}

func Encode(wr io.Writer, d []Deployment) error {

	// build an encoder
	enc := json.NewEncoder(wr)
	enc.SetIndent("", "  ")

	// do the encoding
	return enc.Encode(d)
}
