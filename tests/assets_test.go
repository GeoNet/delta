package delta_test

import (
	"sort"
	"testing"

	"github.com/GeoNet/delta/meta"
)

var testAssets = map[string]func([]meta.Asset) func(t *testing.T){

	"check for duplicate asset numbers": func(assets []meta.Asset) func(t *testing.T) {
		return func(t *testing.T) {
			reference := make(map[string]string)
			for _, a := range assets {
				if a.Number == "" {
					continue
				}
				if x, ok := reference[a.Number]; ok {
					t.Errorf("duplicate asset number %s: %s [%s]", a.String(), a.Number, x)
				}
				reference[a.Number] = a.String()
			}
		}
	},

	"check for duplicate equipment": func(assets []meta.Asset) func(t *testing.T) {
		return func(t *testing.T) {
			for i := 0; i < len(assets); i++ {
				for j := i + 1; j < len(assets); j++ {
					switch {
					case assets[i].Model != assets[j].Model:
					case assets[i].Serial != assets[j].Serial:
					default:
						t.Errorf("duplicate equipment %s: %s", assets[i].String(), assets[j].String())
					}
				}
			}
		}
	},
}

func TestAssets(t *testing.T) {

	var assets meta.AssetList

	files := map[string]string{
		"antennas":    "../assets/antennas.csv",
		"cameras":     "../assets/cameras.csv",
		"dataloggers": "../assets/dataloggers.csv",
		"metsensors":  "../assets/metsensors.csv",
		"radomes":     "../assets/radomes.csv",
		"receivers":   "../assets/receivers.csv",
		"recorders":   "../assets/recorders.csv",
		"sensors":     "../assets/sensors.csv",
	}

	for k, v := range files {
		t.Run("load asset files: "+k, func(t *testing.T) {
			var a meta.AssetList
			loadListFile(t, v, &a)
			assets = append(assets, a...)
		})
	}

	sort.Sort(assets)

	for k, fn := range testAssets {
		t.Run(k, fn(assets))
	}
}
