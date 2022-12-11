package main

import (
	"encoding/json"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBuild_Network(t *testing.T) {

	// load in the test data and convert to stationxml indented text
	b1, err := ioutil.ReadFile("./testdata/impact.json")
	if err != nil {
		t.Fatalf("error: unable to load test amplitudes file: %v", err)
	}

	streams, err := buildStreams("./testdata", "[EBH][NH]Z")
	if err != nil {
		t.Fatalf("problem loading streams: %v", err)
	}

	b2, err := json.MarshalIndent(streams, "", "   ")
	if err != nil {
		t.Fatalf("problem marshalling streams: %v", err)
	}

	// compare stored with computed
	if strings.TrimSpace(string(b1)) != strings.TrimSpace(string(b2)) {
		t.Errorf("unexpected content -got/+exp\n%s", cmp.Diff(string(b1), string(b2)))
	}
}
