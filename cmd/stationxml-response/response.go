package main

import (
	"encoding/json"
	"io"
	"os"

	stationxml "github.com/GeoNet/delta/internal/stationxml/v1.2"
)

// Responses holds a map of Stream slices per installation.
// It also includes lookup maps for the sensor poles and zeros
// entries as well as the decimation fir filters used.
type Responses struct {
	PolesZeros map[string]PolesZeros `json:"poles_zeros,omitempty"`
	FirFilters map[string]FirFilter  `json:"fir_filters,omitempty"`
	Streams    map[string][]Stream   `json:"streams,omitempty"`
}

func (r *Responses) GetPolesZeros(stage stationxml.ResponseStageType) (string, bool) {
	paz, ok := GetPolesZeros(stage)
	if !ok {
		return "", false
	}

	if r.PolesZeros == nil {
		r.PolesZeros = make(map[string]PolesZeros)
	}

	r.PolesZeros[paz.Name] = paz

	return paz.Name, true
}

func (r *Responses) GetFirFilter(stage stationxml.ResponseStageType) (string, bool) {
	fir, ok := GetFirFilter(stage)
	if !ok {
		return "", false
	}

	if r.FirFilters == nil {
		r.FirFilters = make(map[string]FirFilter)
	}

	r.FirFilters[fir.Name] = fir

	return fir.Name, true
}

func (r *Responses) GetStream(srcname string, channel stationxml.ChannelType) (Stream, bool) {

	var sensor string
	if channel.Sensor != nil {
		sensor = channel.Sensor.Model
	}

	var datalogger string
	if channel.DataLogger != nil {
		datalogger = channel.DataLogger.Model
	}

	var paz string
	var filters []DecimationFilter

	for _, s := range channel.Response.Stage {
		if pz, ok := r.GetPolesZeros(s); ok {
			paz = pz
		}
		if _, ok := r.GetFirFilter(s); !ok {
			continue
		}
		if filter, ok := GetDecimationFilter(s); ok {
			filters = append(filters, filter)
		}
	}

	stream := Stream{
		Srcname:     srcname,
		Sensor:      sensor,
		Datalogger:  datalogger,
		SampleRate:  channel.SampleRate.Value,
		Sensitivity: channel.Response.InstrumentSensitivity.Value,
		Frequency:   channel.Response.InstrumentSensitivity.Frequency,
		InputUnits:  channel.Response.InstrumentSensitivity.InputUnits.Name,
		OutputUnits: channel.Response.InstrumentSensitivity.OutputUnits.Name,

		PolesZeros:        paz,
		DecimationFilters: filters,

		StartDate: channel.StartDate.Time,
		EndDate:   channel.EndDate.Time,
	}

	return stream, true
}

func (r *Responses) AddStream(srcname string, channel stationxml.ChannelType) {
	if stream, ok := r.GetStream(srcname, channel); ok {
		if r.Streams == nil {
			r.Streams = make(map[string][]Stream)
		}
		r.Streams[stream.Srcname] = append(r.Streams[stream.Srcname], stream)
	}
}

func (r *Responses) WriteFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return r.Write(file)
}

func (r *Responses) Write(wr io.Writer) error {
	return json.NewEncoder(wr).Encode(r)
}
