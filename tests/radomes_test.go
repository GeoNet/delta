package delta_test

import (
	"sort"
	"testing"

	"github.com/GeoNet/delta/meta"
)

func TestRadomes(t *testing.T) {
	var radomes meta.InstalledRadomeList

	if err := meta.LoadList("../install/radomes.csv", &radomes); err != nil {
		t.Fatal(err)
	}

	var marks meta.MarkList

	if err := meta.LoadList("../network/marks.csv", &marks); err != nil {
		t.Fatal(err)
	}

	var assets meta.AssetList

	if err := meta.LoadList("../assets/radomes.csv", &assets); err != nil {
		t.Fatal(err)
	}

	t.Run("Check for radomes installation equipment overlaps", func(t *testing.T) {
		installs := make(map[string]meta.InstalledRadomeList)
		for _, s := range radomes {
			_, ok := installs[s.Model]
			if ok {
				installs[s.Model] = append(installs[s.Model], s)

			} else {
				installs[s.Model] = meta.InstalledRadomeList{s}
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
					if v[i].End.Equal(v[j].Start) || v[i].Start.Equal(v[j].End) {
						continue
					}
					t.Errorf("radomes %s at %-5s has mark %s overlap between %s and %s",
						v[i].Model, v[i].Serial, v[i].Mark,
						v[i].Start.Format(meta.DateTimeFormat),
						v[i].End.Format(meta.DateTimeFormat))
				}
			}
		}
	})

	t.Run("Check for missing radome marks", func(t *testing.T) {
		keys := make(map[string]interface{})

		for _, m := range marks {
			keys[m.Code] = true
		}

		for _, c := range radomes {
			if _, ok := keys[c.Mark]; ok {
				continue
			}
			t.Errorf("unable to find radome mark %-5s", c.Mark)
		}
	})

	t.Run("Check for radome assets", func(t *testing.T) {
		for _, r := range radomes {
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
				t.Errorf("unable to find radome asset: %s [%s]", r.Model, r.Serial)
			}
		}
	})

}
