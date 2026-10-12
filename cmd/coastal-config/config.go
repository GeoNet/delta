package main

import (
	"encoding/json"
	"io"
	"time"
)

// Coastal configures an individual coastal tsunami site including de-tiding parameters.
type Coastal struct {
	// Network is the external network code
	Network string `json:"network"`
	// Station is the station code
	Station string `json:"station"`
	// Location is the site location code
	Location string `json:"location"`
	// Latitude is the site latitude
	Latitude float64 `json:"latitude"`
	// Longitude is the site longitude
	Longitude float64 `json:"longitude"`
	// Factor is scale factor for raw miniseed samples
	Factor float64 `json:"factor"`
	// Bias is the offset to scaled data samples
	Bias float64 `json:"bias"`
	// Units is the input units of the scale factor and bias.
	Units string `json:"units"`
	// Detide holds the tidal constituents and details used for detiding
	Detide *Detide `json:"detide,omitempty"`
	// Start is the time from which the configuration is valid
	Start time.Time `json:"start"`
	// End is the time before which the configuration was valid
	End time.Time `json:"end"`
}

func Encode(wr io.Writer, c []Coastal) error {

	// build an encoder
	enc := json.NewEncoder(wr)
	enc.SetIndent("", "  ")

	// do the encoding
	return enc.Encode(c)
}
