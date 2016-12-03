package delta_test

import (
	"sort"
	"testing"

	"github.com/GeoNet/delta/meta"
)

func TestCameras(t *testing.T) {
	var cameras meta.InstalledCameraList

	if err := meta.LoadList("../install/cameras.csv", &cameras); err != nil {
		t.Fatal(err)
	}

	var mounts meta.MountList

	if err := meta.LoadList("../network/mounts.csv", &mounts); err != nil {
		t.Fatal(err)
	}

	var assets meta.AssetList

	if err := meta.LoadList("../assets/cameras.csv", &assets); err != nil {
		t.Fatal(err)
	}

	t.Run("Check for cameras installation equipment overlaps", func(t *testing.T) {
		installs := make(map[string]meta.InstalledCameraList)
		for _, s := range cameras {
			_, ok := installs[s.Model]
			if ok {
				installs[s.Model] = append(installs[s.Model], s)

			} else {
				installs[s.Model] = meta.InstalledCameraList{s}
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
					if v[i].Serial != v[j].Serial {
						continue
					}
					if v[i].End.Before(v[j].Start) || v[i].Start.After(v[j].End) {
						continue
					}
					if v[i].End.Equal(v[j].Start) || v[i].Start.Equal(v[j].End) {
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
		keys := make(map[string]interface{})

		for _, m := range mounts {
			keys[m.Code] = true
		}

		for _, c := range cameras {
			if _, ok := keys[c.Mount]; ok {
				continue
			}
			t.Errorf("unable to find camera mount %-5s", c.Mount)
		}
	})

	t.Run("Check for camera assets", func(t *testing.T) {
		for _, r := range cameras {
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
			if !found {
				t.Errorf("unable to find camera asset: %s [%s]", r.Model, r.Serial)
			}
		}
	})

}
