package delta_test

import (
	"testing"

	"github.com/GeoNet/delta"
	"github.com/GeoNet/delta/meta"
)

var dartChecks = map[string]func(*meta.Set) func(t *testing.T){

	// check for session overlaps, there can't be two sessions running at the same mark for the same sampling interval.
	"check for duplicated darts": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			darts := set.Darts()
			for i := 0; i < len(darts); i++ {
				for j := i + 1; j < len(darts); j++ {
					if darts[i].Station != darts[j].Station {
						continue
					}
					t.Errorf("dart duplication: %s", darts[i].Station)
				}
			}
		}
	},
}

func TestDarts(t *testing.T) {
	set, err := delta.New()
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range dartChecks {
		t.Run(k, v(set))
	}
}
