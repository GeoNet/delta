package main

import (
	"io/ioutil"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBuild(t *testing.T) {

	// load in the test data and convert to stationxml indented text
	b1, err := ioutil.ReadFile("./testdata/amplitude.xml")
	if err != nil {
		t.Fatalf("error: unable to load test amplitudes file: %v", err)
	}

	cfgs, err := loadConfig("./testdata/chart-amplitude.yaml")
	if err != nil {
		t.Fatalf("error: unable to load test config file: %v", err)
	}

	amplitudes, err := buildAmplitudes(cfgs, "./testdata", "/work/chart/amplitude")
	if err != nil {
		t.Fatalf("error: unable to build test amplitudes file: %v", err)
	}

	b2, err := encodeAmplitudes(amplitudes)
	if err != nil {
		t.Fatalf("error: unable to encode test amplitudes file: %v", err)
	}

	// compare stored with computed
	if string(b1) != string(b2) {
		t.Errorf("unexpected content -got/+exp\n%s", cmp.Diff(string(b1), string(b2)))
	}
}
