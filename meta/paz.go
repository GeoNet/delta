package meta

import (
	"fmt"
	"io/ioutil"

	"github.com/go-yaml/yaml"
)

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

/*
var params struct {
	        SkipHeaderValidation bool `yaml:"skip-header-validation"`
		    }
*/

func (c Complex) MarshalYAML() (interface{}, error) {
	return fmt.Sprintf("%g", complex128(c)), nil
}

/*
func (c Complex) UnarshalYAML() (interface{}, error) {
	return fmt.Sprintf("%g + %gi", real(c), imag(c)), nil
}
*/

type PAZ struct {
	Name  string    `yaml:"name"`
	Code  string    `yaml:"code"`
	Type  string    `yaml:"type"`
	Notes *string   `yaml:"notes,omitempty"`
	Poles []Complex `yaml:"poles,omitempty"`
	Zeros []Complex `yaml:"zeros,omitempty"`
}

func LoadPAZ(file string) (map[string]PAZ, error) {
	p := make(map[string]PAZ)

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
