package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/GeoNet/delta"
	"github.com/GeoNet/delta/internal/ntrip"
)

func TestFiles_Config(t *testing.T) {

	// what do we expect to see
	raw, err := os.ReadFile(filepath.Join("testdata", "config.yaml"))
	if err != nil {
		t.Fatal(err)
	}

	// set recovers the delta tables
	set, err := delta.NewBase("./testdata")
	if err != nil {
		t.Fatal(err)
	}

	// recover ntrip input files
	caster, err := ntrip.NewCaster("./testdata", "./testdata")
	if err != nil {
		t.Fatal(err)
	}

	// build config from test files
	config, err := NewConfig(set, caster)
	if err != nil {
		t.Fatal(err)
	}

	// yaml encoded output
	ans, err := config.Marshal()
	if err != nil {
		t.Fatal(err)
	}

	// they should be the same
	if !bytes.Equal(raw, ans) {
		t.Errorf("unexpected content -got/+exp\n%s", cmp.Diff(string(raw), string(ans)))
	}
}
