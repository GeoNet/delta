package delta_test

import (
	"testing"

	"github.com/GeoNet/delta/meta"
)

func TestMetSensors(t *testing.T) {

	var installedMetsensors meta.InstalledMetSensorList
	loadListFile(t, "../install/metsensors.csv", &installedMetsensors)

	t.Run("check for metsensors installation equipment overlaps", func(t *testing.T) {
		installs := make(map[string]meta.InstalledMetSensorList)
		for _, s := range installedMetsensors {
			if _, ok := installs[s.Model]; !ok {
				installs[s.Model] = meta.InstalledMetSensorList{}
			}
			installs[s.Model] = append(installs[s.Model], s)
		}
		for _, v := range installs {
			for i := 0; i < len(v); i++ {
				for j := i + 1; j < len(v); j++ {
					if v[i].Serial != v[j].Serial {
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

					t.Errorf("metsensors %s at %-5s has mark %s overlap between %s and %s",
						v[i].Model, v[i].Serial, v[i].Mark,
						v[i].Start.Format(meta.DateTimeFormat),
						v[i].End.Format(meta.DateTimeFormat))
				}
			}
		}
	})

	t.Run("Check for missing metsensor marks", func(t *testing.T) {
		var marks meta.MarkList
		loadListFile(t, "../network/marks.csv", &marks)

		keys := make(map[string]interface{})
		for _, m := range marks {
			keys[m.Code] = true
		}

		for _, c := range installedMetsensors {
			if _, ok := keys[c.Mark]; !ok {
				t.Errorf("unable to find metsensor mark %-5s", c.Mark)
			}
		}
	})

	t.Run("check for missing metsensor assets", func(t *testing.T) {
		var assets meta.AssetList
		loadListFile(t, "../assets/metsensors.csv", &assets)

		for _, r := range installedMetsensors {
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
			if found {
				continue
			}
			t.Errorf("unable to find metsensor asset: %s [%s]", r.Model, r.Serial)
		}
	})
}
