package delta_test

import (
	"testing"

	"github.com/GeoNet/delta"
	"github.com/GeoNet/delta/meta"
)

var polaritiesChecks = map[string]func(*meta.Set) func(t *testing.T){
	"check for duplicate polarities": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			polarities := set.Polarities()
			for i := 0; i < len(polarities); i++ {
				for j := i + 1; j < len(polarities); j++ {
					if polarities[i].Station != polarities[j].Station {
						continue
					}
					if polarities[i].Location != polarities[j].Location {
						continue
					}
					if polarities[i].Sublocation != polarities[j].Sublocation {
						continue
					}
					if polarities[i].Subsource != polarities[j].Subsource {
						continue
					}
					if polarities[i].Primary && !polarities[j].Primary {
						continue
					}
					if !polarities[i].Primary && polarities[j].Primary {
						continue
					}
					if !polarities[i].Span.Overlaps(polarities[j].Span) {
						continue
					}

					t.Errorf("polarity overlap, %s %s %s", polarities[i].Station, polarities[i].Location, polarities[i].Span)
				}
			}
		}
	},
}

func TestPolarities(t *testing.T) {

	set, err := delta.New()
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range polaritiesChecks {
		t.Run(k, v(set))
	}
}
