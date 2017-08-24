package delta_test

import (
	"sort"
	"testing"

	"github.com/GeoNet/delta/meta"
)

func TestRadomes(t *testing.T) {

	var radomes meta.InstalledRadomeList
	t.Log("Load deployed radomes file")
	{
		if err := meta.LoadList("../install/radomes.csv", &radomes); err != nil {
			t.Fatal(err)
		}
	}

	t.Log("Check for radomes installation equipment overlaps")
	{
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
			v := installs[k]

			for i, n := 0, len(v); i < n; i++ {
				for j := i + 1; j < n; j++ {
					switch {
					case v[i].Serial != v[j].Serial:
					case v[i].End.Before(v[j].Start):
					case v[i].Start.After(v[j].End):
					default:
						t.Errorf("radomes %s at %-5s has mark %s overlap between %s and %s",
							v[i].Model, v[i].Serial, v[i].Mark, v[i].Start.Format(meta.DateTimeFormat), v[i].End.Format(meta.DateTimeFormat))
					}
				}
			}
		}
	}

	t.Log("Check for overlapping radomes installations")
	{
		installs := make(map[string]meta.InstalledRadomeList)
		for _, s := range radomes {
			_, ok := installs[s.Mark]
			if ok {
				installs[s.Mark] = append(installs[s.Mark], s)

			} else {
				installs[s.Mark] = meta.InstalledRadomeList{s}
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
					switch {
					case v[i].Serial != v[j].Serial:
					case v[i].End.Before(v[j].Start):
					case v[i].Start.After(v[j].End):
					default:
						t.Errorf("mark %-5s has radome %s/%s overlap wth %s/%s between %s and %s",
							v[i].Mark, v[i].Model, v[i].Serial, v[j].Model, v[j].Serial, v[i].Start.Format(meta.DateTimeFormat), v[i].End.Format(meta.DateTimeFormat))
					}
				}
			}
		}
	}

	t.Log("Check for missing radome marks")
	{
		var marks meta.MarkList

		if err := meta.LoadList("../network/marks.csv", &marks); err != nil {
			t.Fatal(err)
		}

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
	}

	var assets meta.AssetList
	t.Log("Load radome assets file")
	{
		if err := meta.LoadList("../assets/radomes.csv", &assets); err != nil {
			t.Fatal(err)
		}
	}

	t.Log("Check for radome assets")
	{
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
	}

}
