package delta_test

import (
	"testing"

	"github.com/GeoNet/delta"
	"github.com/GeoNet/delta/meta"
)

var doasChecks = map[string]func(*meta.Set) func(t *testing.T){

	"check for doases installation equipment overlaps": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {

			installs := make(map[string]meta.InstalledDoasList)
			for _, c := range set.Doases() {
				installs[c.Model] = append(installs[c.Model], c)
			}

			for _, v := range installs {
				for i := 0; i < len(v); i++ {
					for j := i + 1; j < len(v); j++ {
						if v[i].Serial != v[j].Serial {
							continue
						}
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

						t.Errorf("doases %s at %-5s has mount %s overlap between %s and %s",
							v[i].Model, v[i].Serial, v[i].Mount,
							v[i].Start.Format(meta.DateTimeFormat),
							v[i].End.Format(meta.DateTimeFormat))
					}
				}
			}
		}
	},

	"check for doases installation mount overlaps": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			keys := make(map[string]interface{})
			for _, m := range set.Mounts() {
				keys[m.Code] = true
			}

			for _, c := range set.Doases() {
				if _, ok := keys[c.Mount]; !ok {
					t.Errorf("unable to find doas mount %-5s", c.Mount)
				}
			}
		}
	},

	"check for doases installation views": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			type view struct{ m, c string }
			keys := make(map[view]interface{})
			for _, m := range set.Views() {
				keys[view{m.Mount, m.Code}] = true
			}

			for _, c := range set.Doases() {
				if _, ok := keys[view{c.Mount, c.View}]; !ok {
					t.Errorf("unable to find doas mount %-5s (%-2s)", c.Mount, c.View)
				}
			}
		}
	},

	"check doases assets": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {

			for _, r := range set.Doases() {
				var found bool
				for _, a := range set.Assets() {
					if a.Model != r.Model {
						continue
					}
					if a.Serial != r.Serial {
						continue
					}
					found = true
				}
				if found {
					continue
				}
				t.Errorf("unable to find doas asset: %s [%s]", r.Model, r.Serial)
			}
		}
	},
}

func TestDoases(t *testing.T) {

	set, err := delta.New()
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range doasChecks {
		t.Run(k, v(set))
	}
}
