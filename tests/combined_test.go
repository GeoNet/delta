package delta_test

import (
	"sort"
	"testing"

	"github.com/GeoNet/delta/meta"
)

func TestCombined(t *testing.T) {

	var combined meta.InstalledSensorList
	t.Log("Load installed sensors file")
	{
		if err := meta.LoadList("../install/sensors.csv", &combined); err != nil {
			t.Fatal(err)
		}
	}

	var recorders meta.InstalledRecorderList
	t.Log("Load installed recorders file")
	{
		if err := meta.LoadList("../install/recorders.csv", &recorders); err != nil {
			t.Fatal(err)
		}
	}

	for _, r := range recorders {
		combined = append(combined, meta.InstalledSensor{
			Install:  r.Install,
			Station:  r.Station,
			Location: r.Location,
		})
	}

	t.Log("Check for sensor/recorder installation location overlaps")
	{
		installs := make(map[string]meta.InstalledSensorList)
		for _, s := range combined {
			_, ok := installs[s.Station]
			if ok {
				installs[s.Station] = append(installs[s.Station], s)

			} else {
				installs[s.Station] = meta.InstalledSensorList{s}
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
					case v[i].Location != v[j].Location:
					case v[i].End.Before(v[j].Start):
					case v[i].Start.After(v[j].End):
					case v[i].End.Equal(v[j].Start):
					case v[i].Start.Equal(v[j].End):
					default:
						t.Errorf("sensor/recorder %s/%s at %-5s has location %-2s overlap between %s and %s",
							v[i].Model, v[i].Serial, v[i].Station, v[i].Location, v[i].Start.Format(meta.DateTimeFormat), v[i].End.Format(meta.DateTimeFormat))
					}
				}
			}
		}
	}

}
