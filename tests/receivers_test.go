package delta_test

import (
	"testing"

	"github.com/GeoNet/delta/meta"
)

var testDeployedReceivers = map[string]func([]meta.DeployedReceiver) func(t *testing.T){
	"check for receiver installation overlaps": func(receivers []meta.DeployedReceiver) func(t *testing.T) {
		return func(t *testing.T) {

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
		}
	},

	"check for receiver installation equipment overlaps": func(receivers []meta.DeployedReceiver) func(t *testing.T) {
		return func(t *testing.T) {

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
		}
	},
}

var testDeployedReceiversMarks = map[string]func([]meta.DeployedReceiver, []meta.Mark) func(t *testing.T){
	"check for missing receiver marks": func(receivers []meta.DeployedReceiver, marks []meta.Mark) func(t *testing.T) {
		return func(t *testing.T) {

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
		}
	},
}

var testDeployedReceiversAssets = map[string]func([]meta.DeployedReceiver, []meta.Asset) func(t *testing.T){
	"check for missing receiver assets": func(receivers []meta.DeployedReceiver, assets []meta.Asset) func(t *testing.T) {
		return func(t *testing.T) {
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
		}
	},
}

func TestDeployedReceivers(t *testing.T) {
	var receivers meta.DeployedReceiverList
	loadListFile(t, "../install/receivers.csv", &receivers)

	for k, fn := range testDeployedReceivers {
		t.Run(k, fn(receivers))
	}
}

func TestDeployedReceiversMarks(t *testing.T) {
	var receivers meta.DeployedReceiverList
	loadListFile(t, "../install/receivers.csv", &receivers)

	var marks meta.MarkList
	loadListFile(t, "../network/marks.csv", &marks)

	for k, fn := range testDeployedReceiversMarks {
		t.Run(k, fn(receivers, marks))
	}
}

func TestDeployedReceiversAssets(t *testing.T) {
	var receivers meta.DeployedReceiverList
	loadListFile(t, "../install/receivers.csv", &receivers)

	var assets meta.AssetList
	loadListFile(t, "../assets/receivers.csv", &assets)

	for k, fn := range testDeployedReceiversAssets {
		t.Run(k, fn(receivers, assets))
	}
}
