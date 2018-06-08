package delta_test

import (
	"testing"

	"github.com/GeoNet/delta/meta"
)

func TestReceivers(t *testing.T) {

	var receivers meta.DeployedReceiverList
	loadListFile(t, "../install/receivers.csv", &receivers)

	t.Run("check for particular receiver installation overlaps", func(t *testing.T) {
		installs := make(map[string]meta.DeployedReceiverList)
		for _, s := range receivers {
			if _, ok := installs[s.Model]; !ok {
				installs[s.Model] = meta.DeployedReceiverList{}
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

					t.Errorf("receiver %s [%s] at %s has overlap with %s between times %s and %s",
						v[i].Model, v[i].Serial, v[i].Mark, v[j].Mark, v[i].Start.Format(meta.DateTimeFormat), v[i].End.Format(meta.DateTimeFormat))
				}
			}
		}
	})

	t.Run("check for receiver installation equipment overlaps", func(t *testing.T) {
		installs := make(map[string]meta.DeployedReceiverList)
		for _, s := range receivers {
			if _, ok := installs[s.Mark]; !ok {
				installs[s.Mark] = meta.DeployedReceiverList{}
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

					t.Errorf("receivers %s [%s] / %s [%s] at %s has overlap between %s and %s",
						v[i].Model, v[i].Serial, v[j].Model, v[j].Serial, v[i].Mark, v[i].Start.Format(meta.DateTimeFormat), v[i].End.Format(meta.DateTimeFormat))
				}
			}
		}
	})

	t.Run("check for missing receiver marks", func(t *testing.T) {
		var marks meta.MarkList
		loadListFile(t, "../network/marks.csv", &marks)

		keys := make(map[string]interface{})
		for _, m := range marks {
			keys[m.Code] = true
		}

		for _, r := range receivers {
			if _, ok := keys[r.Mark]; ok {
				continue
			}
			t.Errorf("unable to find receiver mark %-5s", r.Mark)
		}
	})

	t.Run("check for missing receiver assets", func(t *testing.T) {
		var assets meta.AssetList
		loadListFile(t, "../assets/receivers.csv", &assets)
		for _, r := range receivers {
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
			t.Errorf("unable to find receiver asset: %s [%s]", r.Model, r.Serial)
		}
	})

}
