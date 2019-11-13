package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

const stationStrong = `###
### Delivered by puppet
###
global:strong_lowrate_20
scautopick:strong
`

const stationWeak = `###
### Delivered by puppet
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
			d, err := ioutil.TempDir(os.TempDir(), "test")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(d)

			if err := Store(s.station, d); err != nil {
				t.Fatalf("unable to store key output %s: %v", k, err)
			}

			key, err := ioutil.ReadFile(filepath.Join(d, k))
			if err != nil {
				t.Fatalf("unable to read temp key file %s: %v", d, err)
			}
			if string(key) != s.content {
				t.Errorf("contents mismatch %s", k)
			}
		})
	}
}
