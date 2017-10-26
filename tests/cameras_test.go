package delta_test

import (
	"testing"

	"github.com/GeoNet/delta/meta"
)

func TestCameras(t *testing.T) {

	var installedCameras meta.InstalledCameraList
	loadListFile(t, "../install/cameras.csv", &installedCameras)

	t.Run("Check for cameras installation equipment overlaps", func(t *testing.T) {
		installs := make(map[string]meta.InstalledCameraList)
		for _, s := range installedCameras {
			if _, ok := installs[s.Model]; !ok {
				installs[s.Model] = meta.InstalledCameraList{}
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

					t.Errorf("cameras %s at %-5s has mount %s overlap between %s and %s",
						v[i].Model, v[i].Serial, v[i].Mount,
						v[i].Start.Format(meta.DateTimeFormat),
						v[i].End.Format(meta.DateTimeFormat))
				}
			}
		}
	})

	t.Run("Check for missing camera mounts", func(t *testing.T) {
		var mounts meta.MountList
		loadListFile(t, "../network/mounts.csv", &mounts)

		keys := make(map[string]interface{})
		for _, m := range mounts {
			keys[m.Code] = true
		}

		for _, c := range installedCameras {
			if _, ok := keys[c.Mount]; !ok {
				t.Errorf("unable to find camera mount %-5s", c.Mount)
			}
		}
	})

	t.Run("Load camera assets file", func(t *testing.T) {
		var assets meta.AssetList
		loadListFile(t, "../assets/cameras.csv", &assets)

		for _, r := range installedCameras {
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
			t.Errorf("unable to find camera asset: %s [%s]", r.Model, r.Serial)
		}
	})

}
