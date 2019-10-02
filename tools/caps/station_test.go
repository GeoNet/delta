package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestBuilder(t *testing.T) {

	stns := map[string]Station{
		"station_NZ_NBSS": {
			Network: "NZ",
			Code:    "NBSS",
		},
		"station_NZ_CAW": {
			Network: "NZ",
			Code:    "CAW",
		},
	}

	for k, s := range stns {
		t.Run("check "+k, func(t *testing.T) {
			d, err := ioutil.TempDir(os.TempDir(), "test")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(d)

			if err := s.Output(d); err != nil {
				t.Fatalf("unable to store key output %s: %v", d, err)
			}

			all, err := ioutil.ReadFile(filepath.Join(d, k))
			if err != nil {
				t.Fatalf("unable to read temp key file %s: %v", d, err)
			}
			if string(all) != contents {
				t.Errorf("contents mismatch %s", k)
			}
		})
	}
}
