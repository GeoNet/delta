package main

import (
	"testing"
)

func TestMatcher(t *testing.T) {

	var checks = map[string]map[string]bool{
		`^(NZ)$`: {
			"":   false,
			"NZ": true,
			"SB": false,
		},
		`^()$`: {
			"":   true,
			"NZ": false,
			"SB": false,
		},
		`^[A-Z0-9]+$`: {
			"":   false,
			"NZ": true,
			"SB": true,
			"A0": true,
			"A+": false,
		},
		`^(SB|.X)$`: {
			"SB": true,
			"SX": true,
			"NZ": false,
		},
	}

	// check positive matches
	for k, v := range checks {
		m, err := NewMatcher(k)
		if err != nil {
			t.Fatal(err)
		}
		for s, b := range v {
			if ok := m.MatchString(s); ok != b {
				t.Errorf("mismatch for %q with %q, expected %v but got %v", k, s, b, ok)
			}
		}
	}

	// check negative matches
	for k, v := range checks {
		m, err := NewMatcher("!" + k)
		if err != nil {
			t.Fatal(err)
		}
		for s, b := range v {
			if ok := m.MatchString(s); ok == b {
				t.Errorf("mismatch for !%q with %q, expected %v but got %v", k, s, !b, ok)
				t.Error(k, m)
			}
		}
	}
}
