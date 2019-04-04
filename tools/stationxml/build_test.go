package main

import (
	"encoding/xml"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"

	"github.com/ozym/fdsn/stationxml"
)

func TestBuilder(t *testing.T) {

	// load in the test data and convert to stationxml indented text
	raw, err := ioutil.ReadFile("./testdata/test.xml")
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
		t.Error("**** stationxml mismatch ****")

		f1, err := ioutil.TempFile(os.TempDir(), "tmp")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(f1.Name())
		if _, err := f1.Write(b1); err != nil {
			t.Fatal(err)
		}

		f2, err := ioutil.TempFile(os.TempDir(), "tmp")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(f2.Name())
		if _, err := f2.Write(b2); err != nil {
			t.Fatal(err)
		}

		cmd := exec.Command("diff", "-c", f1.Name(), f2.Name())
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			t.Fatal(err)
		}
		err = cmd.Start()
		if err != nil {
			t.Fatal(err)
		}
		diff, err := ioutil.ReadAll(stdout)
		if err != nil {
			t.Fatal(err)
		}
		t.Error(string(diff))

		if err := cmd.Wait(); err != nil {
			t.Fatal(err)
		}
	}
}
