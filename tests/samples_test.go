package delta_test

import (
	"testing"

	"github.com/GeoNet/delta"
	"github.com/GeoNet/delta/meta"
)

var sampleChecks = map[string]func(*meta.Set) func(t *testing.T){
	"check for duplicate sampling site": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			samples := set.Samples()
			for i := 0; i < len(samples); i++ {
				for j := i + 1; j < len(samples); j++ {
					if !samples[i].Overlaps(samples[j]) {
						continue
					}
					t.Errorf("error: sample overlap: %s", samples[i].Id())
				}
			}
		}
	},
	"check for missing sample networks": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			for _, s := range set.Samples() {
				if _, ok := set.Network(s.Network); ok {
					continue
				}
				t.Logf("warning: missing network: %s", s.Id())
			}
		}
	},
}

func TestSamples(t *testing.T) {

	set, err := delta.New()
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range sampleChecks {
		t.Run(k, v(set))
	}
}
