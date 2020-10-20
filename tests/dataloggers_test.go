package delta_test

import (
	"testing"

	"github.com/GeoNet/delta/meta"
)

var testDeployedDataloggers = map[string]func([]meta.DeployedDatalogger) func(t *testing.T){
	"check for datalogger installation place overlaps": func(dataloggers []meta.DeployedDatalogger) func(t *testing.T) {
		return func(t *testing.T) {
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
		}
	},
	"check for datalogger installation equipment overlaps": func(dataloggers []meta.DeployedDatalogger) func(t *testing.T) {
		return func(t *testing.T) {
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
		}
	},
}

var testDeployedDataloggersAssets = map[string]func([]meta.DeployedDatalogger, []meta.Asset) func(t *testing.T){
	"check for missing datalogger assets": func(dataloggers []meta.DeployedDatalogger, assets []meta.Asset) func(t *testing.T) {
		return func(t *testing.T) {
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
		}
	},
}

func TestDataloggers(t *testing.T) {
	var dataloggers meta.DeployedDataloggerList
	loadListFile(t, "../install/dataloggers.csv", &dataloggers)

	for k, fn := range testDeployedDataloggers {
		t.Run(k, fn(dataloggers))
	}

}

func TestDataloggers_Assets(t *testing.T) {
	var dataloggers meta.DeployedDataloggerList
	loadListFile(t, "../install/dataloggers.csv", &dataloggers)

	var assets meta.AssetList
	loadListFile(t, "../assets/dataloggers.csv", &assets)

	for k, fn := range testDeployedDataloggersAssets {
		t.Run(k, fn(dataloggers, assets))
	}
}
