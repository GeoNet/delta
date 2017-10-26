package delta_test

import (
	"testing"

	"github.com/GeoNet/delta/meta"
)

func TestAntennas(t *testing.T) {

	var installedAntennas meta.InstalledAntennaList
	loadListFile(t, "../install/antennas.csv", &installedAntennas)

	t.Run("check for antenna installation equipment overlaps", func(t *testing.T) {
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
	})

	t.Run("check for antenna installation mark overlaps", func(t *testing.T) {
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
	})

	t.Run("check for missing antenna marks", func(t *testing.T) {
		var marks meta.MarkList
		loadListFile(t, "../network/marks.csv", &marks)

		keys := make(map[string]interface{})
		for _, m := range marks {
			keys[m.Code] = true
		}

		for _, c := range installedAntennas {
			if _, ok := keys[c.Mark]; !ok {
				t.Errorf("unable to find antenna mark %-5s", c.Mark)
			}
		}
	})

	t.Run("check for missing antenna assets", func(t *testing.T) {
		var assets meta.AssetList
		loadListFile(t, "../assets/antennas.csv", &assets)

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
	})

	t.Run("check for missing antenna sessions", func(t *testing.T) {
		var sessions meta.SessionList
		loadListFile(t, "../install/sessions.csv", &sessions)
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
	})
}
