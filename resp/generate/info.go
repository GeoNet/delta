package main

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/mitchellh/hashstructure"
)

var header = `
package resp

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

`

func hash(v interface{}) string {
	i, err := hashstructure.Hash(v, nil)
	if err != nil {
		panic(err)
	}
	return strconv.FormatUint(i, 16)
}

func join(s string) string {
	s = strings.Split(s, "/")[0]
	s = strings.Title(strings.ToLower(s))
	s = strings.Join(strings.Fields(s), "-")
	for _, k := range []string{"/", "#"} {
		s = strings.Replace(s, k, "-", -1)
	}
	return s
}

// yaml is unable to handle complex numbers
type Complex64 complex64

func (c *Complex64) UnmarshalText(text []byte) error {
	_, err := fmt.Sscanf(string(text), "%v", c)
	return err
}

func (c Complex64) MarshalText() ([]byte, error) {
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
	Response        map[string]Response        `yaml:"response"`
}

func NewResponseInfo() *ResponseInfo {
	return &ResponseInfo{
		PAZ:             make(map[string]PAZ),
		Polynomial:      make(map[string]Polynomial),
		FIR:             make(map[string]FIR),
		DataloggerModel: make(map[string]DataloggerModel),
		SensorModel:     make(map[string]SensorModel),
		Filter:          make(map[string][]ResponseStage),
		Response:        make(map[string]Response),
	}
}

func (r *ResponseInfo) Merge(resp ResponseInfo) {
	for k, v := range resp.PAZ {
		v.ResourceId = fmt.Sprintf("smi:geonet.org.nz/ResponsePAZ#%s", hash(v))
		r.PAZ[k] = v
	}
	for k, v := range resp.Polynomial {
		v.ResourceId = fmt.Sprintf("smi:geonet.org.nz/ResponsePolynomial#%s", hash(v))
		r.Polynomial[k] = v
	}
	for k, v := range resp.FIR {
		v.ResourceId = fmt.Sprintf("smi:geonet.org.nz/ResponseFIR#%s", hash(v))
		r.FIR[k] = v
	}
	for k, v := range resp.DataloggerModel {
		v.ResourceId = fmt.Sprintf("smi:geonet.org.nz/Datalogger#%s-%s", hash(v), join(k))
		r.DataloggerModel[k] = v
	}
	for k, v := range resp.SensorModel {
		v.ResourceId = fmt.Sprintf("smi:geonet.org.nz/Sensor#%s-%s", hash(v), join(k))
		r.SensorModel[k] = v
	}
	for k, v := range resp.Filter {
		r.Filter[k] = v
	}
	for k, v := range resp.Response {
		r.Response[k] = v
	}
}

func (r ResponseInfo) Generate(w io.Writer) error {
	g := Generate{
		ResponseMap: responseMap(r.Response),
		FilterMap:   filterMap(r.Filter),
		PazMap:      pazMap(r.PAZ),
		FirMap:      firMap(r.FIR),
		PolyMap:     polynomialMap(r.Polynomial),
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
