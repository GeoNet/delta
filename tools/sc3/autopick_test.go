package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

const broadband = `###
### Delivered by puppet
###
# Defines the filter to be used for picking.
detecFilter = "FILTER"

# The time correction applied to a detected pick.
timeCorr = -0.15

# Defines whether or not the streams are picked or not
detecEnable = true
`

func TestAutoPick(t *testing.T) {

	picks := map[string]struct {
		pick    AutoPick
		content string
	}{
		"scautopick/profile_broadband": {
			pick: AutoPick{
				Style:      "broadband",
				Filter:     "FILTER",
				Correction: -0.15,
			},
			content: broadband,
		},
		"scautopick/profile_weak": {
			pick: AutoPick{
				Style:      "weak",
				Filter:     "FILTER",
				Correction: -0.15,
			},
			content: broadband,
		},
	}

	for k, p := range picks {
		t.Run("check "+k, func(t *testing.T) {
			d, err := ioutil.TempDir(os.TempDir(), "test")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(d)

			if err := Store(p.pick, d); err != nil {
				t.Fatalf("unable to store key output %s: %v", k, err)
			}

			key, err := ioutil.ReadFile(filepath.Join(d, k))
			if err != nil {
				t.Fatalf("unable to read temp key file %s: %v", d, err)
			}
			if string(key) != p.content {
				t.Errorf("contents mismatch %s", k)
			}
		})
	}
}
