package resp

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	"github.com/BurntSushi/toml"
)

type DataloggerModel struct {
	Model        string // FDSN StationXML Datalogger Model
	Type         string // FDSN StationXML Datalogger Type
	Description  string // FDSN StationXML Datalogger Description
	Manufacturer string // FDSN StationXML Datalogger Manufacturer
	Vendor       string // FDSN StationXML Datalogger Vendor
}

type dataloggerModelList struct {
	Models []DataloggerModel `toml:"datalogger"`
}

type DataloggerModels []DataloggerModel

func (c DataloggerModels) Len() int      { return len(c) }
func (c DataloggerModels) Swap(i, j int) { c[i], c[j] = c[j], c[i] }
func (c DataloggerModels) Less(i, j int) bool {
	switch {
	case c[i].Model < c[j].Model:
		return true
	case c[i].Model > c[j].Model:
		return false
	default:
		return false
	}
}

func LoadDataloggerModelFile(path string) ([]DataloggerModel, error) {
	var dataloggers dataloggerModelList
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if _, err := toml.Decode(string(b), &dataloggers); err != nil {
		return nil, err
	}

	sort.Sort(DataloggerModels(dataloggers.Models))

	return dataloggers.Models, nil
}

func LoadDataloggerModelFiles(dirname, filename string) ([]DataloggerModel, error) {

	var dataloggers []DataloggerModel
	err := filepath.Walk(dirname, func(path string, fi os.FileInfo, err error) error {
		if err == nil && filepath.Base(path) == filename {
			s, err := LoadDataloggerModelFile(path)
			if err != nil {
				return err
			}
			dataloggers = append(dataloggers, s...)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Sort(DataloggerModels(dataloggers))

	return dataloggers, nil
}

func StoreDataloggerModelFile(path string, dataloggers []DataloggerModel) error {

	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(dataloggerModelList{dataloggers}); err != nil {
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
