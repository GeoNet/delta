package delta_test

import (
	"testing"

	"github.com/GeoNet/delta/meta"
)

var testInstalledDoases = map[string]func([]meta.InstalledDoas) func(t *testing.T){

	"check for doases installation equipment overlaps": func(doases []meta.InstalledDoas) func(t *testing.T) {
		return func(t *testing.T) {

			installs := make(map[string]meta.InstalledDoasList)
			for _, c := range doases {
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

						t.Errorf("doases %s at %-5s has mount %s overlap between %s and %s",
							v[i].Model, v[i].Serial, v[i].Mount,
							v[i].Start.Format(meta.DateTimeFormat),
							v[i].End.Format(meta.DateTimeFormat))
					}
				}
			}
		}
	},
}

var testInstalledDoasesMounts = map[string]func([]meta.InstalledDoas, []meta.Mount) func(t *testing.T){

	"check for doases installation equipment overlaps": func(doases []meta.InstalledDoas, mounts []meta.Mount) func(t *testing.T) {
		return func(t *testing.T) {
			keys := make(map[string]interface{})
			for _, m := range mounts {
				keys[m.Code] = true
			}

			for _, c := range doases {
				if _, ok := keys[c.Mount]; !ok {
					t.Errorf("unable to find doas mount %-5s", c.Mount)
				}
			}
		}
	},
}

var testInstalledDoasesViews = map[string]func([]meta.InstalledDoas, []meta.View) func(t *testing.T){

	"check for doases installation views": func(doases []meta.InstalledDoas, views []meta.View) func(t *testing.T) {
		return func(t *testing.T) {
			type view struct{ m, c string }
			keys := make(map[view]interface{})
			for _, m := range views {
				keys[view{m.Mount, m.Code}] = true
			}

			for _, c := range doases {
				if _, ok := keys[view{c.Mount, c.View}]; !ok {
					t.Errorf("unable to find doas mount %-5s (%-2s)", c.Mount, c.View)
				}
			}
		}
	},
}

var testInstalledDoasesAssets = map[string]func([]meta.InstalledDoas, []meta.Asset) func(t *testing.T){
	"check doases assets": func(doases []meta.InstalledDoas, assets []meta.Asset) func(t *testing.T) {
		return func(t *testing.T) {

			for _, r := range doases {
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
				t.Errorf("unable to find doas asset: %s [%s]", r.Model, r.Serial)
			}
		}
	},
}

func TestDoases(t *testing.T) {

	var doases meta.InstalledDoasList
	loadListFile(t, "../install/doases.csv", &doases)

	for k, fn := range testInstalledDoases {
		t.Run(k, fn(doases))
	}
}

func TestDoases_Mounts(t *testing.T) {

	var doases meta.InstalledDoasList
	loadListFile(t, "../install/doases.csv", &doases)

	var mounts meta.MountList
	loadListFile(t, "../network/mounts.csv", &mounts)

	for k, fn := range testInstalledDoasesMounts {
		t.Run(k, fn(doases, mounts))
	}
}

func TestDoases_Views(t *testing.T) {

	var doases meta.InstalledDoasList
	loadListFile(t, "../install/doases.csv", &doases)

	var views meta.ViewList
	loadListFile(t, "../network/views.csv", &views)

	for k, fn := range testInstalledDoasesViews {
		t.Run(k, fn(doases, views))
	}
}

func TestDoases_Assets(t *testing.T) {

	var doases meta.InstalledDoasList
	loadListFile(t, "../install/doases.csv", &doases)

	var assets meta.AssetList
	loadListFile(t, "../assets/doases.csv", &assets)

	for k, fn := range testInstalledDoasesAssets {
		t.Run(k, fn(doases, assets))
	}
}
