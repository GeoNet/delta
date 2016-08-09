package main

import (
	"fmt"
	"io"
)

var header = `
package main

/*
 *  WARNING: CODE GENERATED AUTOMATICALLY.
 *
 *  To update: edit or add yaml file(s) in the responses directory.
 *  Commit these changes and run "go generate" in the main project
 *  directory. Changes to this file should then also be commited.
 *
 *  THIS FILE SHOULD NOT BE EDITED BY HAND.
 *
 */

import(
	"github.com/ozym/fdsn/stationxml"
)

`

// yaml is unable to handle complex numbers
type Complex128 complex128

func (c *Complex128) UnmarshalText(text []byte) error {
	_, err := fmt.Sscanf(string(text), "%v", c)
	return err
}

func (c Complex128) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf("%v", c)), nil
}

// manage response information
type ResponseInfo struct {
	PAZ             map[string]PAZ             `yaml:"paz"`
	Polynomial      map[string]Polynomial      `yaml:"polynomial"`
	FIR             map[string]FIR             `yaml:"fir"`
	DataloggerModel map[string]DataloggerModel `yaml:"dataloggermodel"`
	SensorModel     map[string]SensorModel     `yaml:"sensormodel"`
	Filter          map[string][]ResponseStage `yaml:"filter"`
	Response        []Response                 `yaml:"response"`
}

func NewResponseInfo() *ResponseInfo {
	return &ResponseInfo{
		PAZ:             make(map[string]PAZ),
		Polynomial:      make(map[string]Polynomial),
		FIR:             make(map[string]FIR),
		DataloggerModel: make(map[string]DataloggerModel),
		SensorModel:     make(map[string]SensorModel),
		Filter:          make(map[string][]ResponseStage),
	}
}

func (r *ResponseInfo) Merge(resp ResponseInfo) {
	for k, v := range resp.PAZ {
		r.PAZ[k] = v
	}
	for k, v := range resp.Polynomial {
		r.Polynomial[k] = v
	}
	for k, v := range resp.FIR {
		r.FIR[k] = v
	}
	for k, v := range resp.DataloggerModel {
		r.DataloggerModel[k] = v
	}
	for k, v := range resp.SensorModel {
		r.SensorModel[k] = v
	}
	for k, v := range resp.Filter {
		r.Filter[k] = v
	}
	r.Response = append(r.Response, resp.Response...)
}

func (r ResponseInfo) Generate(w io.Writer) error {
	g := Generate{
		ResponseList: r.Response,
		FilterMap:    filterMap(r.Filter),
		PazMap:       pazMap(r.PAZ),
		FirMap:       firMap(r.FIR),
		PolyMap:      polynomialMap(r.Polynomial),
	}
	if _, err := w.Write([]byte(header)); err != nil {
		return err
	}

	if err := dataloggermodel(w, r.DataloggerModel); err != nil {
		return err
	}
	if err := sensormodel(w, r.SensorModel); err != nil {
		return err
	}
	if err := g.generate(w); err != nil {
		return err
	}

	return nil
}
