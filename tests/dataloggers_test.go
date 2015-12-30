package meta_test

import (
	"sort"
	"testing"

	"github.com/GeoNet/delta/meta"
)

func TestDataloggers(t *testing.T) {
	var dataloggers meta.DeployedDataloggers

	t.Log("Load deployed dataloggers file")
	{
		if err := meta.LoadLists("../equipment/dataloggers", "dataloggers.csv", &dataloggers); err != nil {
			t.Fatal(err)
		}
	}

	t.Log("Check for datalogger installation place overlaps")
	{
		installs := make(map[string]meta.DeployedDataloggers)
		for _, d := range dataloggers {
			_, ok := installs[d.Place]
			if ok {
				installs[d.Place] = append(installs[d.Place], d)

			} else {
				installs[d.Place] = meta.DeployedDataloggers{d}
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
					case v[i].EndTime.Before(v[j].StartTime):
					case v[i].StartTime.After(v[j].EndTime):
					//case v[i].EndTime.Equal(v[j].StartTime):
					//case v[i].StartTime.Equal(v[j].EndTime):
					default:
						t.Errorf("datalogger %s:[%s] at %-32s has place overlap between %s and %s",
							v[i].Model, v[i].Serial, v[i].Place, v[i].StartTime.Format(meta.DateTimeFormat), v[i].EndTime.Format(meta.DateTimeFormat))
					}
				}
			}
		}
	}

	t.Log("Check for datalogger installation equipment overlaps")
	{
		installs := make(map[string]meta.DeployedDataloggers)
		for _, s := range dataloggers {
			_, ok := installs[s.Model]
			if ok {
				installs[s.Model] = append(installs[s.Model], s)

			} else {
				installs[s.Model] = meta.DeployedDataloggers{s}
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
					case v[i].EndTime.Before(v[j].StartTime):
					case v[i].StartTime.After(v[j].EndTime):
						//		case v[i].EndTime.Equal(v[j].StartTime):
						//		case v[i].StartTime.Equal(v[j].EndTime):
					default:
						t.Errorf("datalogger %s:[%s] at %-32s has installation overlap between %s and %s",
							v[i].Model, v[i].Serial, v[i].Place, v[i].StartTime.Format(meta.DateTimeFormat), v[i].EndTime.Format(meta.DateTimeFormat))
					}
				}
			}
		}
	}

}
