package resp

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	"github.com/BurntSushi/toml"
)

type SensorComponent struct {
	Azimuth float64
	Dip     float64
}

type SensorModel struct {
	Model        string // FDSN StationXML Sensor Model
	Type         string // FDSN StationXML Sensor Type
	Description  string // FDSN StationXML Sensor Description
	Manufacturer string // FDSN StationXML Vendor Description
	Vendor       string // FDSN StationXML Vendor Description

	Pins []SensorComponent `toml:"component"`
}

type sensorModelList struct {
	Models []SensorModel `toml:"sensor"`
}

type SensorModels []SensorModel

func (c SensorModels) Len() int      { return len(c) }
func (c SensorModels) Swap(i, j int) { c[i], c[j] = c[j], c[i] }
func (c SensorModels) Less(i, j int) bool {
	switch {
	case c[i].Model < c[j].Model:
		return true
	case c[i].Model > c[j].Model:
		return false
	default:
		return false
	}
}

func LoadSensorModelFile(path string) ([]SensorModel, error) {
	var sensors sensorModelList
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if _, err := toml.Decode(string(b), &sensors); err != nil {
		return nil, err
	}

	sort.Sort(SensorModels(sensors.Models))

	return sensors.Models, nil
}

func LoadSensorModelFiles(dirname, filename string) ([]SensorModel, error) {

	var sensors []SensorModel
	err := filepath.Walk(dirname, func(path string, fi os.FileInfo, err error) error {
		if err == nil && filepath.Base(path) == filename {
			s, err := LoadSensorModelFile(path)
			if err != nil {
				return err
			}
			sensors = append(sensors, s...)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Sort(SensorModels(sensors))

	return sensors, nil
}

func StoreSensorModelFile(path string, sensors []SensorModel) error {

	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(sensorModelList{sensors}); err != nil {
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
