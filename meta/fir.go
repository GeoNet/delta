package meta

import (
//"io/ioutil"

//	"gopkg.in/yaml.v2"
)

type FIR struct {
	Name       string    `yaml:"name"`
	Causal     bool      `yaml:"causal"`
	Symmetry   string    `yaml:"symmetry"`
	Decimation float64   `yaml:"decimation"`
	Gain       float64   `yaml:gain"`
	Notes      *string   `yaml:"notes,omitempty"`
	Factors    []float64 `yaml:"factors,omitempty"`
}

/*
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
*/
