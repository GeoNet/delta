package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
)

func TestBuild(t *testing.T) {

	// load in the test data and convert to stationxml indented text
	b1, err := ioutil.ReadFile("./testdata/spectra.xml")
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
		t.Error("**** spectras xml mismatch ****")

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
