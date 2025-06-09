package delta_test

import (
	"testing"

	"github.com/GeoNet/delta"
	"github.com/GeoNet/delta/meta"
)

var metSensorChecks = map[string]func(*meta.Set) func(t *testing.T){

	"check for metsensors installation equipment overlaps": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {

			installs := make(map[string]meta.InstalledMetSensorList)
			for _, s := range set.InstalledMetSensors() {
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
						if v[i].End.Equal(v[j].Start) {
							continue
						}
						if v[i].Start.Equal(v[j].End) {
							continue
						}

						t.Errorf("metsensors %s at %-5s has mark %s overlap between %s and %s",
							v[i].Model, v[i].Serial, v[i].Mark,
							v[i].Start.Format(meta.DateTimeFormat),
							v[i].End.Format(meta.DateTimeFormat))
					}
				}
			}
		}
	},
	"check for invalid installation dates": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			for _, i := range set.InstalledMetSensors() {
				if i.End.After(i.Start) {
					continue
				}
				t.Errorf("installed metsensor is removed before it has been installed: %s", i.String())
			}
		}
	},
	"check for missing metsensor marks": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			keys := make(map[string]interface{})
			for _, m := range set.Marks() {
				keys[m.Code] = true
			}

			for _, c := range set.InstalledMetSensors() {
				if _, ok := keys[c.Mark]; !ok {
					t.Errorf("unable to find metsensor mark %-5s", c.Mark)
				}
			}
		}
	},
	"check for missing metsensor assets": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {

			for _, r := range set.InstalledMetSensors() {
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
				t.Errorf("unable to find metsensor asset: %s [%s]", r.Model, r.Serial)
			}
		}
	},
}

func TestMetSensors(t *testing.T) {

	set, err := delta.New()
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range metSensorChecks {
		t.Run(k, v(set))
	}
}
