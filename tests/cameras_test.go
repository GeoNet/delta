package delta_test

import (
	"testing"
	"time"

	"github.com/GeoNet/delta"
	"github.com/GeoNet/delta/meta"
)

var cameraChecks = map[string]func(*meta.Set) func(t *testing.T){

	"check for cameras installation equipment overlaps": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {

			installs := make(map[string]meta.InstalledCameraList)
			for _, c := range set.InstalledCameras() {
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

	"check for invalid installation dates": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			for _, i := range set.InstalledCameras() {
				if i.End.After(i.Start) {
					continue
				}
				t.Errorf("installed camera is removed before it has been installed: %s", i.String())
			}
		}
	},

	"check for cameras installation mount overlaps": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			keys := make(map[string]interface{})
			for _, m := range set.Mounts() {
				keys[m.Code] = true
			}

			for _, c := range set.InstalledCameras() {
				if _, ok := keys[c.Mount]; !ok {
					t.Errorf("unable to find camera mount %-5s", c.Mount)
				}
			}
		}
	},

	"check for cameras installation views": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			type view struct{ m, c string }
			keys := make(map[view][]meta.View)
			for _, m := range set.Views() {
				keys[view{m.Mount, m.Code}] = append(keys[view{m.Mount, m.Code}], m)
			}

			for _, c := range set.InstalledCameras() {
				if _, ok := keys[view{c.Mount, c.View}]; !ok {
					t.Errorf("unable to find camera mount %-5s (%-2s)", c.Mount, c.View)
				}
			}

			for _, c := range set.InstalledCameras() {
				views, ok := keys[view{c.Mount, c.View}]
				if !ok {
					continue
				}
				var found bool
				for _, v := range views {
					if c.Start.Before(v.Start) {
						continue
					}
					if c.End.After(v.End) {
						continue
					}
					found = true
				}
				if found {
					continue
				}
				t.Errorf("installed camera is not within a known view: %s (%s/%s) %s", c.String(), c.Mount, c.View, c.Start.Format(time.RFC3339))
			}
		}
	},

	"check cameras assets": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			for _, r := range set.InstalledCameras() {
				var found bool
				for _, a := range set.Assets() {
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

	set, err := delta.New()
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range cameraChecks {
		t.Run(k, v(set))
	}
}
