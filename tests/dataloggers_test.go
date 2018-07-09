package delta_test

import (
	"testing"

	"github.com/GeoNet/delta/meta"
)

func TestDataloggers(t *testing.T) {
	var dataloggers meta.DeployedDataloggerList

	loadListFile(t, "../install/dataloggers.csv", &dataloggers)

	t.Run("check for datalogger installation place overlaps", func(t *testing.T) {
		installs := make(map[string]meta.DeployedDataloggerList)
		for _, d := range dataloggers {
			if _, ok := installs[d.Place]; !ok {
				installs[d.Place] = meta.DeployedDataloggerList{}
			}
			installs[d.Place] = append(installs[d.Place], d)
		}

		for _, v := range installs {
			for i := 0; i < len(v); i++ {
				for j := i + 1; j < len(v); j++ {
					if v[i].Place != v[j].Place {
						continue
					}
					if v[i].Role != v[j].Role {
						continue
					}
					if v[i].End.Before(v[j].Start) {
						continue
					}
					if v[i].Start.After(v[j].End) {
						continue
					}

					t.Errorf("datalogger %s:[%s] at %-32s has place overlap between %s and %s",
						v[i].Model, v[i].Serial, v[i].Place,
						v[i].Start.Format(meta.DateTimeFormat),
						v[i].End.Format(meta.DateTimeFormat))
				}
			}
		}
	})

	t.Run("check for datalogger installation equipment overlaps", func(t *testing.T) {
		installs := make(map[string]meta.DeployedDataloggerList)
		for _, s := range dataloggers {
			if _, ok := installs[s.Model]; !ok {
				installs[s.Model] = meta.DeployedDataloggerList{}
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

					t.Errorf("datalogger %s:[%s] at %-32s has installation overlap between %s and %s",
						v[i].Model, v[i].Serial, v[i].Place,
						v[i].Start.Format(meta.DateTimeFormat),
						v[i].End.Format(meta.DateTimeFormat))
				}
			}
		}
	})

	t.Run("check for missing datalogger assets", func(t *testing.T) {
		var assets meta.AssetList
		loadListFile(t, "../assets/dataloggers.csv", &assets)
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
			if found {
				continue
			}
			t.Errorf("unable to find datalogger asset: %s [%s]", r.Model, r.Serial)
		}
	})
}
