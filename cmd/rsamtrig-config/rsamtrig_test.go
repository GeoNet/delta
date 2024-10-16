package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

const rsam_330 = `###
### Delivered by puppet
###
# Defines a list of filters to apply.
rsamFilters = f1-4

# Waveform filter string to apply before RSAM calculation.
rsamFilter.f1-4.filter = BW(4,1,4)

# Minimum amplitude level. Used as starting point and lower barrier.
rsamFilter.f1-4.baseLevel = 330
`

func TestRsamTrig(t *testing.T) {

	trigs := map[string]struct {
		trig    RsamTrig
		content string
	}{
		"rsamtrig/profile_rsam_mavz": {
			trig: RsamTrig{
				Station: "MAVZ",
				Name:    DefaultName,
				Filter:  DefaultFilter,
				Base:    330,
			},
			content: rsam_330,
		},
	}

	for k, v := range trigs {
		t.Run("check "+k, func(t *testing.T) {
			d, err := os.MkdirTemp(os.TempDir(), "test")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(d)

			if err := Store(v.trig, d); err != nil {
				t.Fatalf("unable to store key output %s: %v", k, err)
			}

			key, err := os.ReadFile(filepath.Join(d, k))
			if err != nil {
				t.Fatalf("unable to read temp key file %s: %v", d, err)
			}
			if s := string(key); s != v.content {
				t.Errorf("unexpected content -got/+exp\n%s", cmp.Diff(s, v.content))
			}
		})
	}
}
