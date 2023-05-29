package delta_test

import (
	"testing"

	"github.com/GeoNet/delta"
	"github.com/GeoNet/delta/meta"
)

var deployedReceiverChecks = map[string]func(*meta.Set) func(t *testing.T){
	"check for receiver installation overlaps": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {

			installs := make(map[string]meta.DeployedReceiverList)
			for _, s := range set.DeployedReceivers() {
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

						t.Errorf("receiver %s [%s] at %s has overlap with %s between times %s and %s",
							v[i].Model, v[i].Serial, v[i].Mark, v[j].Mark, v[i].Start.Format(meta.DateTimeFormat), v[i].End.Format(meta.DateTimeFormat))
					}
				}
			}
		}
	},

	"check for receiver installation equipment overlaps": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {

			installs := make(map[string]meta.DeployedReceiverList)
			for _, s := range set.DeployedReceivers() {
				installs[s.Mark] = append(installs[s.Mark], s)
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

						t.Errorf("receivers %s [%s] / %s [%s] at %s has overlap between %s and %s",
							v[i].Model, v[i].Serial, v[j].Model, v[j].Serial, v[i].Mark, v[i].Start.Format(meta.DateTimeFormat), v[i].End.Format(meta.DateTimeFormat))
					}
				}
			}
		}
	},

	"check for missing receiver marks": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {

			keys := make(map[string]interface{})
			for _, m := range set.Marks() {
				keys[m.Code] = true
			}

			for _, r := range set.DeployedReceivers() {
				if _, ok := keys[r.Mark]; ok {
					continue
				}
				t.Errorf("unable to find receiver mark %-5s", r.Mark)
			}
		}
	},

	"check for missing receiver assets": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			for _, r := range set.DeployedReceivers() {
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
				t.Errorf("unable to find receiver asset: %s [%s]", r.Model, r.Serial)
			}
		}
	},
}

func TestDeployedReceivers(t *testing.T) {

	set, err := delta.New()
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range deployedReceiverChecks {
		t.Run(k, v(set))
	}
}
