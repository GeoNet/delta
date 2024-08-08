package main

import (
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/GeoNet/delta"
)

func TestSkeleton(t *testing.T) {
	// set recovers the delta tables
	set, err := delta.NewBase("./testdata")
	if err != nil {
		t.Fatal(err)
	}

	tm := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	b, err := skeleton("AVLN", set, tm.UTC().Unix())
	if err != nil {
		t.Fatal(err)
	}

	f, err := os.ReadFile("./testdata/AVLN00NZL.SKL")
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal([]byte(b), f) {
		t.Errorf("generated skeleton file does not match(%d vs %d):\n%s", len(b), len(f), b)
	}
}

// tests the case where we unable to generate a skeleton file by given metadata,
// and the generic skeleton file is used instead.
func TestGenericSkeleton(t *testing.T) {
	set, err := delta.NewBase("./testdata")
	if err != nil {
		t.Fatal(err)
	}

	// set the reference time to before the mark's install time
	tm := time.Date(2005, 1, 1, 0, 0, 0, 0, time.UTC)
	b, err := skeleton("AVLN", set, tm.UTC().Unix())
	if err == nil {
		t.Fatalf("expected to get error but got nil")
	}
	f, err := os.ReadFile("./testdata/AVLN-generic.SKL")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal([]byte(b), f) {
		t.Errorf("generated skeleton file does not match(%d vs %d)", len(b), len(f))
	}
}
