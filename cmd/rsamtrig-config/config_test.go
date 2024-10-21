package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestConfig(t *testing.T) {

	var good = map[string]Config{
		"MAVZ/11:330": {
			Station:  "MAVZ",
			Location: "11",
			Level:    330.0,
		},
		"TRVZ:130": {
			Station: "TRVZ",
			Level:   130.0,
		},
		"WHVZ/11": {
			Station:  "WHVZ",
			Location: "11",
		},
		"WHVZ": {
			Station: "WHVZ",
		},
	}

	for k, v := range good {
		t.Run(k, func(t *testing.T) {

			config, err := NewConfig(k)
			if err != nil {
				t.Fatal(err)
			}

			if !cmp.Equal(v, config) {
				t.Errorf("unexpected %q content -got/+exp\n%s", k, cmp.Diff(v, config))
			}
		})
	}

	var bad = []string{
		" ",
		"",
		"/",
		":",
		"?",
		" A ",
		" A:100/20 ",
	}

	for _, k := range bad {
		t.Run(k, func(t *testing.T) {

			if _, err := NewConfig(k); err == nil {
				t.Fatalf("unexpected %q content - should have failed to compile", k)
			}
		})
	}
}
