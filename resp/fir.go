package resp

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	"github.com/BurntSushi/toml"
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

type firList struct {
	Filters []FIR `toml:"fir"`
}

type FIRs []FIR

func (p FIRs) Len() int      { return len(p) }
func (p FIRs) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p FIRs) Less(i, j int) bool {
	return p[i].Name < p[j].Name
}

func LoadFIRFile(path string) ([]FIR, error) {
	var firs firList
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if _, err := toml.Decode(string(b), &firs); err != nil {
		return nil, err
	}

	sort.Sort(FIRs(firs.Filters))

	return firs.Filters, nil
}

func LoadFIRFiles(dirname, filename string) ([]FIR, error) {

	var filters []FIR
	err := filepath.Walk(dirname, func(path string, fi os.FileInfo, err error) error {
		if err == nil && filepath.Base(path) == filename {
			f, err := LoadFIRFile(path)
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

	sort.Sort(FIRs(filters))

	return filters, nil
}

func StoreFIRFile(path string, filters []FIR) error {

	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(firList{filters}); err != nil {
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
