package main

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBuild(t *testing.T) {

	// load in the test data and convert to stationxml indented text
	b1, err := os.ReadFile("./testdata/spectra.xml")
	if err != nil {
		t.Fatalf("error: unable to load test spectras file: %v", err)
	}

	cfgs, err := loadConfig("./testdata/chart-spectra.yaml")
	if err != nil {
		t.Fatalf("error: unable to load test config file: %v", err)
	}

	spectras, err := buildSpectras(cfgs, "./testdata", "/work/chart/spectra")
	if err != nil {
		t.Fatalf("error: unable to build test spectras file: %v", err)
	}

	b2, err := encodeSpectras(spectras)
	if err != nil {
		t.Fatalf("error: unable to encode test spectras file: %v", err)
	}

	// compare stored with computed
	if string(b1) != string(b2) {
		t.Errorf("unexpected content -got/+exp\n%s", cmp.Diff(string(b1), string(b2)))
	}
}
