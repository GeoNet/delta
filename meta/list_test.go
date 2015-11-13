package meta

import (
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

func TestList(t *testing.T) {

	var liststart, _ = time.Parse(DateTimeFormat, "2010-01-01T12:00:00Z")
	var listend, _ = time.Parse(DateTimeFormat, "2012-01-01T12:00:00Z")

	var listtests = []struct {
		f     string
		l     list
		load  func(string, interface{}) ([]byte, error)
		lists func(string, string, interface{}) ([]byte, error)
	}{
		{
			"testdata/stations.csv",
			Stations{
				Station{
					Network:   "AA",
					Code:      "AAAA",
					Name:      "A Name",
					Latitude:  -41.5,
					Longitude: 173.5,
					StartTime: liststart,
					EndTime:   listend,
				},
				Station{
					Network:   "BB",
					Code:      "BBBB",
					Name:      "B Name",
					Latitude:  -42.5,
					Longitude: 174.5,
					StartTime: liststart.Add(time.Hour),
					EndTime:   listend.Add(time.Hour),
				},
			},
			func(f string, l interface{}) ([]byte, error) {
				err := load(f, l.(*Stations))
				return marshal(*l.(*Stations)), err
			},
			func(d, f string, l interface{}) ([]byte, error) {
				err := lists(d, f, l.(*Stations))
				return marshal(*l.(*Stations)), err
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
			if string(b) != string(marshal(tt.l)) {
				t.Errorf("stations file text mismatch: %s [\n%s\n]", tt.f, bdiff(b, marshal(tt.l)))
			}
		}
		t.Log("Check list file: " + tt.f)
		{
			s, err := tt.load(tt.f, reflect.New(reflect.ValueOf(tt.l).Type()).Interface())
			if err != nil {
				t.Fatal(err)
			}
			if string(marshal(tt.l)) != string(s) {
				t.Errorf("stations file list mismatch: %s [\n%s\n]", tt.f, diff(string(marshal(tt.l)), string(s)))
			}
		}
		t.Log("Check loading files: " + tt.f)
		{
			s, err := tt.lists(filepath.Dir(tt.f), filepath.Base(tt.f), reflect.New(reflect.ValueOf(tt.l).Type()).Interface())
			if err != nil {
				t.Fatal(err)
			}
			if string(marshal(tt.l)) != string(s) {
				t.Errorf("stations file load mismatch: [\n%s\n]", diff(string(marshal(tt.l)), string(s)))
			}
		}
	}
}
