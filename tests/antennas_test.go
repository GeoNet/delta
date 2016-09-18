package delta_test

import (
	"sort"
	"testing"

	"github.com/GeoNet/delta/meta"
)

func TestAntennas(t *testing.T) {

	var antennas meta.InstalledAntennaList
	t.Log("Load deployed antennas file")
	{
		if err := meta.LoadList("../install/antennas.csv", &antennas); err != nil {
			t.Fatal(err)
		}
	}

	t.Log("Check for antenna installation equipment overlaps")
	{
		installs := make(map[string]meta.InstalledAntennaList)
		for _, s := range antennas {
			if _, ok := installs[s.Model]; !ok {
				installs[s.Model] = meta.InstalledAntennaList{}
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
						t.Errorf("antennas %s [%s] at %s has overlap at %s between %s and %s",
							v[i].Model, v[i].Serial, v[i].Mark, v[j].Mark, v[i].Start.Format(meta.DateTimeFormat), v[i].End.Format(meta.DateTimeFormat))
					}
				}
			}
		}
	}

	t.Log("Check for antenna installation mark overlaps")
	{
		installs := make(map[string]meta.InstalledAntennaList)
		for _, s := range antennas {
			if _, ok := installs[s.Mark]; !ok {
				installs[s.Mark] = meta.InstalledAntennaList{}
			}
			installs[s.Mark] = append(installs[s.Mark], s)
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
						t.Errorf("antennas %s [%s] and %s [%s] at %s has overlap between %s and %s",
							v[i].Model, v[i].Serial, v[j].Model, v[j].Serial, v[i].Mark,
							v[i].Start.Format(meta.DateTimeFormat), v[i].End.Format(meta.DateTimeFormat))
					}
				}
			}
		}
	}

	t.Log("Check for missing antenna marks")
	{
		var marks meta.MarkList

		if err := meta.LoadList("../network/marks.csv", &marks); err != nil {
			t.Fatal(err)
		}

		keys := make(map[string]interface{})

		for _, m := range marks {
			keys[m.Code] = true
		}

		for _, c := range antennas {
			if _, ok := keys[c.Mark]; ok {
				continue
			}
			t.Errorf("unable to find antenna mark %-5s", c.Mark)
		}
	}

	var assets meta.AssetList
	t.Log("Load antenna assets file")
	{
		if err := meta.LoadList("../assets/antennas.csv", &assets); err != nil {
			t.Fatal(err)
		}
	}

	t.Log("Check for antenna assets")
	{
		for _, r := range antennas {
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
				t.Errorf("unable to find antenna asset: %s [%s]", r.Model, r.Serial)
			}
		}
	}

	var sessions meta.SessionList
	t.Log("Load session list")
	{
		if err := meta.LoadList("../install/sessions.csv", &sessions); err != nil {
			t.Fatal(err)
		}
	}

	t.Log("Check sessions exist ...")
	{
		for _, r := range antennas {
			var found bool
			for _, s := range sessions {
				if s.End.Before(r.Start) {
					continue
				}
				if s.Start.After(r.End) {
					continue
				}
				found = true
			}
			if !found {
				t.Log(r)
			}
		}
	}
}
