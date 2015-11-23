package meta

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Complex struct {
	Value complex128
}

func (c *Complex) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var params string
	if err := unmarshal(&params); err != nil {
		return err
	}

	var v complex128
	if _, err := fmt.Sscanf(params, "%g", &v); err != nil {
		return err
	}

	//*c = Complex(v)
	c.Value = v

	return nil
}

func (c Complex) MarshalYAML() (interface{}, error) {
	return fmt.Sprintf("%g", complex128(c.Value)), nil
}

func (c Complex) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf("%g", complex128(c.Value))), nil
}

func (c *Complex) UnmarshalText(text []byte) error {
	var v complex128

	if _, err := fmt.Sscanf(string(text), "%g", &v); err != nil {
		return err
	}

	//*c = Complex(v)
	c.Value = v

	return nil
}

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
