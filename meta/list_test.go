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
		file string
		test List
		list func() List
	}{
		{
			"testdata/networks.csv",
			Networks{
				Network{
					Code:       "AA",
					Map:        "XX",
					Name:       "A Name",
					Restricted: false,
				},
				Network{
					Code:       "BB",
					Map:        "XX",
					Name:       "B Name",
					Restricted: true,
				},
			},
			func() List { return &Networks{} },
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
			func() List { return &Stations{} },
		},
	}

	for _, tt := range listtests {
		res := MarshalList(tt.test)

		t.Log("Compare raw list file: " + tt.file)
		{
			b, err := ioutil.ReadFile(tt.file)
			if err != nil {
				t.Fatal(err)
			}
			if string(res) != string(b) {
				t.Errorf("list file text mismatch: %s [\n%s\n]", tt.file, diff(string(res), string(b)))
			}
		}
		t.Log("Check encode/decode list: " + tt.file)
		{
			list := tt.list()
			if err := UnmarshalList(res, list); err != nil {
				t.Fatal(err)
			}
			s := MarshalList(list)

			if string(res) != string(s) {
				t.Errorf("list encode/reencode mismatch: %s [\n%s\n]", tt.file, diff(string(res), string(s)))
			}
		}

		t.Log("Check list file: " + tt.file)
		{
			list := tt.list()
			if err := LoadList(tt.file, list); err != nil {
				t.Fatal(err)
			}
			s := MarshalList(list)

			if string(res) != string(s) {
				t.Errorf("list file list mismatch: %s [\n%s\n]", tt.file, diff(string(res), string(s)))
			}
		}

		t.Log("Check loading files: " + tt.file)
		{
			list := tt.list()
			if err := LoadLists(filepath.Dir(tt.file), filepath.Base(tt.file), list); err != nil {
				t.Fatal(err)
			}

			s := MarshalList(list)
			if string(res) != string(s) {
				t.Errorf("list file load mismatch: [\n%s\n]", diff(string(res), string(s)))
			}
		}
	}
}
