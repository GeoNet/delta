package delta_test

import (
	"testing"

	"github.com/GeoNet/delta/meta"
)

func TestCombined(t *testing.T) {

	var combined meta.InstalledSensorList
	loadListFile(t, "../install/sensors.csv", &combined)

	var recorders meta.InstalledRecorderList
	loadListFile(t, "../install/recorders.csv", &recorders)

	for _, r := range recorders {
		combined = append(combined, meta.InstalledSensor{
			Install:  r.Install,
			Station:  r.Station,
			Location: r.Location,
		})
	}

	t.Run("check for sensor/recorder installation location overlaps", func(t *testing.T) {
		installs := make(map[string]meta.InstalledSensorList)
		for _, s := range combined {
			if _, ok := installs[s.Station]; !ok {
				installs[s.Station] = meta.InstalledSensorList{}
			}
			installs[s.Station] = append(installs[s.Station], s)
		}
		for _, v := range installs {
			for i := 0; i < len(v); i++ {
				for j := i + 1; j < len(v); j++ {
					if v[i].Location != v[j].Location {
						continue
					}
					if v[i].End.Before(v[j].Start) {
						continue
					}
					if v[i].Start.After(v[j].End) {
						continue
					}
					if v[i].End.Equal(v[j].Start) {
						continue
					}
					if v[i].Start.Equal(v[j].End) {
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
