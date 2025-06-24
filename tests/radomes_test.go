package delta_test

import (
	"testing"

	"github.com/GeoNet/delta"
	"github.com/GeoNet/delta/meta"
)

var installedRadomeChecks = map[string]func(*meta.Set) func(t *testing.T){

	"check for radomes installation equipment overlaps": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			installs := make(map[string]meta.InstalledRadomeList)
			for _, s := range set.InstalledRadomes() {
				installs[s.Model] = append(installs[s.Model], s)
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
						t.Errorf("radomes %s at %-5s has mark %s overlap between %s and %s",
							v[i].Model, v[i].Serial, v[i].Mark, v[i].Start.Format(meta.DateTimeFormat), v[i].End.Format(meta.DateTimeFormat))
					}
				}
			}
		}
	},
	"check for overlapping radomes installations": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			installs := make(map[string]meta.InstalledRadomeList)
			for _, s := range set.InstalledRadomes() {
				installs[s.Mark] = append(installs[s.Mark], s)
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

						t.Errorf("mark %-5s has radome %s/%s overlap wth %s/%s between %s and %s",
							v[i].Mark, v[i].Model, v[i].Serial, v[j].Model, v[j].Serial, v[i].Start.Format(meta.DateTimeFormat), v[i].End.Format(meta.DateTimeFormat))
					}
				}
			}
		}
	},

	"check for invalid installation dates": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			for _, i := range set.InstalledRadomes() {
				if i.End.After(i.Start) {
					continue
				}
				t.Errorf("installed radome is removed before it has been installed: %s", i.String())
			}
		}
	},

	"check for missing radome marks": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {

			keys := make(map[string]interface{})
			for _, m := range set.Marks() {
				keys[m.Code] = true
			}

			for _, c := range set.InstalledRadomes() {
				if _, ok := keys[c.Mark]; !ok {
					t.Errorf("unable to find radome mark %-5s", c.Mark)
				}
			}
		}
	},

	"check for missing radome assets": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {

			for _, r := range set.InstalledRadomes() {
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
				t.Errorf("unable to find radome asset: %s [%s]", r.Model, r.Serial)
			}
		}
	},
}

func TestInstalledRadomes(t *testing.T) {

	set, err := delta.New()
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range installedRadomeChecks {
		t.Run(k, v(set))
	}
}
