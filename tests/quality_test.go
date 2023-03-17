package delta_test

import (
	"testing"

	"github.com/GeoNet/delta"
	"github.com/GeoNet/delta/meta"
)

var qualityChecks = map[string]func(*meta.Set) func(t *testing.T){

	"check for overlapping qualities": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			qualities := set.Qualities()
			for i := 0; i < len(qualities); i++ {
				for j := i + 1; j < len(qualities); j++ {
					if qualities[i].Station != qualities[j].Station {
						continue
					}
					if qualities[i].Location != qualities[j].Location {
						continue
					}
					if !qualities[i].Span.Overlaps(qualities[j].Span) {
						continue
					}

					t.Errorf("error: quality overlap for %q", qualities[i].String())
				}
			}
		}
	},
}

func TestQualities(t *testing.T) {

	set, err := delta.New()
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range qualityChecks {
		t.Run(k, v(set))
	}
}
