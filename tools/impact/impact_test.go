package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"
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
		t.Error("**** impact json mismatch ****")

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
