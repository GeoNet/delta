package delta_test

import (
	"sort"
	"testing"

	"github.com/GeoNet/delta/meta"
)

func TestRecorders(t *testing.T) {
	var recorders meta.InstalledRecorderList

	t.Log("Load installed recorders file")
	{
		if err := meta.LoadList("../install/recorders.csv", &recorders); err != nil {
			t.Fatal(err)
		}
	}

	t.Log("Check for recorder installation equipment overlaps")
	{
		installs := make(map[string]meta.InstalledRecorderList)
		for _, s := range recorders {
			_, ok := installs[s.Model]
			if ok {
				installs[s.Model] = append(installs[s.Model], s)

			} else {
				installs[s.Model] = meta.InstalledRecorderList{s}
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
					case v[i].End.Equal(v[j].Start):
					case v[i].Start.Equal(v[j].End):
					default:
						t.Errorf("recorder %s/%s at %-5s has location %-2s overlap between %s and %s",
							v[i].Model, v[i].Serial, v[i].StationCode, v[i].LocationCode, v[i].Start.Format(meta.DateTimeFormat), v[i].End.Format(meta.DateTimeFormat))
					}
				}
			}
		}
	}

	t.Log("Check for missing recorder stations")
	{
		var stations meta.StationList

		if err := meta.LoadList("../network/stations.csv", &stations); err != nil {
			t.Fatal(err)
		}

		keys := make(map[string]interface{})

		for _, s := range stations {
			keys[s.Code] = true
		}

		for _, s := range recorders {
			if _, ok := keys[s.StationCode]; ok {
				continue
			}
			t.Errorf("unable to find recorder installed station %-5s", s.StationCode)
		}
	}

}
