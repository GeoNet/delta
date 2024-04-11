package main

import (
	"encoding/xml"

	"github.com/GeoNet/delta/internal/stationxml"
	"github.com/GeoNet/delta/meta"
	"github.com/GeoNet/delta/resp"
)

// ResponseInfo is used to store decoded stream conversion details.
type ResponseInfo struct {
	Sensitivity float64
	Gain        float64
	Bias        float64
	Input       string
	Output      string
}

// Response builds a ResponseInfor for the given response base, a possible sampling rate, and the component and channel codes.
// If the sampling rate is given then only the component response information is used as this is likely a derived response.
func Response(base string, collection meta.Collection) (*ResponseInfo, error) {

	// simple response
	info := ResponseInfo{
		Sensitivity: 1.0,
		Gain:        1.0,
	}

	// check both component and channel responses, other than if the rate is given.
	for _, lookup := range []string{collection.Component.Response, collection.Channel.Response} {
		if collection.Component.SamplingRate != 0.0 && lookup == collection.Channel.Response {
			continue
		}

		// find the associated StationXML snippet, skip if missing
		data, err := resp.LookupBase(base, lookup)
		if err != nil || len(data) == 0 {
			continue
		}

		// decode the response into a simple form.
		var res stationxml.ResponseType
		if err := xml.Unmarshal(data, &res); err != nil {
			return nil, err
		}

		// extract the response details into the response info.
		switch {
		case res.InstrumentSensitivity != nil:
			if info.Input == "" {
				info.Input = res.InstrumentSensitivity.InputUnits.Name
			}
			info.Output = res.InstrumentSensitivity.OutputUnits.Name
			info.Sensitivity *= res.InstrumentSensitivity.Value
		case res.InstrumentPolynomial != nil:
			if info.Input == "" {
				info.Input = res.InstrumentPolynomial.InputUnits.Name
			}
			info.Output = res.InstrumentPolynomial.OutputUnits.Name
			if coefs := res.PolynomialCoefficients(); len(coefs) == 2 {
				info.Bias = coefs[0].Value
				info.Gain = coefs[1].Value
			}
		}
	}

	return &info, nil
}
