package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

const stationContent = `###
### Delivered by puppet
###
capslink:sl4caps
`

func TestStation(t *testing.T) {

	stations := map[string]struct {
		station Station
		content string
	}{
		"station_NZ_MAGS": {
			station: Station{
				Network: "NZ",
				Code:    "MAGS",
			},
			content: stationContent,
		},
		"station_NZ_CAW": {
			station: Station{
				Network: "NZ",
				Code:    "CAW",
			},
			content: stationContent,
		},
	}

	for k, s := range stations {
		t.Run("check "+k, func(t *testing.T) {
			d, err := os.MkdirTemp(os.TempDir(), "test")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(d)

			if err := s.station.Output(d); err != nil {
				t.Fatalf("unable to store key output %s: %v", k, err)
			}

			key, err := os.ReadFile(filepath.Join(d, k))
			if err != nil {
				t.Fatalf("unable to read temp key file %s: %v", d, err)
			}
			if v := string(key); v != s.content {
				t.Errorf("unexpected content -got/+exp\n%s", cmp.Diff(v, s.content))
			}
		})
	}
}
