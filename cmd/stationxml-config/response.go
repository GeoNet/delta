package main

import (
	"fmt"

	"github.com/GeoNet/delta/internal/build/v1.2"
	"github.com/GeoNet/delta/internal/stationxml/v1.2"

	"github.com/GeoNet/delta/meta"
	"github.com/GeoNet/delta/resp"
)

type Responses struct {
	resps map[string][]byte
}

func NewResponses(lookup string, keys ...string) (*Responses, error) {
	resps := make(map[string][]byte)
	for _, k := range keys {
		data, err := resp.LookupBase(lookup, k)
		if err != nil {
			return nil, err
		}
		if data == nil {
			continue
		}
		resps[k] = data
	}
	responses := Responses{
		resps: resps,
	}

	return &responses, nil
}

func (r *Responses) Prefix(c meta.Collection) string {
	at := c.InstalledSensor.Start
	if c.DeployedDatalogger.Start.After(at) {
		at = c.DeployedDatalogger.Start
	}
	if c.Connection != nil && c.Connection.Start.After(at) {
		at = c.Connection.Start
	}

	return fmt.Sprintf("%s.%s.%s.%s.", c.InstalledSensor.Station, c.InstalledSensor.Location, c.Code(), at.Format("2006.002"))
}

func (r *Responses) DerivedResponseType(c meta.Collection) (*stationxml.ResponseType, error) {
	resp := build.NewResponse(r.Prefix(c), c.InstalledSensor.Serial, LegacyFrequency(c.Code()))

	derived, ok := r.resps[c.Component.Response]
	if !ok || !(len(derived) > 0) {
		return nil, nil
	}

	return resp.Derived(derived)
}

func (r *Responses) PairedResponseType(c meta.Collection, v meta.Correction) (*stationxml.ResponseType, error) {
	resp := build.NewResponse(r.Prefix(c), c.InstalledSensor.Serial, LegacyFrequency(c.Code()))

	if c.InstalledSensor.Station == "KAVZ" {
		resp = build.NewResponse(r.Prefix(c), c.InstalledSensor.Serial, 1.0)
	}

	var gain, bias float64
	if v.Gain != nil {
		gain, bias = v.Gain.Factor, v.Gain.Bias
	}

	var preamp float64
	if v.Preamp != nil {
		preamp = v.Preamp.Gain
	}

	sensor, ok := r.resps[c.Component.Response]
	if !ok || !(len(sensor) > 0) {
		return nil, nil
	}
	if err := resp.Sensor(gain, bias, sensor); err != nil {
		return nil, err
	}
	datalogger, ok := r.resps[c.Channel.Response]
	if !ok || !(len(datalogger) > 0) {
		return nil, nil
	}
	if err := resp.Datalogger(preamp, datalogger); err != nil {
		return nil, err
	}

	return resp.ResponseType()
}

func (r *Responses) Response(c meta.Collection, v meta.Correction) (*stationxml.ResponseType, error) {
	switch {
	case c.Component.SamplingRate != 0:
		return r.DerivedResponseType(c)
	default:
		return r.PairedResponseType(c, v)
	}
}
