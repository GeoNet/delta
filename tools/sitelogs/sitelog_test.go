package main

import (
	//	"bytes"
	"encoding/xml"
	"io/ioutil"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSitelog_Read(t *testing.T) {

	// read in a raw file
	raw, err := ioutil.ReadFile("testdata/yald_20200929.xml")
	if err != nil {
		t.Fatal(err)
	}

	var input SiteLogInput
	if err := xml.Unmarshal(raw, &input); err != nil {
		t.Fatal(err)
	}

	// translate to standard sitelog (same details)
	sitelog := input.SiteLog()

	// write as an expected sitelog
	out, err := sitelog.MarshalLegacy()
	if err != nil {
		t.Fatal(err)
	}

	if string(raw) != string(out) {
		t.Errorf("unexpected content -got/+exp\n%s", cmp.Diff(string(raw), string(out)))
	}
}
