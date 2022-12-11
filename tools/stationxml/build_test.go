package main

import (
	"encoding/xml"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/ozym/fdsn/stationxml"
)

func TestBuilder(t *testing.T) {

	// load in the test data and convert to stationxml indented text
	raw, err := os.ReadFile("./testdata/test.xml")
	if err != nil {
		t.Fatalf("error: unable to load test stationxml file: %v", err)
	}

	var x stationxml.FDSNStationXML
	if err := xml.Unmarshal(raw, &x); err != nil {
		t.Fatalf("error: unable to unmarshal stationxml file: %v", err)
	}

	b1, err := xml.MarshalIndent(&x, "", "  ")
	if err != nil {
		t.Fatalf("error: unable to unmarshal stationxml file: %v", err)
	}

	var builder Builder

	// build networks and construct stationxml
	n, err := builder.Construct("./testdata")
	if err != nil {
		t.Fatalf("error: unable to build networks list: %v", err)
	}

	c, err := stationxml.ParseDateTime("2016-12-13T11:40:46")
	if err != nil {
		t.Fatalf("error: unable to parse creation date: %v", err)
	}

	y := stationxml.FDSNStationXML{
		NameSpace:     stationxml.FDSNNameSpace,
		SchemaVersion: stationxml.FDSNSchemaVersion,
		Source:        "GeoNet",
		Sender:        "WEL(GNS_Test)",
		Module:        "Delta",
		Networks:      n,
		Created:       c,
	}

	b2, err := xml.MarshalIndent(&y, "", "  ")
	if err != nil {
		t.Fatalf("error: unable to marshal stationxml data: %v", err)
	}

	// compare stored with computed
	if string(b1) != string(b2) {
		t.Errorf("unexpected content -got/+exp\n%s", cmp.Diff(string(b1), string(b2)))
	}
}
