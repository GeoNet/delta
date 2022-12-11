package main

import (
	"io/ioutil"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBuild_Tsunami(t *testing.T) {

	// load in the test data and convert to stationxml indented text
	b1, err := ioutil.ReadFile("./testdata/tsunami-gauge.xml")
	if err != nil {
		t.Fatalf("error: unable to load test amplitudes file: %v", err)
	}

	plots, err := ReadPlots("./testdata/chart-tsunami.yaml")
	if err != nil {
		t.Fatalf("error: unable to load test tsunami config file: %v", err)
	}
	var pages []Page
	for plot, cfg := range plots.Configs {
		for _, p := range cfg.Pages {
			res, err := p.Tsunami("./testdata")
			if err != nil {
				t.Fatalf("problem build tsunami pages %s: %v", plot, err)
			}
			pages = append(pages, res...)
		}
	}

	c := Chart{
		Pages: pages,
	}

	b2, err := c.Marshal()
	if err != nil {
		t.Fatalf("error: unable to encode test tsunami file: %v", err)
	}

	// compare stored with computed
	if string(b1) != string(b2) {
		t.Errorf("unexpected content -got/+exp\n%s", cmp.Diff(string(b1), string(b2)))
	}
}
