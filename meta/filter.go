package meta

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	"github.com/BurntSushi/toml"
)

type Stage struct {
	Type   string
	Lookup string
	Freq   float64
	Gain   float64
	Scale  float64
	Input  string
	Output string
}

type Filter struct {
	Name   string
	Stages []Stage `toml:"stage"`
}

type filters struct {
	Filters []Filter `toml:"filter"`
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

func LoadFilterFile(path string) ([]Filter, error) {
	var filts filters
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if _, err := toml.Decode(string(b), &filts); err != nil {
		return nil, err
	}

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

	return filts, nil
}

func StoreFilterFile(path string, filts []Filter) error {

	sort.Sort(Filters(filts))

	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(filters{filts}); err != nil {
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
