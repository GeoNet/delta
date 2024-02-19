package main

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"time"
)

// Stream is used to identify and describe streams that may be used for chart drawing.
type Stream struct {
	// Srcname is the Stream identification encoded as per NN_SSSS_LL_CCC.
	Srcname string `json:"srcname"`
	// NetworkCode is the short code of Stream network (NN).
	NetworkCode string `json:"network-code"`
	// StationCode is the short code of the Stream station (SSSS).
	StationCode string `json:"station-code"`
	// LocationCode is the short code of the Stream site location (LL).
	LocationCode string `json:"location-code"`
	// ChannelCode is the short code of the Stream channel (CCC).
	ChannelCode string `json:"channel-code"`
	// StationName is the long name of the Stream station.
	StationName string `json:"station-name"`
	// InternalNetwok is the shore code of the Stream network lookup.
	InternalNetwork string `json:"internal-network"`
	// NetworkDescription is the long name of the Stream internal network.
	NetworkDescription string `json:"network-description,omitempty"`
	// Latitude is the Stream site latitude in decimal degrees.
	Latitude float64 `json:"latitude"`
	// Longitude is the Stream site longitude in decimal degrees.
	Longitude float64 `json:"longitude"`
	// Elevation is the Stream site height in meters above sea level.
	Elevation float64 `json:"elevation,omitempty"`
	// Depth is the Stream site depth in meters below water level.
	Depth float64 `json:"depth,omitempty"`
	// SamplingPeriod is the time interval between samples, in nanoseconds.
	SamplingPeriod time.Duration `json:"sampling-period,omitempty"`
	// TidalLag is lag of the primary tide for the site if appropriate.
	TidalLag float64 `json:"tidal-lag,omitempty"`
	// Sensitivity is the initial conversion factor to convert from counts to the desired units, or volts.
	Sensitivity float64 `json:"sensitivity,omitempty"`
	// Gain is an optional factor to convert from, usually, volts to the desired units.
	Gain float64 `json:"gain,omitempty"`
	// Bias is an optional factor to add to the output after scaling to offset the signal.
	Bias float64 `json:"bias,omitempty"`
	// InputUnits describes the expected input signal units.
	InputUnits string `json:"input_units,omitempty"`
	// OutputUnits describes the expected output signal units, which is usually counts.
	OutputUnits string `json:"output_units,omitempty"`
}

// Config is a slice of Stream vales.
type Config []Stream

// Write sends a JSON encoded byte array to a Writer interface.
func (c Config) Write(wr io.Writer) error {
	enc := json.NewEncoder(wr)

	return enc.Encode(c)
}

// Marshal converts a Config into a JSON encoded byte array.
func (c Config) Marshal() ([]byte, error) {
	var buf bytes.Buffer
	if err := c.Write(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Write sends a JSON encoded byte array to a file.
func (c Config) WriteFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return c.Write(file)
}
