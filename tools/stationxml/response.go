package main

import (
//	"bytes"
//	"io/ioutil"
//	"os"
//	"path/filepath"
//
//	"github.com/BurntSushi/toml"
)

type Datalogger struct {
	Dataloggers   []string
	Type          string
	Label         string
	Rate          float64
	Frequency     float64
	StorageFormat string
	ClockDrift    float64
	Filters       []string `toml:"filters"`
	Match         string
	Skip          string
}

type DataloggerModel struct {
	Type         string // FDSN StationXML Datalogger Type
	Description  string // FDSN StationXML Datalogger Description
	Manufacturer string // FDSN StationXML Datalogger Manufacturer
	Vendor       string // FDSN StationXML Datalogger Vendor
}

type Sensor struct {
	Sensors  []string
	Filters  []string `toml:"filters"`
	Channels string
	Reversed bool
	Match    string
	Skip     string
}

type SensorComponent struct {
	Azimuth float64
	Dip     float64
}

type SensorModel struct {
	Type         string // FDSN StationXML Sensor Type
	Description  string // FDSN StationXML Sensor Description
	Manufacturer string // FDSN StationXML Vendor Description
	Vendor       string // FDSN StationXML Vendor Description

	Components []SensorComponent `toml:"component"`
}

type Response struct {
	Sensors     []Sensor     `toml:"sensors"`
	Dataloggers []Datalogger `toml:"dataloggers"`
}

type Stream struct {
	Datalogger
	Sensor
}

type ResponseStage struct {
	Type        string
	Lookup      string
	Frequency   float64
	SampleRate  float64
	Decimate    int32
	Gain        float64
	Scale       float64
	Correction  float64
	Delay       float64
	InputUnits  string
	OutputUnits string
}

type Filter struct {
	Stages []ResponseStage `toml:"stage"`
}

type PAZ struct {
	Code  string
	Type  string
	Notes string
	Poles []complex128
	Zeros []complex128
}

type FIR struct {
	Causal     bool      `yaml:"causal"`
	Symmetry   string    `yaml:"symmetry"`
	Decimation float64   `yaml:"decimation"`
	Gain       float64   `yaml:gain"`
	Notes      *string   `yaml:"notes,omitempty"`
	Factors    []float64 `yaml:"factors,omitempty"`
}

type Coefficient struct {
	Value float64 `yaml:"value"`
}

type Polynomial struct {
	Gain                    float64 `yaml:"gain"`
	ApproximationType       string  `yaml:"approximation_type"`
	FrequencyLowerBound     float64 `yaml:"frequency_lower_bound"`
	FrequencyUpperBound     float64 `yaml:"frequency_upper_bound"`
	ApproximationLowerBound float64 `yaml:"approximation_lower_bound"`
	ApproximationUpperBound float64 `yaml:"approximation_upper_bound"`
	MaximumError            float64 `yaml:"maximum_error"`
	Notes                   *string `yaml:"notes,omitempty"`

	Coefficients []Coefficient `yaml:"coefficients,omitempty" toml:"coefficient"`
}

/*
type responseList struct {
	Responses []Response `toml:"response"`
}

type Responses []Response

func LoadResponseFile(path string) ([]Response, error) {
	var resp responseList
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if _, err := toml.Decode(string(b), &resp); err != nil {
		return nil, err
	}

	return resp.Responses, nil
}

func LoadResponseFiles(dirname, filename string) ([]Response, error) {

	var resp []Response
	err := filepath.Walk(dirname, func(path string, fi os.FileInfo, err error) error {
		if err == nil && filepath.Base(path) == filename {
			r, err := LoadResponseFile(path)
			if err != nil {
				return err
			}
			resp = append(resp, r...)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func StoreResponseFile(path string, resp []Response) error {

	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(responseList{resp}); err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(buf.Bytes())
	if err != nil {
		return err
	}

	return nil
}
*/
