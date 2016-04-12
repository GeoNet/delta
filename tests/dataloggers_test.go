package delta_test

import (
	"sort"
	"testing"

	"github.com/GeoNet/delta/meta"
)

func TestDataloggers(t *testing.T) {
	var dataloggers meta.DeployedDataloggerList

	t.Log("Load deployed dataloggers file")
	{
		if err := meta.LoadList("../install/dataloggers.csv", &dataloggers); err != nil {
			t.Fatal(err)
		}
	}

	t.Log("Check for datalogger installation place overlaps")
	{
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
					switch {
					case v[i].Place != v[j].Place:
					case v[i].Role != v[j].Role:
					//case v[i].Model != v[j].Model:
					case v[i].End.Before(v[j].Start):
					case v[i].Start.After(v[j].End):
					//case v[i].End.Equal(v[j].Start):
					//case v[i].Start.Equal(v[j].End):
					default:
						t.Errorf("datalogger %s:[%s] at %-32s has place overlap between %s and %s",
							v[i].Model, v[i].Serial, v[i].Place, v[i].Start.Format(meta.DateTimeFormat), v[i].End.Format(meta.DateTimeFormat))
					}
				}
			}
		}
	}

	t.Log("Check for datalogger installation equipment overlaps")
	{
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
			v := installs[k]

			for i, n := 0, len(v); i < n; i++ {
				for j := i + 1; j < n; j++ {
					switch {
					case v[i].Serial != v[j].Serial:
					case v[i].End.Before(v[j].Start):
					case v[i].Start.After(v[j].End):
						//		case v[i].End.Equal(v[j].Start):
						//		case v[i].Start.Equal(v[j].End):
					default:
						t.Errorf("datalogger %s:[%s] at %-32s has installation overlap between %s and %s",
							v[i].Model, v[i].Serial, v[i].Place, v[i].Start.Format(meta.DateTimeFormat), v[i].End.Format(meta.DateTimeFormat))
					}
				}
			}
		}
	}

	var assets meta.AssetList
	t.Log("Load datalogger assets file")
	{
		if err := meta.LoadList("../assets/dataloggers.csv", &assets); err != nil {
			t.Fatal(err)
		}
	}

	t.Log("Check for datalogger assets")
	{
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
	}

}
