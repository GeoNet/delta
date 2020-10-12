package delta_test

import (
	"testing"

	"github.com/GeoNet/delta/meta"
)

var testInstalledRadomes = map[string]func([]meta.InstalledRadome) func(t *testing.T){

	"check for radomes installation equipment overlaps": func(radomes []meta.InstalledRadome) func(t *testing.T) {
		return func(t *testing.T) {
			installs := make(map[string]meta.InstalledRadomeList)
			for _, s := range radomes {
				if _, ok := installs[s.Model]; !ok {
					installs[s.Model] = meta.InstalledRadomeList{}
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
						t.Errorf("radomes %s at %-5s has mark %s overlap between %s and %s",
							v[i].Model, v[i].Serial, v[i].Mark, v[i].Start.Format(meta.DateTimeFormat), v[i].End.Format(meta.DateTimeFormat))
					}
				}
			}
		}
	},

	"check for overlapping radomes installations": func(radomes []meta.InstalledRadome) func(t *testing.T) {
		return func(t *testing.T) {
			installs := make(map[string]meta.InstalledRadomeList)
			for _, s := range radomes {
				if _, ok := installs[s.Mark]; !ok {
					installs[s.Mark] = meta.InstalledRadomeList{}
				}
				installs[s.Mark] = append(installs[s.Mark], s)
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

						t.Errorf("mark %-5s has radome %s/%s overlap wth %s/%s between %s and %s",
							v[i].Mark, v[i].Model, v[i].Serial, v[j].Model, v[j].Serial, v[i].Start.Format(meta.DateTimeFormat), v[i].End.Format(meta.DateTimeFormat))
					}
				}
			}
		}
	},
}

var testInstalledRadomesMarks = map[string]func([]meta.InstalledRadome, []meta.Mark) func(t *testing.T){

	"check for missing radome marks": func(radomes []meta.InstalledRadome, marks []meta.Mark) func(t *testing.T) {
		return func(t *testing.T) {

			keys := make(map[string]interface{})
			for _, m := range marks {
				keys[m.Code] = true
			}

			for _, c := range radomes {
				if _, ok := keys[c.Mark]; !ok {
					t.Errorf("unable to find radome mark %-5s", c.Mark)
				}
			}
		}
	},
}

var testInstalledRadomesAssets = map[string]func([]meta.InstalledRadome, []meta.Asset) func(t *testing.T){
	"check for missing radome assets": func(radomes []meta.InstalledRadome, assets []meta.Asset) func(t *testing.T) {
		return func(t *testing.T) {

			for _, r := range radomes {
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
				t.Errorf("unable to find radome asset: %s [%s]", r.Model, r.Serial)
			}
		}
	},
}

func TestInstalledRadomes(t *testing.T) {

	var radomes meta.InstalledRadomeList
	loadListFile(t, "../install/radomes.csv", &radomes)

	for k, fn := range testInstalledRadomes {
		t.Run(k, fn(radomes))
	}
}

func TestInstalledRadomes_Marks(t *testing.T) {

	var radomes meta.InstalledRadomeList
	loadListFile(t, "../install/radomes.csv", &radomes)

	var marks meta.MarkList
	loadListFile(t, "../network/marks.csv", &marks)

	for k, fn := range testInstalledRadomesMarks {
		t.Run(k, fn(radomes, marks))
	}
}

func TestInstalledRadomes_Assets(t *testing.T) {

	var radomes meta.InstalledRadomeList
	loadListFile(t, "../install/radomes.csv", &radomes)

	var assets meta.AssetList
	loadListFile(t, "../assets/radomes.csv", &assets)

	for k, fn := range testInstalledRadomesAssets {
		t.Run(k, fn(radomes, assets))
	}
}
