package main

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/GeoNet/delta"
)

func TestTide(t *testing.T) {

	var checks = map[string]string{
		"./testdata/depth.tmpl":  "./testdata/depth.out",
		"./testdata/detide.tmpl": "./testdata/detide.out",
	}

	set, err := delta.NewBase("./testdata")
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range checks {
		t.Run("check "+k, func(t *testing.T) {

			// read the template file
			tmpl, err := os.ReadFile(k)
			if err != nil {
				t.Fatal(err)
			}

			// read the expected results
			expected, err := os.ReadFile(v)
			if err != nil {
				t.Fatal(err)
			}

			// build the output results
			output, err := Parse(set, tmpl)
			if err != nil {
				t.Fatal(err)
			}

			if !cmp.Equal(expected, output) {
				t.Errorf("unexpected %q content -got/+exp\n%s", k, cmp.Diff(expected, output))
			}
		})
	}
}
