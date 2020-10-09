package delta_test

import (
	"testing"

	"github.com/GeoNet/delta/meta"
)

var testInstalledMetSensors = map[string]func([]meta.InstalledMetSensor) func(t *testing.T){

	"check for metsensors installation equipment overlaps": func(installedMetSensors []meta.InstalledMetSensor) func(t *testing.T) {
		return func(t *testing.T) {

			installs := make(map[string]meta.InstalledMetSensorList)
			for _, s := range installedMetSensors {
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
		}
	},
}

var testInstalledMetSensorsMarks = map[string]func([]meta.InstalledMetSensor, []meta.Mark) func(t *testing.T){

	"check for missing metsensor marks": func(installedMetSensors []meta.InstalledMetSensor, marks []meta.Mark) func(t *testing.T) {
		return func(t *testing.T) {
			keys := make(map[string]interface{})
			for _, m := range marks {
				keys[m.Code] = true
			}

			for _, c := range installedMetSensors {
				if _, ok := keys[c.Mark]; !ok {
					t.Errorf("unable to find metsensor mark %-5s", c.Mark)
				}
			}
		}
	},
}

var testInstalledMetSensorsAssets = map[string]func([]meta.InstalledMetSensor, []meta.Asset) func(t *testing.T){
	"check for missing metsensor assets": func(installedMetSensors []meta.InstalledMetSensor, assets []meta.Asset) func(t *testing.T) {
		return func(t *testing.T) {

			for _, r := range installedMetSensors {
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
		}
	},
}

func TestInstalledMetSensors(t *testing.T) {

	var installedMetSensors meta.InstalledMetSensorList
	loadListFile(t, "../install/metsensors.csv", &installedMetSensors)

	for k, fn := range testInstalledMetSensors {
		t.Run(k, fn(installedMetSensors))
	}
}

func TestInstalledMetSensors_Marks(t *testing.T) {

	var installedMetSensors meta.InstalledMetSensorList
	loadListFile(t, "../install/metsensors.csv", &installedMetSensors)

	var marks meta.MarkList
	loadListFile(t, "../network/marks.csv", &marks)

	for k, fn := range testInstalledMetSensorsMarks {
		t.Run(k, fn(installedMetSensors, marks))
	}
}

func TestInstalledMetSensors_Assets(t *testing.T) {

	var installedMetSensors meta.InstalledMetSensorList
	loadListFile(t, "../install/metsensors.csv", &installedMetSensors)

	var assets meta.AssetList
	loadListFile(t, "../assets/metsensors.csv", &assets)

	for k, fn := range testInstalledMetSensorsAssets {
		t.Run(k, fn(installedMetSensors, assets))
	}
}
