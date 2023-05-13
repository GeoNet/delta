package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

const stationStrong = `###
### Warning, this file is automaticaly generated and may be overwritten.
###
global:strong_lowrate_20
scautopick:strong
`

const stationWeak = `###
### Warning, this file is automaticaly generated and may be overwritten.
###
global:weak_10
scautopick:weak
`

func TestStation(t *testing.T) {

	stations := map[string]struct {
		station Station
		content string
	}{
		"station_NZ_MAGS": {
			station: Station{
				Global: Global{
					Location: "20",
					Stream:   "BN",
				},
				Code:    "MAGS",
				Network: "NZ",
				AutoPick: AutoPick{
					Style:      "strong",
					Filter:     DefaultFilter,
					Correction: DefaultCorrection,
				},
			},
			content: stationStrong,
		},
		"station_NZ_CAW": {
			station: Station{
				Global: Global{
					Location: "10",
					Stream:   "EH",
				},
				Code:    "CAW",
				Network: "NZ",
				AutoPick: AutoPick{
					Style:      "weak",
					Filter:     DefaultFilter,
					Correction: DefaultCorrection,
				},
			},
			content: stationWeak,
		},
	}

	for k, s := range stations {
		t.Run("check "+k, func(t *testing.T) {
			d, err := os.MkdirTemp(os.TempDir(), "test")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(d)

			if err := Store(s.station, d); err != nil {
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
