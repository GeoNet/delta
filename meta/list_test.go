package meta

import (
	"io/ioutil"
	"path/filepath"
	"testing"
	"time"
)

func TestList(t *testing.T) {

	var liststart, _ = time.Parse(DateTimeFormat, "2010-01-01T12:00:00Z")
	var listend, _ = time.Parse(DateTimeFormat, "2012-01-01T12:00:00Z")

	var listtests = []struct {
		f     string
		l     list
		load  func(string) ([]byte, error)
		lists func(string, string) ([]byte, error)
	}{
		{
			"testdata/networks.csv",
			Networks{
				Network{
					Code: "AA",
					Map:  "XX",
					Name: "A Network",
				},
				Network{
					Code:       "BB",
					Map:        "XX",
					Name:       "B Network",
					Restricted: true,
				},
			},
			func(f string) ([]byte, error) {
				var n Networks
				err := load(f, &n)
				return marshal(n), err
			},
			func(d, f string) ([]byte, error) {
				var n Networks
				err := lists(d, f, &n)
				return marshal(n), err
			},
		},
		{
			"testdata/stations.csv",
			Stations{
				Station{
					Code:      "AAAA",
					Network:   "AA",
					Name:      "A Name",
					Latitude:  -41.5,
					Longitude: 173.5,
					StartTime: liststart,
					EndTime:   listend,
				},
				Station{
					Code:      "BBBB",
					Network:   "BB",
					Name:      "B Name",
					Latitude:  -42.5,
					Longitude: 174.5,
					StartTime: liststart.Add(time.Hour),
					EndTime:   listend.Add(time.Hour),
				},
			},
			func(f string) ([]byte, error) {
				var s Stations
				err := load(f, &s)
				return marshal(s), err
			},
			func(d, f string) ([]byte, error) {
				var s Stations
				err := lists(d, f, &s)
				return marshal(s), err
			},
		},
	}

	for _, tt := range listtests {
		t.Log("Compare list file: " + tt.f)
		{
			b, err := ioutil.ReadFile(tt.f)
			if err != nil {
				t.Fatal(err)
			}
			m := marshal(tt.l)
			if string(m) != string(b) {
				t.Errorf("stations file text mismatch: %s [\n%s\n]", tt.f, diff(string(m), string(b)))
			}
		}
		t.Log("Check list file: " + tt.f)
		{
			s, err := tt.load(tt.f)
			if err != nil {
				t.Fatal(err)
			}
			m := marshal(tt.l)
			if string(m) != string(s) {
				t.Errorf("stations file list mismatch: %s [\n%s\n]", tt.f, diff(string(m), string(s)))
			}
		}

		t.Log("Check loading files: " + tt.f)
		{
			s, err := tt.lists(filepath.Dir(tt.f), filepath.Base(tt.f))
			if err != nil {
				t.Fatal(err)
			}
			m := marshal(tt.l)
			if string(m) != string(s) {
				t.Errorf("stations file load mismatch: [\n%s\n]", diff(string(m), string(s)))
			}
		}
	}
}
