package delta_test

import (
	"testing"

	"github.com/GeoNet/delta"
	"github.com/GeoNet/delta/meta"
)

var preampsChecks = map[string]func(*meta.Set) func(t *testing.T){
	"check for duplicate preamps": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			preamps := set.Preamps()
			for i := 0; i < len(preamps); i++ {
				for j := i + 1; j < len(preamps); j++ {
					if preamps[i].Station != preamps[j].Station {
						continue
					}
					if preamps[i].Location != preamps[j].Location {
						continue
					}
					if preamps[i].Subsource != preamps[j].Subsource {
						continue
					}
					if !preamps[i].Span.Overlaps(preamps[j].Span) {
						continue
					}

					t.Errorf("polarity overlap, %s %s %s", preamps[i].Station, preamps[i].Location, preamps[i].Span)
				}
			}
		}
	},
}

func TestPreamps(t *testing.T) {

	set, err := delta.New()
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range preampsChecks {
		t.Run(k, v(set))
	}
}
