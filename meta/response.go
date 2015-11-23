package meta

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	"github.com/BurntSushi/toml"
)

/*
type Response struct {
	Datalogger string `csv:"Datalogger Model"`
	Sensor     string `csv:"Sensor Model"`
	Reversed   bool   `csv:"Reversed Connection"`
	Lookup     string `csv:"Response Lookup"`
	Match      string `csv:"Site Match"`
}

type Responses []Response

func (r Responses) Len() int      { return len(r) }
func (r Responses) Swap(i, j int) { r[i], r[j] = r[j], r[i] }
func (r Responses) Less(i, j int) bool {
	switch {
	case r[i].Datalogger < r[j].Datalogger:
		return true
	case r[i].Datalogger > r[j].Datalogger:
		return false
	case r[i].Sensor < r[j].Sensor:
		return true
	case r[i].Sensor > r[j].Sensor:
		return false
	case r[i].Match < r[j].Match:
		return true
	case r[i].Match > r[j].Match:
		return false
	default:
		return false
	}
}

func (r Responses) List()      {}
func (r Responses) Sort() List { sort.Sort(r); return r }
*/

type Stream struct {
	Type          string
	Label         string
	Channels      string
	Rate          float64
	Frequencey    float64
	StorageFormat string
	ClockDrift    float64
	Stages        string
}

type Response struct {
	Dataloggers []string
	Sensors     []string
	Reversed    bool
	Lookup      string
	Match       string

	Streams []Stream `toml:"streams"`
}

type responses struct {
	Responses []Response `toml:"response"`
}

type Responses []Response

func (r Responses) Len() int      { return len(r) }
func (r Responses) Swap(i, j int) { r[i], r[j] = r[j], r[i] }
func (r Responses) Less(i, j int) bool {
	switch {
	case r[i].Lookup < r[j].Lookup:
		return true
	case r[i].Lookup > r[j].Lookup:
		return false
	case r[i].Match < r[j].Match:
		return true
	case r[i].Match > r[j].Match:
		return false
	default:
		return false
	}
}

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

	sort.Sort(Responses(resp))

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
