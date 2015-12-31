package resp

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Datalogger struct {
	Dataloggers   []string
	Type          string
	Label         string
	Rate          float64
	Frequencey    float64
	StorageFormat string
	ClockDrift    float64
	DLFilters     []string `toml:"filters"`
	Match         string
	Skip          string
}

type Sensor struct {
	Sensors   []string
	SNFilters []string `toml:"filters"`
	Channels  string
	Reversed  bool
}

type Response struct {
	Sensors     []Sensor     `toml:"sensors"`
	Dataloggers []Datalogger `toml:"dataloggers"`
}

type Stream struct {
	Datalogger
	Sensor
}

type responses struct {
	Responses []Response `toml:"response"`
}

type Responses []Response

func LoadResponseFile(path string) ([]Response, error) {
	var resp responses
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if _, err := toml.Decode(string(b), &resp); err != nil {
		return nil, err
	}

	return resp.Responses, nil
}

func LoadResponseFiles(dirname, filename string) ([]Response, error) {

	var resp []Response
	err := filepath.Walk(dirname, func(path string, fi os.FileInfo, err error) error {
		if err == nil && filepath.Base(path) == filename {
			r, err := LoadResponseFile(path)
			if err != nil {
				return err
			}
			resp = append(resp, r...)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func StoreResponseFile(path string, resp []Response) error {

	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(responses{resp}); err != nil {
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
