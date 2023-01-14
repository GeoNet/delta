package main

import (
	"fmt"

	"github.com/GeoNet/delta/meta"
	"github.com/GeoNet/delta/resp"

	"github.com/GeoNet/delta/internal/stationxml"
)

// ResponseFinder describes how to find a response segment given a key name.
type ResponseFinder interface {
	Lookup(string) ([]byte, error)
}

// LookupFunc describes a function that can find a response segment given a key name.
type LookupFunc func(key string) ([]byte, error)

// Lookup implements ResponseFinder
func (l LookupFunc) Lookup(key string) ([]byte, error) {
	return l(key)
}

// Lookup uses the in built delta response finding mechanism.
func Lookup(base string) LookupFunc {
	return func(key string) ([]byte, error) {
		return resp.LookupBase(base, key)
	}
}

// Builder is used to generate StationXML files.
type Builder struct {
	finders []ResponseFinder
	cache   map[string][]byte
}

// NewBuilder returns a Builder pointer with the given set of ResponseFinder elements.
func NewBuilder(finder ...ResponseFinder) *Builder {
	return &Builder{
		finders: finder,
		cache:   make(map[string][]byte),
	}
}

// Lookup returns the first response returned by the Builder ResponseFinder elements, or nil if no response could be found.
func (b *Builder) Lookup(key string) ([]byte, error) {
	if r, ok := b.cache[key]; ok {
		return r, nil
	}
	for _, f := range b.finders {
		data, err := f.Lookup(key)
		if err != nil {
			return nil, err
		}
		if data == nil {
			continue
		}
		b.cache[key] = data
		return data, nil
	}
	return nil, nil
}

// Response returns a stationxml ResponseType based on whether the Component has a sampling rate or not.
func (b *Builder) Response(collection meta.Collection, correction meta.Correction) (*stationxml.ResponseType, error) {
	switch {
	case collection.Component.SamplingRate != 0:
		return b.DerivedResponseType(collection)
	default:
		return b.PairedResponseType(collection, correction)
	}
}

// Prefix is used when building the XML identifier tags.
func (b *Builder) Prefix(c meta.Collection) string {
	at := c.InstalledSensor.Start
	if c.DeployedDatalogger.Start.After(at) {
		at = c.DeployedDatalogger.Start
	}
	return fmt.Sprintf("%s.%s.%s.%s.", c.InstalledSensor.Station, c.InstalledSensor.Location, c.Code(), at.Format("2006.002"))
}

// DerivedResponseType builds a general ResponseType from a response without a sensor which can then be encoded into the required version.
func (b *Builder) DerivedResponseType(c meta.Collection) (*stationxml.ResponseType, error) {
	pair := stationxml.NewResponse(
		stationxml.Prefix(b.Prefix(c)),
		stationxml.Serial(c.InstalledSensor.Serial),
		stationxml.Frequency(ChannelFrequency(c.Code())),
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

	return ans, nil
}

// PairedResponseType builds a general ResponseType for a paired sensor and datalogger which can then be encoded into the required version.
func (b *Builder) PairedResponseType(c meta.Collection, v meta.Correction) (*stationxml.ResponseType, error) {
	pair := stationxml.NewResponse(
		stationxml.Prefix(b.Prefix(c)),
		stationxml.Serial(c.InstalledSensor.Serial),
		stationxml.Frequency(ChannelFrequency(c.Code())),
	)

	// adjust for corrections
	if v.SensorCalibration != nil {
		pair.SetCalibration(v.SensorCalibration.ScaleFactor, v.SensorCalibration.ScaleBias)
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

	return ans, nil
}
