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
