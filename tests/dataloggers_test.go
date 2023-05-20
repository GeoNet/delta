package delta_test

import (
	"testing"

	"github.com/GeoNet/delta"
	"github.com/GeoNet/delta/meta"
)

var dataloggerChecks = map[string]func(*meta.Set) func(t *testing.T){
	"check for datalogger installation place overlaps": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			installs := make(map[string]meta.DeployedDataloggerList)
			for _, d := range set.DeployedDataloggers() {
				installs[d.Place] = append(installs[d.Place], d)
			}

			for _, v := range installs {
				for i := 0; i < len(v); i++ {
					for j := i + 1; j < len(v); j++ {
						if v[i].Place != v[j].Place {
							continue
						}
						if v[i].Role != v[j].Role {
							continue
						}
						if v[i].End.Before(v[j].Start) {
							continue
						}
						if v[i].Start.After(v[j].End) {
							continue
						}

						t.Errorf("datalogger %s:[%s] at %-32s has place overlap between %s and %s",
							v[i].Model, v[i].Serial, v[i].Place,
							v[i].Start.Format(meta.DateTimeFormat),
							v[i].End.Format(meta.DateTimeFormat))
					}
				}
			}
		}
	},
	"check for datalogger installation equipment overlaps": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			installs := make(map[string]meta.DeployedDataloggerList)
			for _, s := range set.DeployedDataloggers() {
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

						t.Errorf("datalogger %s:[%s] at %-32s has installation overlap between %s and %s",
							v[i].Model, v[i].Serial, v[i].Place,
							v[i].Start.Format(meta.DateTimeFormat),
							v[i].End.Format(meta.DateTimeFormat))
					}
				}
			}
		}
	},

	"check for missing datalogger assets": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			for _, r := range set.DeployedDataloggers() {
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
				t.Errorf("unable to find datalogger asset: %s [%s]", r.Model, r.Serial)
			}
		}
	},
}

func TestDataloggers(t *testing.T) {

	set, err := delta.New()
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range dataloggerChecks {
		t.Run(k, v(set))
	}
}
