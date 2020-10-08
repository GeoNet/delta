package delta_test

import (
	"testing"

	"github.com/GeoNet/delta/meta"
)

var testInstalledCameras = map[string]func([]meta.InstalledCamera) func(t *testing.T){

	"check for cameras installation equipment overlaps": func(cameras []meta.InstalledCamera) func(t *testing.T) {
		return func(t *testing.T) {

			installs := make(map[string]meta.InstalledCameraList)
			for _, c := range cameras {
				installs[c.Model] = append(installs[c.Model], c)
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
		}
	},
}

var testInstalledCamerasMounts = map[string]func([]meta.InstalledCamera, []meta.Mount) func(t *testing.T){

	"check for cameras installation equipment overlaps": func(cameras []meta.InstalledCamera, mounts []meta.Mount) func(t *testing.T) {
		return func(t *testing.T) {
			keys := make(map[string]interface{})
			for _, m := range mounts {
				keys[m.Code] = true
			}

			for _, c := range cameras {
				if _, ok := keys[c.Mount]; !ok {
					t.Errorf("unable to find camera mount %-5s", c.Mount)
				}
			}
		}
	},
}

var testInstalledCamerasAssets = map[string]func([]meta.InstalledCamera, []meta.Asset) func(t *testing.T){
	"check cameras assets": func(cameras []meta.InstalledCamera, assets []meta.Asset) func(t *testing.T) {
		return func(t *testing.T) {

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
				if found {
					continue
				}
				t.Errorf("unable to find camera asset: %s [%s]", r.Model, r.Serial)
			}
		}
	},
}

func TestCameras(t *testing.T) {

	var cameras meta.InstalledCameraList
	loadListFile(t, "../install/cameras.csv", &cameras)

	for k, fn := range testInstalledCameras {
		t.Run(k, fn(cameras))
	}
}

func TestCameras_Mounts(t *testing.T) {

	var cameras meta.InstalledCameraList
	loadListFile(t, "../install/cameras.csv", &cameras)

	var mounts meta.MountList
	loadListFile(t, "../network/mounts.csv", &mounts)

	for k, fn := range testInstalledCamerasMounts {
		t.Run(k, fn(cameras, mounts))
	}
}

func TestCameras_Assets(t *testing.T) {

	var cameras meta.InstalledCameraList
	loadListFile(t, "../install/cameras.csv", &cameras)

	var assets meta.AssetList
	loadListFile(t, "../assets/cameras.csv", &assets)

	for k, fn := range testInstalledCamerasAssets {
		t.Run(k, fn(cameras, assets))
	}
}
