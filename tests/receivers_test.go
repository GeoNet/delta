package delta_test

import (
	"sort"
	"testing"

	"github.com/GeoNet/delta/meta"
)

func TestReceivers(t *testing.T) {

	var receivers meta.DeployedReceiverList
	t.Log("Load deployed receivers file")
	{
		if err := meta.LoadList("../install/receivers.csv", &receivers); err != nil {
			t.Fatal(err)
		}
	}

	t.Log("Check for particular receiver installation overlaps")
	{
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
			v := installs[k]

			for i, n := 0, len(v); i < n; i++ {
				for j := i + 1; j < n; j++ {
					switch {
					case v[i].Serial != v[j].Serial:
					case v[i].End.Before(v[j].Start):
					case v[i].Start.After(v[j].End):
					case v[i].End.Equal(v[j].Start):
					case v[i].Start.Equal(v[j].End):
					default:
						t.Errorf("receiver %s [%s] at %s has overlap with %s between times %s and %s",
							v[i].Model, v[i].Serial, v[i].Mark, v[j].Mark, v[i].Start.Format(meta.DateTimeFormat), v[i].End.Format(meta.DateTimeFormat))
					}
				}
			}
		}
	}

	t.Log("Check for receiver sites installation equipment overlaps")
	{
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
			v := installs[k]

			for i, n := 0, len(v); i < n; i++ {
				for j := i + 1; j < n; j++ {
					switch {
					case v[i].End.Before(v[j].Start):
					case v[i].Start.After(v[j].End):
					case v[i].End.Equal(v[j].Start):
					case v[i].Start.Equal(v[j].End):
					default:
						t.Errorf("receivers %s [%s] / %s [%s] at %s has overlap between %s and %s",
							v[i].Model, v[i].Serial, v[j].Model, v[j].Serial, v[i].Mark, v[i].Start.Format(meta.DateTimeFormat), v[i].End.Format(meta.DateTimeFormat))
					}
				}
			}
		}
	}

	t.Log("Check for missing receiver marks")
	{
		var marks meta.MarkList

		if err := meta.LoadList("../network/marks.csv", &marks); err != nil {
			t.Fatal(err)
		}

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
	}

	var assets meta.AssetList
	t.Log("Load receiver assets file")
	{
		if err := meta.LoadList("../assets/receivers.csv", &assets); err != nil {
			t.Fatal(err)
		}
	}

	t.Log("Check for receiver assets")
	{
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
	}

}
