package delta_test

import (
	"sort"
	"testing"

	"github.com/GeoNet/delta/meta"
)

func TestReceivers(t *testing.T) {
	var receivers meta.DeployedReceiverList

	if err := meta.LoadList("../install/receivers.csv", &receivers); err != nil {
		t.Fatal(err)
	}

	var marks meta.MarkList

	if err := meta.LoadList("../network/marks.csv", &marks); err != nil {
		t.Fatal(err)
	}

	var assets meta.AssetList

	if err := meta.LoadList("../assets/receivers.csv", &assets); err != nil {
		t.Fatal(err)
	}

	t.Run("Check for particular receiver installation overlaps", func(t *testing.T) {
		installs := make(map[string]meta.DeployedReceiverList)
		for _, s := range receivers {
			if _, ok := installs[s.Model]; !ok {
				installs[s.Model] = meta.DeployedReceiverList{}
			}
			installs[s.Model] = append(installs[s.Model], s)
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
					t.Errorf("receiver %s [%s] at %s has overlap with %s between times %s and %s",
						v[i].Model, v[i].Serial, v[i].Mark, v[j].Mark,
						v[i].Start.Format(meta.DateTimeFormat), v[i].End.Format(meta.DateTimeFormat))
				}
			}
		}
	})

	t.Run("Check for receiver sites installation equipment overlaps", func(t *testing.T) {
		installs := make(map[string]meta.DeployedReceiverList)
		for _, s := range receivers {
			if _, ok := installs[s.Mark]; !ok {
				installs[s.Mark] = meta.DeployedReceiverList{}
			}
			installs[s.Mark] = append(installs[s.Model], s)
		}

		var keys []string
		for k, _ := range installs {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			for i, v, n := 0, installs[k], len(installs[k]); i < n; i++ {
				for j := i + 1; j < n; j++ {
					if v[i].End.Before(v[j].Start) || v[i].Start.After(v[j].End) {
						continue
					}
					if v[i].End.Equal(v[j].Start) || v[i].Start.Equal(v[j].End) {
						continue
					}
					t.Errorf("receivers %s [%s] / %s [%s] at %s has overlap between %s and %s",
						v[i].Model, v[i].Serial, v[j].Model, v[j].Serial, v[i].Mark,
						v[i].Start.Format(meta.DateTimeFormat), v[i].End.Format(meta.DateTimeFormat))
				}
			}
		}
	})

	t.Run("Check for missing receiver marks", func(t *testing.T) {
		keys := make(map[string]interface{})

		for _, m := range marks {
			keys[m.Code] = true
		}

		for _, r := range receivers {
			if _, ok := keys[r.Mark]; ok {
				continue
			}
			t.Errorf("unable to find receiver mark %-5s", r.Mark)
		}
	})

	t.Run("Check for receiver assets", func(t *testing.T) {
		for _, r := range receivers {
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
				t.Errorf("unable to find receiver asset: %s [%s]", r.Model, r.Serial)
			}
		}
	})

}
