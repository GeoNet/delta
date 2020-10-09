package delta_test

import (
	"testing"
	"time"

	"github.com/GeoNet/delta/meta"
)

var testFirmwareHistory = map[string]func([]meta.FirmwareHistory) func(t *testing.T){

	"check for firmware history overlaps": func(firmwares []meta.FirmwareHistory) func(t *testing.T) {
		return func(t *testing.T) {

			installs := make(map[string]meta.FirmwareHistoryList)
			for _, s := range firmwares {
				if _, ok := installs[s.Model]; !ok {
					installs[s.Model] = meta.FirmwareHistoryList{}
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

						t.Errorf("firmware %s / %s has overlap between %s and %s",
							v[i].Model, v[i].Serial, v[i].Start.Format(meta.DateTimeFormat), v[i].End.Format(meta.DateTimeFormat))
					}
				}
			}
		}
	},

	"check for firmware non-changes": func(firmwares []meta.FirmwareHistory) func(t *testing.T) {
		return func(t *testing.T) {

			for i := 0; i < len(firmwares)-1; i++ {
				current, next := firmwares[i], firmwares[i+1]

				switch {
				case current.Model != next.Model:
				case current.Make != next.Make:
				case current.Serial != next.Serial:
				case current.Version != next.Version:
				default:
					t.Errorf("likely invalid firmware change (line %d): %s %s %s %s %s",
						i+2, current.Make, current.Model, current.Start, current.Serial, current.Version)
				}
			}
		}
	},
}

var testFirmwareHistoryAssets = map[string]func([]meta.FirmwareHistory, []meta.Asset) func(t *testing.T){
	"check for firmware assets": func(firmwares []meta.FirmwareHistory, assets []meta.Asset) func(t *testing.T) {
		return func(t *testing.T) {
			for _, r := range firmwares {
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
				t.Errorf("unable to find firmware receiver asset: %s [%s]", r.Model, r.Serial)
			}
		}
	},
}

var testFirmwareHistoryDeployedReceiver = map[string]func([]meta.FirmwareHistory, []meta.DeployedReceiver) func(t *testing.T){
	"check for latest installed receiver firmware": func(firmwares []meta.FirmwareHistory, receivers []meta.DeployedReceiver) func(t *testing.T) {
		return func(t *testing.T) {

			installs := make(map[string]meta.FirmwareHistory)
			for _, s := range firmwares {
				if s.End.Before(time.Now()) {
					continue
				}
				installs[s.Model+"/"+s.Serial] = s
			}
			for _, r := range receivers {
				if r.End.Before(time.Now()) {
					continue
				}
				if _, ok := installs[r.Model+"/"+r.Serial]; !ok {
					t.Errorf("deployed receiver has no current firmware %s / %s at %s between %s and %s",
						r.Model, r.Serial, r.Mark,
						r.Start.Format(meta.DateTimeFormat),
						r.End.Format(meta.DateTimeFormat))
				}
			}
		}
	},
}

func TestFirmwareHistory(t *testing.T) {

	var firmwares meta.FirmwareHistoryList
	loadListFile(t, "../install/firmware.csv", &firmwares)

	for k, fn := range testFirmwareHistory {
		t.Run(k, fn(firmwares))
	}
}

func TestFirmwareHistory_Assets(t *testing.T) {

	var firmwares meta.FirmwareHistoryList
	loadListFile(t, "../install/firmware.csv", &firmwares)

	var assets meta.AssetList

	files := map[string]string{
		"cameras":     "../assets/cameras.csv",
		"dataloggers": "../assets/dataloggers.csv",
		"receivers":   "../assets/receivers.csv",
		"recorders":   "../assets/recorders.csv",
	}

	for _, v := range files {
		var a meta.AssetList
		loadListFile(t, v, &a)

		assets = append(assets, a...)
	}

	for k, fn := range testFirmwareHistoryAssets {
		t.Run(k, fn(firmwares, assets))
	}
}

func TestFirmwareHistory_DeployedReceiver(t *testing.T) {

	var firmwares meta.FirmwareHistoryList
	loadListFile(t, "../install/firmware.csv", &firmwares)

	var receivers meta.DeployedReceiverList
	loadListFile(t, "../install/receivers.csv", &receivers)

	for k, fn := range testFirmwareHistoryDeployedReceiver {
		t.Run(k, fn(firmwares, receivers))
	}
}
