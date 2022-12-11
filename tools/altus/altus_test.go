package main

import (
	"io/ioutil"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBuild(t *testing.T) {

	// load in the test data and convert to stationxml indented text
	b1, err := ioutil.ReadFile("./testdata/altus.xml")
	if err != nil {
		t.Fatalf("error: unable to load test altus file: %v", err)
	}

	sites, err := buildSites("./testdata")
	if err != nil {
		t.Fatalf("error: unable to build test altus file: %v", err)
	}

	b2, err := encodeSites(sites)
	if err != nil {
		t.Fatalf("error: unable to encode test altus file: %v", err)
	}

	// compare stored with computed
	if string(b1) != string(b2) {
		t.Errorf("unexpected content -got/+exp\n%s", cmp.Diff(string(b1), string(b2)))
	}
}
