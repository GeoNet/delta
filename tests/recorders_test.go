package delta_test

import (
	"testing"

	"github.com/GeoNet/delta/meta"
)

func TestRecorders(t *testing.T) {
	var recorders meta.InstalledRecorderList
	loadListFile(t, "../install/recorders.csv", &recorders)

	t.Run("check for recorder installation equipment overlaps", func(t *testing.T) {
		installs := make(map[string]meta.InstalledRecorderList)
		for _, s := range recorders {
			if _, ok := installs[s.Model]; !ok {
				installs[s.Model] = meta.InstalledRecorderList{}
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
					t.Errorf("recorder %s/%s at %-5s has location %-2s overlap between %s and %s",
						v[i].Model, v[i].Serial,
						v[i].Station, v[i].Location,
						v[i].Start.Format(meta.DateTimeFormat),
						v[i].End.Format(meta.DateTimeFormat))
				}
			}
		}
	})

	t.Run("check for missing recorder stations", func(t *testing.T) {
		var stations meta.StationList
		loadListFile(t, "../network/stations.csv", &stations)

		keys := make(map[string]interface{})
		for _, s := range stations {
			keys[s.Code] = true
		}

		for _, s := range recorders {
			if _, ok := keys[s.Station]; ok {
				continue
			}
			t.Errorf("unable to find recorder installed station %-5s", s.Station)
		}
	})

	t.Run("check for recorder assets", func(t *testing.T) {
		var assets meta.AssetList
		loadListFile(t, "../assets/recorders.csv", &assets)
		for _, r := range recorders {
			var found bool
			for _, a := range assets {
				if a.Model != r.DataloggerModel {
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
			t.Errorf("unable to find recorders asset: %s [%s]", r.DataloggerModel, r.Serial)
		}
	})

}
