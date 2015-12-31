package delta_test

import (
	"sort"
	"testing"

	"github.com/GeoNet/delta/meta"
)

func TestSensors(t *testing.T) {
	var sensors meta.InstalledSensorList

	t.Log("Load installed sensors file")
	{
		if err := meta.LoadList("../installs/sensors.csv", &sensors); err != nil {
			t.Fatal(err)
		}
	}

	t.Log("Check for sensor installation location overlaps")
	{
		installs := make(map[string]meta.InstalledSensorList)
		for _, s := range sensors {
			_, ok := installs[s.StationCode]
			if ok {
				installs[s.StationCode] = append(installs[s.StationCode], s)

			} else {
				installs[s.StationCode] = meta.InstalledSensorList{s}
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
					case v[i].LocationCode != v[j].LocationCode:
					case v[i].End.Before(v[j].Start):
					case v[i].Start.After(v[j].End):
					case v[i].End.Equal(v[j].Start):
					case v[i].Start.Equal(v[j].End):
					default:
						t.Errorf("sensor %s/%s at %-5s has location %-2s overlap between %s and %s",
							v[i].Model, v[i].Serial, v[i].StationCode, v[i].LocationCode, v[i].Start.Format(meta.DateTimeFormat), v[i].End.Format(meta.DateTimeFormat))
					}
				}
			}
		}
	}

	t.Log("Check for sensor installation equipment overlaps")
	{
		installs := make(map[string]meta.InstalledSensorList)
		for _, s := range sensors {
			_, ok := installs[s.Model]
			if ok {
				installs[s.Model] = append(installs[s.Model], s)

			} else {
				installs[s.Model] = meta.InstalledSensorList{s}
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
						t.Errorf("sensor %s/%s at %-5s has location %-2s overlap between %s and %s",
							v[i].Model, v[i].Serial, v[i].StationCode, v[i].LocationCode, v[i].Start.Format(meta.DateTimeFormat), v[i].End.Format(meta.DateTimeFormat))
					}
				}
			}
		}
	}

	t.Log("Check for missing sensor stations")
	{
		var stations meta.StationList

		if err := meta.LoadList("../network/stations.csv", &stations); err != nil {
			t.Fatal(err)
		}

		keys := make(map[string]interface{})

		for _, s := range stations {
			keys[s.Code] = true
		}

		for _, s := range sensors {
			if _, ok := keys[s.StationCode]; ok {
				continue
			}
			t.Errorf("unable to find sensor installed station %-5s", s.StationCode)
		}
	}

}
