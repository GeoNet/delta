package meta

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	"github.com/BurntSushi/toml"
)

type Pin struct {
	Pin     int32
	Azimuth float64
	Dip     float64
}

type Model struct {
	Model string
	Pins  []Pin `toml:"pin"`
}

type models struct {
	Models []Model `toml:"model"`
}

type Models []Model

func (c Models) Len() int      { return len(c) }
func (c Models) Swap(i, j int) { c[i], c[j] = c[j], c[i] }
func (c Models) Less(i, j int) bool {
	switch {
	case c[i].Model < c[j].Model:
		return true
	case c[i].Model > c[j].Model:
		return false
	default:
		return false
	}
}

func LoadModelFile(path string) ([]Model, error) {
	var mods models
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if _, err := toml.Decode(string(b), &mods); err != nil {
		return nil, err
	}

	return mods.Models, nil
}

func LoadModelFiles(dirname, filename string) ([]Model, error) {

	var mods []Model
	err := filepath.Walk(dirname, func(path string, fi os.FileInfo, err error) error {
		if err == nil && filepath.Base(path) == filename {
			r, err := LoadModelFile(path)
			if err != nil {
				return err
			}
			mods = append(mods, r...)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return mods, nil
}

func StoreModelFile(path string, mods []Model) error {

	sort.Sort(Models(mods))

	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(models{mods}); err != nil {
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
