package main

import (
	"testing"
)

func TestSettings(t *testing.T) {

	stations := []string{"NZ_MAGS", "WL_CAW", "NZ_TEST"}

	checks := map[string]map[string]interface{}{
		"*_*":     {"NZ_MAGS": true, "WL_CAW": true, "NZ_TEST": true},
		"NZ_*":    {"NZ_MAGS": true, "NZ_TEST": true},
		"NZ_MAGS": {"NZ_MAGS": true},
		"*_MAGS":  {"NZ_MAGS": true},
		"WL_*":    {"WL_CAW": true},
		"WL_CAW":  {"WL_CAW": true},
		"*_CAW":   {"WL_CAW": true},
	}

	for k, v := range checks {
		t.Run("check station includes "+k, func(t *testing.T) {

			settings := Settings{
				includeStations: []Station{NewStation(k)},
			}

			for _, s := range stations {
				_, ok := v[s]

				if m := settings.ExcludeStation(NewStation(s)); m == ok {
					t.Errorf("unexpected content for %q, got %v but expected %v", s, m, ok)
				}
			}
		})

		t.Run("check station excludes "+k, func(t *testing.T) {

			settings := Settings{
				excludeStations: []Station{NewStation(k)},
			}

			for _, s := range stations {
				_, ok := v[s]

				if m := settings.ExcludeStation(NewStation(s)); m != ok {
					t.Errorf("unexpected content for %q, got %v but expected %v", s, m, ok)
				}
			}
		})
	}
}
