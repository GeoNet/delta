package main

import (
	"fmt"

	"github.com/GeoNet/delta/meta"
	"github.com/GeoNet/delta/resp"
)

// Builder is a cache of response files.
type Builder struct {
	lookup     string
	correction bool
	freqs      Frequencies
	resps      map[string][]byte
}

func NewBuilder(lookup string, correction bool, freqs Frequencies) *Builder {
	return &Builder{
		lookup:     lookup,
		correction: correction,
		freqs:      freqs,
		resps:      make(map[string][]byte),
	}
}

// Lookup returns the response data for a given XML lookup key.
func (b *Builder) Lookup(key string) ([]byte, error) {
	if r, ok := b.resps[key]; ok {
		return r, nil
	}
	data, err := resp.LookupBase(b.lookup, key)
	if err != nil {
		return nil, err
	}
	b.resps[key] = data

	return data, nil
}

// Frequency selects the longest matching response frequency.
func (b *Builder) Frequency(code string) float64 {

	if v, ok := b.freqs.Find(code); ok {
		return v
	}

	return 15.0
}

// Response returns a stationxml ResponseType based on whether the Component has a sampling rate or not.
func (b *Builder) Response(collection meta.Collection, correction meta.Correction) (*resp.ResponseType, error) {
	switch {
	case collection.Component.SamplingRate != 0:
		return b.DerivedResponseType(collection)
	default:
		return b.PairedResponseType(collection, correction)
	}
}

// Prefix is used to attach to internal Id values.
func (b *Builder) Prefix(c meta.Collection) string {
	at := c.InstalledSensor.Start
	if c.DeployedDatalogger.Start.After(at) {
		at = c.DeployedDatalogger.Start
	}
	return fmt.Sprintf("%s.%s.%s.%s.", c.InstalledSensor.Station, c.InstalledSensor.Location, c.Code(), at.Format("2006.002"))
}

// DerivedResponseType return a ResponseType pointer for the derived type.
func (b *Builder) DerivedResponseType(c meta.Collection) (*resp.ResponseType, error) {
	pair := resp.NewInstrumentResponse(
		resp.Prefix(b.Prefix(c)),
		resp.Serial(c.InstalledSensor.Serial),
		resp.Frequency(b.Frequency(c.Code())),
	)

	// find the derived response
	derived, err := b.Lookup(c.Component.Response)
	if err != nil || !(len(derived) > 0) {
		return nil, err
	}

	// generate the derived response
	ans, err := pair.Derived(derived)
	if err != nil {
		return nil, err
	}

	// zero out corrections if not wanted
	if !b.correction {
		for i, s := range ans.Stages {
			if s.Decimation == nil {
				continue
			}
			ans.Stages[i].Decimation.Delay = 0.0
			ans.Stages[i].Decimation.Correction = 0.0
		}
	}

	return ans, nil
}

// PairedResponseType return a ResponseType pointer where both a datalogger and sensor will be described.
func (b *Builder) PairedResponseType(c meta.Collection, v meta.Correction) (*resp.ResponseType, error) {
	pair := resp.NewInstrumentResponse(
		resp.Prefix(b.Prefix(c)),
		resp.Serial(c.InstalledSensor.Serial),
		resp.Frequency(b.Frequency(c.Code())),
	)

	// adjust for corrections
	if v.SensorCalibration != nil {
		pair.SetCalibration(v.SensorCalibration.ScaleFactor, v.SensorCalibration.ScaleBias, v.SensorCalibration.ScaleAbsolute)
	}
	if v.Gain != nil {
		pair.SetGain(v.Gain.Scale.Factor, v.Gain.Scale.Bias, v.Gain.Absolute)
	}
	if v.Telemetry != nil {
		pair.SetTelemetry(v.Telemetry.ScaleFactor)
	}
	if v.Preamp != nil {
		pair.SetPreamp(v.Preamp.ScaleFactor)
	}

	// find the sensor and add it to the pair
	sensor, err := b.Lookup(c.Component.Response)
	if err != nil || !(len(sensor) > 0) {
		return nil, err
	}
	if err := pair.SetSensor(sensor); err != nil {
		return nil, err
	}

	// find the datalogger and add it to the pair
	datalogger, err := b.Lookup(c.Channel.Response)
	if err != nil || !(len(datalogger) > 0) {
		return nil, nil
	}
	if err := pair.SetDatalogger(datalogger); err != nil {
		return nil, err
	}

	// generate the response
	ans, err := pair.ResponseType()
	if err != nil {
		return nil, err
	}

	// zero out corrections if not wanted
	if !b.correction {
		for i, s := range ans.Stages {
			if s.Decimation == nil {
				continue
			}
			ans.Stages[i].Decimation.Delay = 0.0
			ans.Stages[i].Decimation.Correction = 0.0
		}
	}

	return ans, nil
}
