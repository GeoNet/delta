package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
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
		t.Error("**** amplitudes xml mismatch ****")

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
