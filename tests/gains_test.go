package delta_test

import (
	"testing"

	"github.com/GeoNet/delta"
	"github.com/GeoNet/delta/meta"
)

var gainChecks = map[string]func(*meta.Set) func(t *testing.T){
	"check for gain installation overlaps": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			installs := make(map[string]meta.GainList)
			for _, s := range set.Gains() {
				for _, c := range s.Gains() {
					installs[c.Id()] = append(installs[c.Id()], c)
				}
			}
			for _, v := range installs {
				for i := 0; i < len(v); i++ {
					for j := i + 1; j < len(v); j++ {
						if v[i].End.Before(v[j].Start) {
							continue
						}
						if v[i].Start.After(v[j].End) {
							continue
						}
						if v[i].End.Equal(v[j].Start) {
							continue
						}
						if v[i].Start.Equal(v[j].End) {
							continue
						}

						t.Errorf("gain %s/%s has component %s overlap between %s and %s",
							v[i].Station, v[i].Location, v[i].Subsource,
							v[i].Start.Format(meta.DateTimeFormat),
							v[i].End.Format(meta.DateTimeFormat))
					}
				}
			}
		}
	},
	"check for missing sites": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			for _, s := range set.Gains() {
				var found bool
				for _, a := range set.Sites() {
					if a.Station != s.Station {
						continue
					}
					if a.Location != s.Location {
						continue
					}
					found = true
				}
				if found {
					continue
				}
				t.Errorf("unable to find gain site: %s [%s]", s.Station, s.Location)
			}
		}
	},
}

func TestGains(t *testing.T) {
	set, err := delta.New()
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range gainChecks {
		t.Run(k, v(set))
	}
}
