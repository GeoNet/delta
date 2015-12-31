package meta

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	"github.com/BurntSushi/toml"
)

type Complex struct {
	Value complex128
}

/*
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
*/

func (c Complex) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf("%g", complex128(c.Value))), nil
}

func (c *Complex) UnmarshalText(text []byte) error {
	var v complex128

	if _, err := fmt.Sscanf(string(text), "%g", &v); err != nil {
		return err
	}

	c.Value = v

	return nil
}

type PAZ struct {
	Name  string
	Code  string
	Type  string
	Notes string
	Poles []Complex
	Zeros []Complex
}

type pazs struct {
	PAZs []PAZ `toml:"paz"`
}

type PAZs []PAZ

func (p PAZs) Len() int      { return len(p) }
func (p PAZs) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p PAZs) Less(i, j int) bool {
	return p[i].Name < p[j].Name
}

func LoadPAZFile(path string) ([]PAZ, error) {
	var pazs pazs
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if _, err := toml.Decode(string(b), &pazs); err != nil {
		return nil, err
	}

	return pazs.PAZs, nil
}

func LoadPAZFiles(dirname, filename string) ([]PAZ, error) {

	var filters []PAZ
	err := filepath.Walk(dirname, func(path string, fi os.FileInfo, err error) error {
		if err == nil && filepath.Base(path) == filename {
			f, err := LoadPAZFile(path)
			if err != nil {
				return err
			}
			filters = append(filters, f...)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return filters, nil
}

func StorePAZFile(path string, filters []PAZ) error {

	sort.Sort(PAZs(filters))

	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(pazs{filters}); err != nil {
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
