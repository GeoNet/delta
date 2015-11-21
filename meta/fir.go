package meta

import (
	//	"fmt"
	"io/ioutil"

	"github.com/go-yaml/yaml"
)

/*
type Complex complex128

//UnmarshalYAML(unmarshal func(interface{}) error) error

func (c *Complex) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var params string
	if err := unmarshal(&params); err != nil {
		return err
	}

	var v complex128
	if _, err := fmt.Sscanf(params, "%g", &v); err != nil {
		return err
	}

	*c = Complex(v)

	return nil
}

func (c Complex) MarshalYAML() (interface{}, error) {
	return fmt.Sprintf("%g", complex128(c)), nil
}
*/

type FIR struct {
	Name       string    `yaml:"name"`
	Causal     bool      `yaml:"causal"`
	Symmetry   string    `yaml:"symmetry"`
	Decimation float64   `yaml:"decimation"`
	Gain       float64   `yaml:gain"`
	Notes      *string   `yaml:"notes,omitempty"`
	Factors    []float64 `yaml:"factors,omitempty"`
	/*
		Code  string    `yaml:"code"`
		Type  string    `yaml:"type"`
		Poles []Complex `yaml:"poles,omitempty"`
		Zeros []Complex `yaml:"zeros,omitempty"`
	*/
}

func LoadFIR(file string) (map[string]FIR, error) {
	p := make(map[string]FIR)

	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(b, p)
	if err != nil {
		return nil, err
	}

	return p, nil
}
