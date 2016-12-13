package main

import (
	"encoding/xml"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"testing"

	"github.com/ozym/fdsn/stationxml"
)

func TestBuild(t *testing.T) {

	// load meta information
	md, err := NewMeta("./testdata/network", "./testdata/install")
	if err != nil {
		t.Fatalf("unable to load meta data: %v", err)
	}

	// select all test networks/stations/channels
	re, err := regexp.Compile("[A-Z0-9]+")
	if err != nil {
		t.Fatalf("unable to compile regexp: %v", err)
	}

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

	// build networks and construct stationxml
	n, err := buildNetworks(md, re, re, re)
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
		f1.Write(b1)

		f2, err := ioutil.TempFile(os.TempDir(), "tmp")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(f2.Name())
		f2.Write(b2)

		cmd := exec.Command("diff", "-c", f1.Name(), f2.Name())
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			t.Fatal(err)
		}
		err = cmd.Start()
		if err != nil {
			t.Fatal(err)
		}
		defer cmd.Wait()
		diff, err := ioutil.ReadAll(stdout)
		if err != nil {
			t.Fatal(err)
		}
		t.Error(string(diff))
	}
}
