package resp

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	"github.com/BurntSushi/toml"
)

type Coefficient struct {
	Value float64 `yaml:"value"`
}

type Polynomial struct {
	Name                    string  `yaml:"name"`
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

type polynomialList struct {
	Polynomials []Polynomial `toml:"polynomial"`
}

type Polynomials []Polynomial

func (p Polynomials) Len() int      { return len(p) }
func (p Polynomials) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p Polynomials) Less(i, j int) bool {
	return p[i].Name < p[j].Name
}

func LoadPolynomialFile(path string) ([]Polynomial, error) {
	var pols polynomialList
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if _, err := toml.Decode(string(b), &pols); err != nil {
		return nil, err
	}

	sort.Sort(Polynomials(pols.Polynomials))

	return pols.Polynomials, nil
}

func LoadPolynomialFiles(dirname, filename string) ([]Polynomial, error) {

	var pols []Polynomial
	err := filepath.Walk(dirname, func(path string, fi os.FileInfo, err error) error {
		if err == nil && filepath.Base(path) == filename {
			p, err := LoadPolynomialFile(path)
			if err != nil {
				return err
			}
			pols = append(pols, p...)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Sort(Polynomials(pols))

	return pols, nil
}

func StorePolynomialFile(path string, pols []Polynomial) error {

	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(polynomialList{pols}); err != nil {
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
