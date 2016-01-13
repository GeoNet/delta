package resp

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	"github.com/BurntSushi/toml"
)

type ResponseStage struct {
	Type        string
	Lookup      string
	Frequency   float64
	SampleRate  float64
	Factor      int32
	Gain        float64
	Scale       float64
	Correction  float64
	Delay       float64
	InputUnits  string
	OutputUnits string
}

type Filter struct {
	Name   string
	Stages []ResponseStage `toml:"stage"`
}

type Filters []Filter

func (f Filters) Len() int      { return len(f) }
func (f Filters) Swap(i, j int) { f[i], f[j] = f[j], f[i] }
func (f Filters) Less(i, j int) bool {
	switch {
	case f[i].Name < f[j].Name:
		return true
	case f[i].Name > f[j].Name:
		return false
	default:
		return false
	}
}

type filterList struct {
	Filters []Filter `toml:"filter"`
}

func LoadFilterFile(path string) ([]Filter, error) {
	var filts filterList
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if _, err := toml.Decode(string(b), &filts); err != nil {
		return nil, err
	}

	sort.Sort(Filters(filts.Filters))

	return filts.Filters, nil
}

func LoadFilterFiles(dirname, filename string) ([]Filter, error) {

	var filts []Filter
	err := filepath.Walk(dirname, func(path string, fi os.FileInfo, err error) error {
		if err == nil && filepath.Base(path) == filename {
			f, err := LoadFilterFile(path)
			if err != nil {
				return err
			}
			filts = append(filts, f...)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Sort(Filters(filts))

	return filts, nil
}

func StoreFilterFile(path string, filts []Filter) error {

	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(filterList{filts}); err != nil {
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
