package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

const stationBroadband = `###
### Delivered by puppet
###
global:broadband_10
rsamtrig:rsam_mavz
`

const stationWeak = `###
### Delivered by puppet
###
global:weak_10
rsamtrig:rsam_caw
`

func TestStation(t *testing.T) {

	stations := map[string]struct {
		station Station
		content string
	}{
		"station_NZ_CAW": {
			station: Station{
				Global: Global{
					Location: "10",
					Stream:   "EHZ",
				},
				Code:    "CAW",
				Network: "NZ",
			},
			content: stationWeak,
		},
		"station_NZ_MAVZ": {
			station: Station{
				Global: Global{
					Location: "10",
					Stream:   "HHZ",
				},
				Code:    "MAVZ",
				Network: "NZ",
			},
			content: stationBroadband,
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
