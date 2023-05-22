package delta_test

import (
	"testing"
	"time"

	"github.com/GeoNet/delta"
	"github.com/GeoNet/delta/meta"
)

var featureChecks = map[string]func(*meta.Set) func(t *testing.T){
	"check for duplicated site features": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			features := set.Features()

			for i := 0; i < len(features); i++ {
				for j := i + 1; j < len(features); j++ {
					if !features[i].Overlaps(features[j]) {
						continue
					}
					t.Error("error: site feature overlaps: " + features[i].Id())
				}
			}
		}
	},
	"check for duplicated features": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			for _, f := range set.Features() {
				if _, ok := set.Station(f.Station); ok {
					continue
				}
				if _, ok := set.Sample(f.Station); ok {
					continue
				}

				t.Error("error: unable to find feature station: " + f.Id())
			}
			for _, f := range set.Features() {
				if s, ok := set.Station(f.Station); ok {
					switch {
					case f.Start.Before(s.Start):
						t.Log("warning: feature start mismatch: " + f.Id() + " before " + s.Start.String())
					case f.End.Before(time.Now()) && f.End.After(s.End):
						t.Log("warning: feature end mismatch: " + f.Id() + " after " + s.End.String())
					}
				}
				if s, ok := set.Sample(f.Station); ok {
					switch {
					case f.Start.Before(s.Start):
						t.Log("warning: feature start mismatch: " + f.Id() + " before " + s.Start.String())
					case f.End.Before(time.Now()) && f.End.After(s.End):
						t.Log("warning: feature end mismatch: " + f.Id() + " after " + s.End.String())
					}
				}
			}
		}
	},
	"check for duplicated feature sites": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			for _, f := range set.Features() {
				if _, ok := set.Site(f.Station, f.Location); !ok {
					t.Error("error: unable to find feature site: " + f.Id())
				}
			}
			for _, f := range set.Features() {
				if s, ok := set.Site(f.Station, f.Location); ok {
					switch {
					case f.Start.Before(s.Start):
						t.Log("warning: feature start mismatch: " + f.Id() + " before " + s.Start.String())
					case s.End.Before(time.Now()) && f.End.After(s.End):
						t.Log("warning: feature end mismatch: " + f.Id() + " after " + s.End.String())
					}
				}
			}
		}
	},
}

func TestFeatures(t *testing.T) {

	set, err := delta.New()
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range featureChecks {
		t.Run(k, v(set))
	}
}
