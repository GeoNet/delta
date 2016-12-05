package delta_test

import (
	"sort"
	"testing"

	"github.com/GeoNet/delta/meta"
)

func testCombined(t *testing.T) {
	var sensors meta.InstalledSensorList

	if err := meta.LoadList("../install/sensors.csv", &sensors); err != nil {
		t.Fatal(err)
	}

	var recorders meta.InstalledRecorderList

	if err := meta.LoadList("../install/recorders.csv", &recorders); err != nil {
		t.Fatal(err)
	}

	t.Run("Check for sensor/recorder installation location overlaps", func(t *testing.T) {
		combined := append(meta.InstalledSensorList{}, sensors...)

		for _, r := range recorders {
			combined = append(combined, meta.InstalledSensor{
				Install:  r.Install,
				Station:  r.Station,
				Location: r.Location,
			})
		}

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
			for i, v, n := 0, installs[k], len(installs[k]); i < n; i++ {
				for j := i + 1; j < n; j++ {
					if v[i].Location != v[j].Location {
						continue
					}
					if v[i].End.Before(v[j].Start) || v[i].Start.After(v[j].End) {
						continue
					}
					if v[i].End.Equal(v[j].Start) || v[i].Start.Equal(v[j].End) {
						continue
					}
					t.Errorf("sensor/recorder %s/%s at %-5s has location %-2s overlap between %s and %s",
						v[i].Model, v[i].Serial,
						v[i].Station, v[i].Location,
						v[i].Start.Format(meta.DateTimeFormat),
						v[i].End.Format(meta.DateTimeFormat))
				}
			}
		}
	})
}
