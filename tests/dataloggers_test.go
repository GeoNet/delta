package delta_test

import (
	"sort"
	"testing"

	"github.com/GeoNet/delta/meta"
)

func TestDataloggers(t *testing.T) {
	var dataloggers meta.DeployedDataloggerList

	if err := meta.LoadList("../install/dataloggers.csv", &dataloggers); err != nil {
		t.Fatal(err)
	}

	var assets meta.AssetList

	if err := meta.LoadList("../assets/dataloggers.csv", &assets); err != nil {
		t.Fatal(err)
	}

	t.Run("Check for datalogger installation place overlaps", func(t *testing.T) {
		installs := make(map[string]meta.DeployedDataloggerList)
		for _, d := range dataloggers {
			_, ok := installs[d.Place]
			if ok {
				installs[d.Place] = append(installs[d.Place], d)

			} else {
				installs[d.Place] = meta.DeployedDataloggerList{d}
			}
		}

		var keys []string
		for k, _ := range installs {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			v := installs[k]

			for i, n := 0, len(v); i < n; i++ {
				for j := i + 1; j < n; j++ {
					if v[i].Place != v[j].Place {
						continue
					}
					if v[i].Role != v[j].Role {
						continue
					}
					if v[i].End.Before(v[j].Start) || v[i].Start.After(v[j].End) {
						continue
					}
					t.Errorf("datalogger %s:[%s] at %-32s has place overlap between %s and %s",
						v[i].Model, v[i].Serial, v[i].Place,
						v[i].Start.Format(meta.DateTimeFormat),
						v[i].End.Format(meta.DateTimeFormat))
				}
			}
		}
	})

	t.Run("Check for datalogger installation equipment overlaps", func(t *testing.T) {
		installs := make(map[string]meta.DeployedDataloggerList)
		for _, s := range dataloggers {
			_, ok := installs[s.Model]
			if ok {
				installs[s.Model] = append(installs[s.Model], s)

			} else {
				installs[s.Model] = meta.DeployedDataloggerList{s}
			}
		}

		var keys []string
		for k, _ := range installs {
			keys = append(keys, k)
		}

		sort.Strings(keys)

		for _, k := range keys {
			for i, v, n := 0, installs[k], len(installs[k]); i < n; i++ {
				for j := i + 1; j < n; j++ {
					if v[i].Serial != v[j].Serial {
						continue
					}
					if v[i].End.Before(v[j].Start) || v[i].Start.After(v[j].End) {
						continue
					}
					t.Errorf("datalogger %s:[%s] at %-32s has installation overlap between %s and %s",
						v[i].Model, v[i].Serial, v[i].Place,
						v[i].Start.Format(meta.DateTimeFormat),
						v[i].End.Format(meta.DateTimeFormat))
				}
			}
		}
	})

	t.Run("Check for datalogger assets", func(t *testing.T) {
		for _, r := range dataloggers {
			var found bool
			for _, a := range assets {
				if a.Model != r.Model {
					continue
				}
				if a.Serial != r.Serial {
					continue
				}
				found = true
			}
			if !found {
				t.Errorf("unable to find datalogger asset: %s [%s]", r.Model, r.Serial)
			}
		}
	})
}
