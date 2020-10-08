package delta_test

import (
	"testing"

	"github.com/GeoNet/delta/meta"
)

var testInstalledAntennas = map[string]func([]meta.InstalledAntenna) func(t *testing.T){

	// check for session overlaps, there can't be two sessions running at the same mark for the same sampling interval.
	"check antenna installation overlap": func(installedAntennas []meta.InstalledAntenna) func(t *testing.T) {
		return func(t *testing.T) {

			installs := make(map[string]meta.InstalledAntennaList)
			for _, s := range installedAntennas {
				if _, ok := installs[s.Model]; !ok {
					installs[s.Model] = meta.InstalledAntennaList{}
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
						t.Errorf("installed antennas %s [%s] at %s has overlap at %s between %s and %s",
							v[i].Model, v[i].Serial, v[i].Mark, v[j].Mark,
							v[i].Start.Format(meta.DateTimeFormat),
							v[i].End.Format(meta.DateTimeFormat))
					}
				}
			}
		}
	},

	"check for antenna installation mark overlaps": func(installedAntennas []meta.InstalledAntenna) func(t *testing.T) {
		return func(t *testing.T) {
			installs := make(map[string]meta.InstalledAntennaList)
			for _, s := range installedAntennas {
				if _, ok := installs[s.Mark]; !ok {
					installs[s.Mark] = meta.InstalledAntennaList{}
				}
				installs[s.Mark] = append(installs[s.Mark], s)
			}

			for _, v := range installs {
				for i := 0; i < len(v); i++ {
					for j := i + 1; j < len(v); j++ {
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
						t.Errorf("installed antennas %s [%s] and %s [%s] at %s has an overlap between %s and %s",
							v[i].Model, v[i].Serial, v[j].Model, v[j].Serial, v[i].Mark,
							v[i].Start.Format(meta.DateTimeFormat), v[i].End.Format(meta.DateTimeFormat))
					}
				}
			}
		}
	},
}

var testInstalledAntennasMarks = map[string]func([]meta.InstalledAntenna, []meta.Mark) func(t *testing.T){
	"check for missing antenna marks": func(installedAntennas []meta.InstalledAntenna, marks []meta.Mark) func(t *testing.T) {
		return func(t *testing.T) {

			keys := make(map[string]interface{})
			for _, m := range marks {
				keys[m.Code] = true
			}

			for _, c := range installedAntennas {
				if _, ok := keys[c.Mark]; !ok {
					t.Errorf("unable to find antenna mark %-5s", c.Mark)
				}
			}
		}
	},
}

var testInstalledAntennasAssets = map[string]func([]meta.InstalledAntenna, []meta.Asset) func(t *testing.T){
	"check for missing antenna assets": func(installedAntennas []meta.InstalledAntenna, assets []meta.Asset) func(t *testing.T) {
		return func(t *testing.T) {
			for _, r := range installedAntennas {
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
				t.Errorf("unable to find antenna asset: %s [%s]", r.Model, r.Serial)
			}
		}
	},
}

var testInstalledAntennasSessions = map[string]func([]meta.InstalledAntenna, []meta.Session) func(t *testing.T){
	"check for missing antenna sessions": func(installedAntennas []meta.InstalledAntenna, sessions []meta.Session) func(t *testing.T) {
		return func(t *testing.T) {

			for _, r := range installedAntennas {
				var found bool
				for _, s := range sessions {
					if s.End.Before(r.Start) {
						continue
					}
					if s.Start.After(r.End) {
						continue
					}
					found = true
				}
				if found {
					continue
				}
				t.Errorf("unable to find antenna session: %s [%s]", r.Model, r.Serial)
			}
		}
	},
}

func TestInstalledAntennas(t *testing.T) {

	var installedAntennas meta.InstalledAntennaList
	loadListFile(t, "../install/antennas.csv", &installedAntennas)

	for k, fn := range testInstalledAntennas {
		t.Run(k, fn(installedAntennas))
	}
}

func TestInstalledAntennas_Marks(t *testing.T) {

	var installedAntennas meta.InstalledAntennaList
	loadListFile(t, "../install/antennas.csv", &installedAntennas)

	var marks meta.MarkList
	loadListFile(t, "../network/marks.csv", &marks)

	for k, fn := range testInstalledAntennasMarks {
		t.Run(k, fn(installedAntennas, marks))
	}

}

func TestInstalledAntennas_Assets(t *testing.T) {

	var installedAntennas meta.InstalledAntennaList
	loadListFile(t, "../install/antennas.csv", &installedAntennas)

	var assets meta.AssetList
	loadListFile(t, "../assets/antennas.csv", &assets)

	for k, fn := range testInstalledAntennasAssets {
		t.Run(k, fn(installedAntennas, assets))
	}
}

func TestInstalledAntennas_Sessions(t *testing.T) {

	var installedAntennas meta.InstalledAntennaList
	loadListFile(t, "../install/antennas.csv", &installedAntennas)

	var sessions meta.SessionList
	loadListFile(t, "../install/sessions.csv", &sessions)

	for k, fn := range testInstalledAntennasSessions {
		t.Run(k, fn(installedAntennas, sessions))
	}
}
