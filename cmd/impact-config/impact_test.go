package main

import (
	"encoding/json"
	"os"
	"regexp"
	"testing"

	"github.com/GeoNet/delta"
	"github.com/GeoNet/delta/resp"

	"github.com/google/go-cmp/cmp"
)

func TestBuild_Network(t *testing.T) {

	settings := Settings{
		skip:     *regexp.MustCompile("^SB$"),
		channels: *regexp.MustCompile("^[EBH][NH]Z$"),
	}

	set, err := delta.NewBase("./testdata")
	if err != nil {
		t.Fatal(err)
	}

	// load in the test data and convert to stationxml indented text
	expected, err := os.ReadFile("./testdata/impact.json")
	if err != nil {
		t.Fatalf("error: unable to load test amplitudes file: %v", err)
	}

	streams, err := settings.ImpactStreams(set, resp.NewResp("./testdata/resp"))
	if err != nil {
		t.Fatalf("problem loading streams: %v", err)
	}

	result, err := json.MarshalIndent(streams, "", "   ")
	if err != nil {
		t.Fatalf("problem marshalling streams: %v", err)
	}

	if !cmp.Equal(expected, result) {
		t.Errorf("unexpected content -got/+exp\n%s", cmp.Diff(expected, result))
	}
}
